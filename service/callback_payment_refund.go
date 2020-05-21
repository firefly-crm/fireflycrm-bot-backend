package service

import (
	"context"
	"fmt"
	tg "github.com/DarthRamone/telegram-bot-api"
	"github.com/firefly-crm/common/bot"
	"github.com/firefly-crm/common/logger"
	tp "github.com/firefly-crm/common/messages/telegram"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
)

func (s Service) processRefundCallback(ctx context.Context, order types.Order, userId, messageId uint64, amount uint32) error {
	log := logger.FromContext(ctx)

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

	defer func() {
		if err := s.deleteHint(ctx, order); err != nil {
			log.Errorf("failed to delete hint: %v", err.Error())
		}
	}()

	err := s.OrderBook.RefundPayment(ctx, paymentId, amount)
	if err != nil {
		return fmt.Errorf("failed to refund payment: %w", err)
	}

	err = s.updateOrderMessage(ctx, userId, messageId, true)
	if err != nil {
		return fmt.Errorf("failed to update order message: %w", err)
	}

	return nil
}

func (s Service) processPartialRefundCallback(ctx context.Context, callback *tp.CallbackEvent) error {
	chatId := int64(callback.UserId)
	messageId := callback.MessageId

	order, err := s.OrderBook.GetOrderByMessageId(ctx, callback.UserId, messageId)
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}
	hintMessage := tg.NewMessage(chatId, bot.ReplyEnterAmount)
	hint, err := s.Bot.Send(hintMessage)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	err = s.OrderBook.UpdateHintMessageForOrder(ctx, order.Id, uint64(hint.MessageID))
	if err != nil {
		return fmt.Errorf("failed to update hint message: %w", err)
	}

	err = s.OrderBook.UpdateOrderEditState(ctx, order.Id, types.EditStateWaitingRefundAmount)
	if err != nil {
		return fmt.Errorf("failed to update order state: %w", err)
	}

	return nil
}
