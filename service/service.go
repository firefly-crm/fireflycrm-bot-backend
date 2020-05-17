package service

import (
	"github.com/firefly-crm/fireflycrm-bot-backend/orderbook"
	"github.com/firefly-crm/fireflycrm-bot-backend/storage"
	"github.com/firefly-crm/fireflycrm-bot-backend/users"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

type (
	Service struct {
		Bot       *tg.BotAPI
		OrderBook orderbook.OrderBook
		Users     users.Users
		Storage   storage.Storage
	}

	Options struct {
		TelegramToken string
	}
)
