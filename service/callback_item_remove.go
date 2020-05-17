package service

import (
	"context"
	"fmt"
	tp "github.com/firefly-crm/common/messages/telegram"
)

func (s Service) processItemRemove(ctx context.Context, callback *tp.CallbackEvent) error {
	err := s.OrderBook.RemoveItem(ctx, callback.EntityId)
	if err != nil {
		return fmt.Errorf("failed to remove receipt item: %w", err)
	}

	err = s.updateOrderMessage(ctx, callback.MessageId, true)
	if err != nil {
		return fmt.Errorf("failed to refresh order message: %w", err)
	}

	return nil
}
