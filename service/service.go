package service

import (
	tg "github.com/DarthRamone/telegram-bot-api"
	"github.com/firefly-crm/fireflycrm-bot-backend/orderbook"
	"github.com/firefly-crm/fireflycrm-bot-backend/storage"
	"github.com/firefly-crm/fireflycrm-bot-backend/users"
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
