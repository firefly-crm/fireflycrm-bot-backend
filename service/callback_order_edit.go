package service

import (
	"context"
	"fmt"
	tp "github.com/firefly-crm/common/messages/telegram"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
)

func (s Service) processOrderEditDueDate(ctx context.Context, callback *tp.CallbackEvent) error {
	order, err := s.OrderBook.GetOrderByMessageId(ctx, callback.UserId, callback.MessageId)
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}

	err = s.OrderBook.UpdateOrderEditState(ctx, order.Id, types.EditStateWaitingOrderDueDate)
	if err != nil {
		return fmt.Errorf("failed to update order state: %w", err)
	}

	return nil
}

func (s Service) processOrderEditDescription(ctx context.Context, callback *tp.CallbackEvent) error {
	order, err := s.OrderBook.GetOrderByMessageId(ctx, callback.UserId, callback.MessageId)
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}

	err = s.OrderBook.UpdateOrderEditState(ctx, order.Id, types.EditStateWaitingOrderDescription)
	if err != nil {
		return fmt.Errorf("failed to update order state: %w", err)
	}

	return nil
}
