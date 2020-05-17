package main

import (
	"context"
	"fmt"
	"github.com/firefly-crm/common/infra"
	"github.com/firefly-crm/common/logger"
	"github.com/firefly-crm/common/rabbit"
	"github.com/firefly-crm/common/rabbit/exchanges"
	"github.com/firefly-crm/common/rabbit/routes"
	"github.com/firefly-crm/fireflycrm-bot-backend/orderbook"
	"github.com/firefly-crm/fireflycrm-bot-backend/service"
	"github.com/firefly-crm/fireflycrm-bot-backend/storage"
	"github.com/firefly-crm/fireflycrm-bot-backend/users"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			logger.Crashf("service exited with error: %v", err)
		}
	}()

	infraCtx := infra.Context()

	//TODO: Config map
	rabbitUsername := os.Getenv("RMQ_USERNAME")
	if rabbitUsername == "" {
		log.Fatalf("rabbit username is empty")
	}

	rabbitPassword := os.Getenv("RMQ_PASSWORD")
	if rabbitPassword == "" {
		log.Fatal("rabbit password is empty")
	}

	rabbitHost := os.Getenv("RMQ_HOST")
	if rabbitHost == "" {
		log.Fatalf("rabbit host is empty")
	}

	rabbitPort := os.Getenv("RMQ_PORT")
	if rabbitPort == "" {
		log.Fatalf("rabbit port is empty")
	}

	rabbitConnString := fmt.Sprintf("amqp://%s:%s@%s:%s", rabbitUsername, rabbitPassword, rabbitHost, rabbitPort)

	rabbitConfig := rabbit.Config{
		Endpoint: rabbitConnString,
	}
	rabbitPrimary := rabbit.MustNew(rabbitConfig)

	go func() {
		errPrimary := <-rabbitPrimary.Done()
		logger.Crashf("primary rabbit client error: %v", errPrimary)
	}()

	tgExchange := exchanges.MustExchangeByID(exchanges.FireflyCRMTelegramUpdates)
	rp := rabbitPrimary.MustNewExchange(tgExchange.Opts)

	pgHost := os.Getenv("POSTGRES_HOST")
	if pgHost == "" {
		panic("pg host is unset; use POSTGRES_HOST env")
	}

	pgUser := os.Getenv("POSTGRES_USER")
	if pgUser == "" {
		panic("pg username is unset; use POSTGRES_USER env")
	}

	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	if pgPassword == "" {
		panic("pg password is unset; use POSTGRES_PASSWORD env")
	}

	pgDBName := os.Getenv("POSTGRES_DB")
	if pgDBName == "" {
		panic("pg db is unset; user POSTGRES_DB env")
	}

	pgPort := "5432"
	envPort := os.Getenv("POSTGRES_PORT")
	if envPort != "" {
		pgPort = envPort
	}

	pgConnString := fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=disable", pgUser, pgPassword, pgDBName, pgPort, pgHost)

	tgToken := os.Getenv("TG_TOKEN")
	if tgToken == "" {
		log.Fatalf("telegram token is not set")
	}

	db, err := sqlx.Connect("postgres", pgConnString)
	if err != nil {
		panic(err)
	}
	stor := storage.NewStorage(db)
	ob := orderbook.MustNewOrderBook(stor)
	u := users.NewUsers(stor)

	bot := service.MustNewBot(tgToken)
	backend := service.Service{
		Bot:       bot,
		OrderBook: ob,
		Users:     u,
		Storage:   stor,
	}

	errGroup, ctx := errgroup.WithContext(infraCtx)

	errGroup.Go(func() error {
		return rp.Queue(routes.TelegramCallbackUpdate).Consume(ctx, backend.ProcessCallbackEvent)
	})
	errGroup.Go(func() error { return rp.Queue(routes.TelegramCommandUpdate).Consume(ctx, backend.ProcessCommandEvent) })
	errGroup.Go(func() error { return rp.Queue(routes.TelegramPromptUpdate).Consume(ctx, backend.ProcessPromptEvent) })

	if err := errGroup.Wait(); err != nil && err != context.Canceled {
		panic(err)
	}
}
