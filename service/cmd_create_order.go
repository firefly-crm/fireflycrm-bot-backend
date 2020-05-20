package service

import (
	"context"
	"fmt"
	tg "github.com/DarthRamone/telegram-bot-api"
	"github.com/firefly-crm/common/logger"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
)

func (s Service) createOrder(ctx context.Context, userId, messageId uint64) error {
	uid := int64(userId)

	log := logger.FromContext(ctx)

	order, err := s.OrderBook.CreateOrder(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	messageText := order.MessageString(nil, types.DisplayModeFull)

	deleteMessage := tg.NewDeleteMessage(uid, int(messageId))
	_, err = s.Bot.DeleteMessage(deleteMessage)
	if err != nil {
		log.Warnf("failed to delete command message: %v", err)
	}

	msg := tg.NewMessage(uid, messageText)
	msg.ParseMode = "markdown"

	var orderMessage tg.Message
	orderMessage, err = s.Bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	err = s.OrderBook.UpdateMessageForOrder(ctx, userId, order.Id, uint64(orderMessage.MessageID))
	if err != nil {
		return fmt.Errorf("failed to update message id for order: %v", err)
	}

	messageReplyMarkup, err := startOrderInlineKeyboard(ctx, s, userId, uint64(orderMessage.MessageID))
	if err != nil {
		return fmt.Errorf("failed to get start order inline markup: %w", err)
	}

	editMessage := tg.NewEditMessageReplyMarkup(uid, orderMessage.MessageID, messageReplyMarkup)
	_, err = s.Bot.Send(editMessage)
	if err != nil {
		return fmt.Errorf("failed to set new markup: %w", err)
	}

	return nil
}
