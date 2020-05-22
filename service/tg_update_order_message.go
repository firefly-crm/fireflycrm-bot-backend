package service

import (
	"context"
	"fmt"
	tg "github.com/DarthRamone/telegram-bot-api"
	"github.com/firefly-crm/common/logger"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
)

func (s Service) updateOrderMessage(ctx context.Context, userId, messageId uint64, flowCompleted bool) error {
	log := logger.FromContext(ctx)

	order, err := s.OrderBook.GetOrderByMessageId(ctx, userId, messageId)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	orderMessage, err := s.OrderBook.GetOrderMessage(ctx, userId, messageId)
	if err != nil {
		return fmt.Errorf("failed to get order message: %w", err)
	}

	var customer *types.Customer
	if order.CustomerId.Valid {
		c, err := s.Users.GetCustomer(ctx, uint64(order.CustomerId.Int64))
		if err != nil {
			log.Errorf("failed to get customer: %w", err)
		}
		customer = &c
	} else {
		log.Warnf("customer is nil")
	}

	chatId := int64(order.UserId)

	editMessage := tg.NewEditMessageText(chatId, int(messageId), order.MessageString(customer, orderMessage.DisplayMode))
	editMessage.ParseMode = "html"
	editMessage.DisableWebPagePreview = true
	var markup tg.InlineKeyboardMarkup
	if flowCompleted {
		markup, err = startOrderInlineKeyboard(ctx, s, userId, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	} else {
		markup = cancelInlineKeyboard()
	}
	editMessage.ReplyMarkup = &markup

	_, err = s.Bot.Send(editMessage)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
