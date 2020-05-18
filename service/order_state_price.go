package service

import (
	"context"
	"fmt"
	tg "github.com/DarthRamone/telegram-bot-api"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
)

func (s Service) setWaitingForPrice(ctx context.Context, order types.Order) error {
	err := s.OrderBook.UpdateOrderEditState(ctx, order.Id, types.EditStateWaitingItemPrice)
	if err != nil {
		return fmt.Errorf("failed to change order state: %w", err)
	}

	if !order.HintMessageId.Valid {
		return fmt.Errorf("hint message is nil")
	}

	hintMessageId := int(order.HintMessageId.Int64)

	hintMessage := tg.NewEditMessageText(int64(order.UserId), hintMessageId, replyEnterItemPrice)
	_, err = s.Bot.Send(hintMessage)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
