package service

import (
	"context"
	"fmt"
	tp "github.com/firefly-crm/common/messages/telegram"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
)

func (s Service) processRefundCallback(ctx context.Context, order types.Order, userId, messageId uint64, amount uint32) error {
	//TODO: Refund payment at ModulBank

	if !order.ActivePaymentId.Valid {
		return fmt.Errorf("active payment id is nil")
	}

	paymentId := uint64(order.ActivePaymentId.Int64)
	if amount == 0 {
		for _, p := range order.Payments {
			if p.Id == paymentId {
				amount = p.Amount
			}
		}
	}

	err := s.OrderBook.RefundPayment(ctx, paymentId, amount)
	if err != nil {
		return fmt.Errorf("failed to refund payment: %w", err)
	}

	err = s.updateOrderMessage(ctx, userId, messageId, nil)
	if err != nil {
		return fmt.Errorf("failed to update order message: %w", err)
	}

	return nil
}

func (s Service) processPartialRefundCallback(ctx context.Context, callback *tp.CallbackEvent) error {
	messageId := callback.MessageId

	order, err := s.OrderBook.GetOrderByMessageId(ctx, callback.UserId, messageId)
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}

	err = s.OrderBook.UpdateOrderEditState(ctx, order.Id, types.EditStateWaitingRefundAmount)
	if err != nil {
		return fmt.Errorf("failed to update order state: %w", err)
	}

	return nil
}
