package service

import (
	"context"
	"fmt"
	tp "github.com/firefly-crm/common/messages/telegram"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
)

func (s Service) processAddPaymentCallback(ctx context.Context, userId, messageId uint64, method types.PaymentMethod) error {
	order, err := s.OrderBook.GetOrderByMessageId(ctx, userId, messageId)
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}

	_, err = s.OrderBook.AddPayment(ctx, order.Id, method)
	if err != nil {
		return fmt.Errorf("failed to add payment to order: %w", err)
	}

	return nil
}

func (s Service) processPartialPaymentCallback(ctx context.Context, callback *tp.CallbackEvent) error {
	userId := int64(callback.UserId)
	messageId := callback.MessageId

	order, err := s.OrderBook.GetOrderByMessageId(ctx, uint64(userId), messageId)
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}

	err = s.OrderBook.UpdateOrderEditState(ctx, order.Id, types.EditStateWaitingPaymentAmount)
	if err != nil {
		return fmt.Errorf("failed to update order state: %w", err)
	}

	return nil
}

//if amount is 0 then full payment
func (s Service) processPaymentCallback(ctx context.Context, userId, messageId uint64, amount uint32) error {
	order, err := s.OrderBook.GetOrderByMessageId(ctx, userId, messageId)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	if !order.ActivePaymentId.Valid {
		return fmt.Errorf("active bill id is nil")
	}

	paymentId := uint64(order.ActivePaymentId.Int64)
	if amount == 0 {
		amount = order.Amount - order.PayedAmount
	}

	var payment types.Payment
	for _, p := range order.Payments {
		if p.Id == paymentId {
			payment = p
			break
		}
	}

	err = s.OrderBook.UpdatePaymentAmount(ctx, paymentId, amount)
	if err != nil {
		return fmt.Errorf("failed to update payment amount: %w", err)
	}

	if payment.PaymentMethod == types.PaymentMethodAcquiring {
		err := s.OrderBook.GeneratePaymentLink(ctx, paymentId)
		if err != nil {
			return fmt.Errorf("failed to generate payment link: %w", err)
		}
	}

	err = s.updateOrderMessage(ctx, userId, messageId, nil)
	if err != nil {
		return fmt.Errorf("failed to update order message: %w", err)
	}

	return nil
}
