package service

import (
	"context"
	"fmt"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
)

func (s Service) processOrderDisplayModeCallback(ctx context.Context, userId, messageId uint64, mode types.DisplayMode) error {
	err := s.OrderBook.UpdateOrderMessageDisplayMode(ctx, userId, messageId, mode)
	if err != nil {
		return fmt.Errorf("failed to update display mode: %w", err)
	}

	err = s.updateOrderMessage(ctx, userId, messageId, true)
	if err != nil {
		return fmt.Errorf("failed to update order message: %w", err)
	}

	return nil
}
