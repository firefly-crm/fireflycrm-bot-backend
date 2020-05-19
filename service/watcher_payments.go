package service

import (
	"context"
	"fmt"
	mb "github.com/DarthRamone/modulbank-go"
	tg "github.com/DarthRamone/telegram-bot-api"
	"github.com/firefly-crm/common/logger"
	"net/http"
	"time"
)

func (s Service) StartPaymentsWatcher(ctx context.Context, interval time.Duration) error {
	log := logger.FromContext(ctx)

	ticker := time.NewTicker(interval)
	defer func() {
		ticker.Stop()
	}()

outsideLoop:
	for {
		select {
		case <-ctx.Done():
			break outsideLoop
		case <-ticker.C:
			err := s.checkPayments(ctx)
			if err != nil {
				log.Errorf("failed to check payments: %v", err.Error())
			}
		}
	}

	return nil
}

func (s Service) checkPayments(ctx context.Context) error {
	log := logger.FromContext(ctx)

	payments, err := s.OrderBook.GetBankPayments(ctx)
	if err != nil {
		return fmt.Errorf("failed to get payments: %w", err)
	}

	for _, p := range payments {
		order, err := s.OrderBook.GetOrder(ctx, p.OrderId)
		if err != nil {
			log.Errorf("failed to get order: %w", err)
			continue
		}

		user, err := s.Storage.GetUser(ctx, order.UserId)
		if err != nil {
			log.Errorf("failed to get user: %w", err)
			continue
		}

		opts := mb.MerchantOptions{
			Merchant:  user.MerchantId,
			SecretKey: user.SecretKey,
		}

		log.Debugf("bank payment id: %s", p.BankPaymentId)

		bill, err := mb.GetBill(ctx, p.BankPaymentId, opts, http.DefaultClient)
		if err != nil {
			log.Errorf("failed to get bill: %w", err)
			continue
		}

		if bill.Paid == 1 {
			err := s.Storage.SetPaymentPaid(ctx, p.Id)
			if err != nil {
				log.Errorf("failed to set payment paid: %w", err)
				continue
			}

			messages, err := s.Storage.GetMessagesForOrder(ctx, user.Id, order.Id)
			if err != nil {
				log.Errorf("failed to get messages for order: %w", err)
				continue
			}

			if len(messages) == 0 {
				log.Errorf("no messages for order found")
				continue
			}

			msg := tg.NewMessage(int64(user.Id), "Заказ оплачен")
			msg.ReplyToMessageID = int(messages[0].Id)
			msg.ReplyMarkup = notifyReadInlineKeyboard()
			_, err = s.Bot.Send(msg)
			if err != nil {
				log.Errorf("failed to send message to chat: %v", err)
				continue
			}

			err = s.updateOrderMessage(ctx, user.Id, messages[0].Id, true)
			if err != nil {
				log.Errorf("failed to update order message: %v", err)
				continue
			}
		} else if bill.Expired == 1 {
			err := s.Storage.SetPaymentExpired(ctx, p.Id)
			if err != nil {
				log.Errorf("failed to set payment expired: %w", err)
				continue
			}
		}
	}

	return nil
}
