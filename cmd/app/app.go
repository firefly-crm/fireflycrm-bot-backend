package main

import (
	"context"
	"github.com/firefly-crm/common/infra"
	"github.com/firefly-crm/common/logger"
	"github.com/firefly-crm/common/rabbit"
	"github.com/firefly-crm/common/rabbit/exchanges"
	"github.com/firefly-crm/common/rabbit/routes"
	"github.com/firefly-crm/fireflycrm-bot-backend/config"
	"github.com/firefly-crm/fireflycrm-bot-backend/orderbook"
	"github.com/firefly-crm/fireflycrm-bot-backend/service"
	"github.com/firefly-crm/fireflycrm-bot-backend/storage"
	"github.com/firefly-crm/fireflycrm-bot-backend/users"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			logger.Crashf("service exited with error: %v", err)
		}
	}()

	infraCtx := infra.Context()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	serviceConfig := config.Config{}
	err = viper.Unmarshal(&serviceConfig)
	if err != nil {
		panic(err)
	}

	rabbitConfig := rabbit.Config{
		Endpoint: serviceConfig.Rabbit,
	}
	rabbitPrimary := rabbit.MustNew(rabbitConfig)

	go func() {
		errPrimary := <-rabbitPrimary.Done()
		logger.Crashf("primary rabbit client error: %v", errPrimary)
	}()

	tgExchange := exchanges.MustExchangeByID(exchanges.FireflyCRMTelegramUpdates)
	rp := rabbitPrimary.MustNewExchange(tgExchange.Opts)

	db, err := sqlx.Connect("postgres", serviceConfig.Db)
	if err != nil {
		panic(err)
	}
	stor := storage.NewStorage(db)
	ob := orderbook.MustNewOrderBook(stor)
	u := users.NewUsers(stor)

	bot := service.MustNewBot(serviceConfig.TgToken)
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
