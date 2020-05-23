package service

import (
	"context"
	"fmt"
	tg "github.com/DarthRamone/telegram-bot-api"
	. "github.com/firefly-crm/common/bot"
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
