package service

import (
	"context"
	"fmt"
	"github.com/firefly-crm/common/bot"
	tp "github.com/firefly-crm/common/messages/telegram"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (s Service) registerMerchant(ctx context.Context, commandEvent *tp.CommandEvent) error {
	userId := commandEvent.UserId

	if len(commandEvent.Arguments) != 2 {
		return fmt.Errorf("wrong arguments count")
	}

	merchantId := commandEvent.Arguments[0]
	secretKey := commandEvent.Arguments[1]

	err := s.Users.RegisterAsMerchant(ctx, userId, merchantId, secretKey)
	if err != nil {
		return fmt.Errorf("failed to register as merchant: %w", err)
	}

	deleteMessage := tg.NewDeleteMessage(int64(userId), int(commandEvent.MessageId))
	_, err = s.Bot.DeleteMessage(deleteMessage)
	if err != nil {
		return fmt.Errorf("failed to delete command message: %w", err)
	}

	msg := tg.NewMessage(int64(userId), bot.ReplyMerchantSuccessfulRegistered)
	msg.ReplyMarkup = merchantStandByKeyboardMarkup()
	msg.ParseMode = "html"

	_, err = s.Bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
