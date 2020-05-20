package service

import (
	"context"
	"fmt"
	tg "github.com/DarthRamone/telegram-bot-api"
	tp "github.com/firefly-crm/common/messages/telegram"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
)

func (s Service) processOrderEditDueDate(ctx context.Context, callback *tp.CallbackEvent) error {
	order, err := s.OrderBook.GetOrderByMessageId(ctx, callback.UserId, callback.MessageId)
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}
	hintMessage := tg.NewMessage(int64(callback.UserId), replyEnterOrderDueDate)
	hint, err := s.Bot.Send(hintMessage)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	err = s.OrderBook.UpdateHintMessageForOrder(ctx, order.Id, uint64(hint.MessageID))
	if err != nil {
		return fmt.Errorf("failed to update hint message: %w", err)
	}

	err = s.OrderBook.UpdateOrderEditState(ctx, order.Id, types.EditStateWaitingOrderDueDate)
	if err != nil {
		return fmt.Errorf("failed to update order state: %w", err)
	}

	return nil
}
