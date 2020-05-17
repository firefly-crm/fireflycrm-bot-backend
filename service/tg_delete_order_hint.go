package service

import (
	"context"
	"fmt"
	"github.com/firefly-crm/fireflycrm-bot-backend/common/logger"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (s Service) deleteHint(ctx context.Context, bot *tg.BotAPI, order types.Order) error {
	log := logger.FromContext(ctx)

	if !order.HintMessageId.Valid {
		log.Infof("hint message already deleted")
		return nil
	}

	deleteMessage := tg.NewDeleteMessage(int64(order.UserId), int(order.HintMessageId.Int64))
	_, err := bot.Send(deleteMessage)
	if err != nil {
		return fmt.Errorf("failed to delete hind: %w", err)
	}

	err = s.OrderBook.UpdateHintMessageForOrder(ctx, order.Id, 0)
	if err != nil {
		return fmt.Errorf("failed to null hint message for order: %w", err)
	}

	return nil
}
