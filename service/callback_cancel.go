package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
)

func (s Service) processCancelCallback(ctx context.Context, messageId uint64) error {

	order, err := s.OrderBook.GetOrderByMessageId(ctx, messageId)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	err = s.OrderBook.UpdateOrderEditState(ctx, order.Id, types.EditStateNone)
	if err != nil {
		return fmt.Errorf("failed to update order state: %w", err)
	}

	if order.ActiveItemId.Valid {
		receiptItem, err := s.OrderBook.GetReceiptItem(ctx, uint64(order.ActiveItemId.Int64))
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to get receipt item: %w", err)
		}

		if err == nil && !receiptItem.Initialised {
			err = s.OrderBook.RemoveItem(ctx, uint64(order.ActiveItemId.Int64))
			if err != nil {
				return fmt.Errorf("failed to delete receipt item: %w", err)
			}
		}
	}

	if order.ActivePaymentId.Valid {
		var payment types.Payment
		for _, p := range order.Payments {
			if p.Id == uint64(order.ActivePaymentId.Int64) {
				payment = p
				break
			}
		}

		if payment.Amount == 0 {
			err = s.OrderBook.RemovePayment(ctx, uint64(order.ActivePaymentId.Int64))
			if err != nil {
				return fmt.Errorf("failed to delete payment item: %w", err)
			}
		}
	}

	err = s.updateOrderMessage(ctx, messageId, true)
	if err != nil {
		return fmt.Errorf("failed to update order message: %w", err)
	}

	err = s.deleteHint(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to delete hint: %w", err)
	}

	return nil
}
