package service

import (
	"context"
	"github.com/firefly-crm/fireflycrm-bot-backend/billmaker"
	"github.com/firefly-crm/fireflycrm-bot-backend/orderbook"
	"github.com/firefly-crm/fireflycrm-bot-backend/storage"
	"github.com/firefly-crm/fireflycrm-bot-backend/users"
)

type (
	Service struct {
		OrderBook orderbook.OrderBook
		BillMaker billmaker.BillMaker
		Users     users.Users
		Storage   storage.Storage
	}

	Options struct {
		TelegramToken string
	}
)

func (s Service) Serve(ctx context.Context, opts Options) error {
	bot := s.startListenTGUpdates(ctx, opts.TelegramToken)
	s.startPaymentsWatcher(ctx, bot)
	return nil
}
