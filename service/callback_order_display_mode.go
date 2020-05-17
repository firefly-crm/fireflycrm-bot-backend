package service

import (
	"context"
	"fmt"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
)

func (s Service) processOrderDisplayModeCallback(ctx context.Context, messageId uint64, mode types.DisplayMode) error {
	err := s.OrderBook.UpdateOrderMessageDisplayMode(ctx, messageId, mode)
	if err != nil {
		return fmt.Errorf("failed to update display mode: %w", err)
	}

	err = s.updateOrderMessage(ctx, messageId, true)
	if err != nil {
		return fmt.Errorf("failed to update order message: %w", err)
	}

	return nil
}
