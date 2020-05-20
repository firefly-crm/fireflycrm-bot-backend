package service

import (
	"context"
	"fmt"
	tg "github.com/DarthRamone/telegram-bot-api"
	"github.com/firefly-crm/common/logger"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
)

func (s Service) deleteHint(ctx context.Context, order types.Order) error {
	log := logger.FromContext(ctx)

	if !order.HintMessageId.Valid {
		log.Infof("hint message already deleted")
		return nil
	}

	deleteMessage := tg.NewDeleteMessage(int64(order.UserId), int(order.HintMessageId.Int64))
	_, err := s.Bot.Send(deleteMessage)
	if err != nil {
		return fmt.Errorf("failed to delete hint: %w", err)
	}

	err = s.OrderBook.UpdateHintMessageForOrder(ctx, order.Id, 0)
	if err != nil {
		return fmt.Errorf("failed to null hint message for order: %w", err)
	}

	return nil
}
