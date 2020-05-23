package service

import (
	"context"
	"fmt"
	tg "github.com/DarthRamone/telegram-bot-api"
	"github.com/firefly-crm/common/bot"
	tp "github.com/firefly-crm/common/messages/telegram"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
)

func (s Service) processAddKnownItem(ctx context.Context, callback *tp.CallbackEvent, data string) error {
	userId := int64(callback.UserId)
	messageId := callback.MessageId

	order, err := s.OrderBook.GetOrderByMessageId(ctx, callback.UserId, messageId)
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}

	hintMessage := tg.NewMessage(userId, bot.ReplyEnterItemPrice)
	hint, err := s.Bot.Send(hintMessage)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	name := "Unknown item"
	itemType := types.ReceiptItemTypeGoods
	switch data {
	case bot.KbDataDelivery:
		name = "Доставка"
		itemType = types.ReceiptItemTypeDelivery
	case bot.KbDataLingerieSet:
		name = "Комплект нижнего белья"
	}

	var itemId uint64
	itemId, err = s.OrderBook.AddItem(ctx, order.Id, itemType)
	if err != nil {
		return fmt.Errorf("failed to add item to order")
	}

	err = s.OrderBook.UpdateReceiptItemName(ctx, name, uint64(userId), itemId)
	if err != nil {
		return fmt.Errorf("failed to set delivery name: %w", err)
	}

	err = s.OrderBook.UpdateHintMessageForOrder(ctx, order.Id, uint64(hint.MessageID))
	if err != nil {
		return fmt.Errorf("failed to update hint message: %w", err)
	}

	err = s.OrderBook.UpdateOrderEditState(ctx, order.Id, types.EditStateWaitingItemPrice)
	if err != nil {
		return fmt.Errorf("failed to update order state: %w", err)
	}

	return nil
}

func (s Service) processAddItemCallack(ctx context.Context, callback *tp.CallbackEvent) error {
	userId := int64(callback.UserId)
	messageId := callback.MessageId

	order, err := s.OrderBook.GetOrderByMessageId(ctx, callback.UserId, messageId)
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}
	hintMessage := tg.NewMessage(userId, bot.ReplyEnterItemName)
	hint, err := s.Bot.Send(hintMessage)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	_, err = s.OrderBook.AddItem(ctx, order.Id, types.ReceiptItemTypeGoods)
	if err != nil {
		return fmt.Errorf("failed to add item to order")
	}

	err = s.OrderBook.UpdateHintMessageForOrder(ctx, order.Id, uint64(hint.MessageID))
	if err != nil {
		return fmt.Errorf("failed to update hint message: %w", err)
	}

	err = s.OrderBook.UpdateOrderEditState(ctx, order.Id, types.EditStateWaitingItemName)
	if err != nil {
		return fmt.Errorf("failed to update order state: %w", err)
	}

	return nil
}
