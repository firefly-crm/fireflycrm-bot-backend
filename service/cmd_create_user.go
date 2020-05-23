package service

import (
	"context"
	"fmt"
	. "github.com/firefly-crm/common/bot"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (s Service) createUser(ctx context.Context, userId uint64) error {
	err := s.Users.CreateUser(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	msg := tg.NewMessage(int64(userId), ReplyWelcome)
	msg.ReplyMarkup = merchantStandByKeyboardMarkup()

	_, err = s.Bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send welcome message: %w", err)
	}

	return nil
}
