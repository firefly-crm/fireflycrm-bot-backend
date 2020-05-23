package service

import (
	"context"
	"fmt"
	tg "github.com/DarthRamone/telegram-bot-api"
	bot "github.com/firefly-crm/common/bot"
	tp "github.com/firefly-crm/common/messages/telegram"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
)

func (s Service) processCustomerEditInstagram(ctx context.Context, callback *tp.CallbackEvent) error {
	return s.processCustomerEditField(ctx, callback, bot.ReplyEnterCustomerInstagram, types.EditStateWaitingCustomerInstagram)
}

func (s Service) processCustomerEditEmail(ctx context.Context, callback *tp.CallbackEvent) error {
	return s.processCustomerEditField(ctx, callback, bot.ReplyEnterCustomerEmail, types.EditStateWaitingCustomerEmail)
}

func (s Service) processCustomerEditPhone(ctx context.Context, callback *tp.CallbackEvent) error {
	return s.processCustomerEditField(ctx, callback, bot.ReplyEnterCustomerPhone, types.EditStateWaitingCustomerPhone)
}

func (s Service) processCustomerEditDescription(ctx context.Context, callback *tp.CallbackEvent) error {
	return s.processCustomerEditField(ctx, callback, bot.ReplyEnterCustomerDescription, types.EditStateWaitingCustomerDescription)
}

func (s Service) processCustomerEditField(ctx context.Context, callback *tp.CallbackEvent, hintText string, state types.EditState) error {
	userId := int64(callback.UserId)
	messageId := callback.MessageId

	order, err := s.OrderBook.GetOrderByMessageId(ctx, uint64(userId), messageId)
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}
	hintMessage := tg.NewMessage(userId, hintText)
	hint, err := s.Bot.Send(hintMessage)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	err = s.OrderBook.UpdateHintMessageForOrder(ctx, order.Id, uint64(hint.MessageID))
	if err != nil {
		return fmt.Errorf("failed to update hint message: %w", err)
	}

	err = s.OrderBook.UpdateOrderEditState(ctx, order.Id, state)
	if err != nil {
		return fmt.Errorf("failed to update order state: %w", err)
	}

	return nil
}
