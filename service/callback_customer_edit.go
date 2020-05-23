package service

import (
	"context"
	"fmt"
	tp "github.com/firefly-crm/common/messages/telegram"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
)

func (s Service) processCustomerEditInstagram(ctx context.Context, callback *tp.CallbackEvent) error {
	return s.processCustomerEditField(ctx, callback, types.EditStateWaitingCustomerInstagram)
}

func (s Service) processCustomerEditEmail(ctx context.Context, callback *tp.CallbackEvent) error {
	return s.processCustomerEditField(ctx, callback, types.EditStateWaitingCustomerEmail)
}

func (s Service) processCustomerEditPhone(ctx context.Context, callback *tp.CallbackEvent) error {
	return s.processCustomerEditField(ctx, callback, types.EditStateWaitingCustomerPhone)
}

func (s Service) processCustomerEditDescription(ctx context.Context, callback *tp.CallbackEvent) error {
	return s.processCustomerEditField(ctx, callback, types.EditStateWaitingCustomerDescription)
}

func (s Service) processCustomerEditField(ctx context.Context, callback *tp.CallbackEvent, state types.EditState) error {
	userId := int64(callback.UserId)
	messageId := callback.MessageId

	order, err := s.OrderBook.GetOrderByMessageId(ctx, uint64(userId), messageId)
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}

	err = s.OrderBook.UpdateOrderEditState(ctx, order.Id, state)
	if err != nil {
		return fmt.Errorf("failed to update order state: %w", err)
	}

	return nil
}
