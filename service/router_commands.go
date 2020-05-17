package service

import (
	"context"
	"fmt"
	"github.com/firefly-crm/common/logger"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

func (s Service) processCommand(ctx context.Context, update tg.Update) error {
	var err error
	var cmd = update.Message.Text

	log := logger.
		FromContext(ctx).
		WithField("user_id", update.Message.From.ID).
		WithField("command", cmd)

	log.Infof("processing command")

	ctx = logger.ToContext(ctx, log)

	if cmd == "/start" {
		err = s.createUser(ctx, update)
	} else if cmd == kbCreateOrder {
		err = s.createOrder(ctx, update)
	} else if strings.HasPrefix(cmd, "/registerAsMerchant") {
		err = s.registerMerchant(ctx, update)
	} else if cmd == kbActiveOrders {
	} else {
		err = s.processPrompt(ctx, update)
	}

	if err != nil {
		return fmt.Errorf("failed process message: %w", err)
	}

	return nil
}
