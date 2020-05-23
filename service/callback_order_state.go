package service

import (
	"context"
	"fmt"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
)

func (s Service) processOrderStateCallback(ctx context.Context, userId, messageId uint64, state types.OrderState) error {
	order, err := s.OrderBook.GetOrderByMessageId(ctx, userId, messageId)
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}

	err = s.OrderBook.UpdateOrderState(ctx, order.Id, state)
	if err != nil {
		return fmt.Errorf("failed to update order state(%v): %w", state, err)
	}

	if state == types.OrderStateDeleted {
		err = s.OrderBook.UpdateOrderMessageDisplayMode(ctx, userId, messageId, types.DisplayModeDeleted)
		if err != nil {
			return fmt.Errorf("failed to update display mode: %w", err)
		}
	}

	if state == types.OrderStateForming {
		orderMessage, err := s.OrderBook.GetOrderMessage(ctx, userId, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order message: %w", err)
		}

		if orderMessage.DisplayMode == types.DisplayModeDeleted {
			err = s.OrderBook.UpdateOrderMessageDisplayMode(ctx, userId, messageId, types.DisplayModeFull)
			if err != nil {
				return fmt.Errorf("failed to update display mode: %w", err)
			}
		}
	}

	err = s.updateOrderMessage(ctx, userId, messageId, nil)
	if err != nil {
		return fmt.Errorf("failed to update order message: %w", err)
	}

	return nil
}
