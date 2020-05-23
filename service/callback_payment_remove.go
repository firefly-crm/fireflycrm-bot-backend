package service

import (
	"context"
	"fmt"
	tp "github.com/firefly-crm/common/messages/telegram"
)

func (s Service) processPaymentRemove(ctx context.Context, callback *tp.CallbackEvent) error {
	err := s.OrderBook.RemovePayment(ctx, callback.EntityId)
	if err != nil {
		return fmt.Errorf("failed to remove payment: %w", err)
	}

	err = s.updateOrderMessage(ctx, callback.UserId, callback.MessageId, nil)
	if err != nil {
		return fmt.Errorf("failed to refresh order message: %w", err)
	}

	return nil
}
