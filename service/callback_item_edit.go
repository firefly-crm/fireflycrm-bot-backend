package service

import (
	"context"
	"fmt"
	tp "github.com/firefly-crm/common/messages/telegram"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
)

func (s Service) processItemEditQty(ctx context.Context, callback *tp.CallbackEvent) error {
	order, err := s.OrderBook.GetOrderByMessageId(ctx, callback.UserId, callback.MessageId)
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}

	err = s.OrderBook.SetActiveItemId(ctx, order.Id, callback.EntityId)
	if err != nil {
		return fmt.Errorf("failed to set active item id: %w", err)
	}

	err = s.OrderBook.UpdateOrderEditState(ctx, order.Id, types.EditStateWaitingItemQuantity)
	if err != nil {
		return fmt.Errorf("failed to update order state: %w", err)
	}

	return nil
}

func (s Service) processItemEditPrice(ctx context.Context, callback *tp.CallbackEvent) error {
	order, err := s.OrderBook.GetOrderByMessageId(ctx, callback.UserId, callback.MessageId)
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}

	err = s.OrderBook.SetActiveItemId(ctx, order.Id, callback.EntityId)
	if err != nil {
		return fmt.Errorf("failed to set active item id: %w", err)
	}

	err = s.OrderBook.UpdateOrderEditState(ctx, order.Id, types.EditStateWaitingItemPrice)
	if err != nil {
		return fmt.Errorf("failed to update order state: %w", err)
	}

	return nil
}

func (s Service) processItemEditName(ctx context.Context, callback *tp.CallbackEvent) error {
	order, err := s.OrderBook.GetOrderByMessageId(ctx, callback.UserId, callback.MessageId)
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}

	err = s.OrderBook.SetActiveItemId(ctx, order.Id, callback.EntityId)
	if err != nil {
		return fmt.Errorf("failed to set active item id: %w", err)
	}

	err = s.OrderBook.UpdateOrderEditState(ctx, order.Id, types.EditStateWaitingItemName)
	if err != nil {
		return fmt.Errorf("failed to update order state: %w", err)
	}

	return nil
}
