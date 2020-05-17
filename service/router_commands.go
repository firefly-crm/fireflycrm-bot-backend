package service

import (
	"context"
	"fmt"
	"github.com/firefly-crm/common/logger"
	tp "github.com/firefly-crm/common/messages/telegram"
)

func (s Service) ProcessCommandEvent(ctx context.Context, commandEvent *tp.CommandEvent) error {
	var err error

	log := logger.
		FromContext(ctx).
		WithField("user_id", commandEvent.UserId).
		WithField("command", tp.CommandType_name[int32(commandEvent.Command)])

	log.Infof("processing command")

	ctx = logger.ToContext(ctx, log)

	switch commandEvent.Command {
	case tp.CommandType_START:
		err = s.createUser(ctx, commandEvent.UserId)
	case tp.CommandType_CREATE_ORDER:
		err = s.createOrder(ctx, commandEvent.UserId, commandEvent.MessageId)
	case tp.CommandType_REGISTER_AS_MERCHANT:
		err = s.registerMerchant(ctx, commandEvent)
	default:
		return fmt.Errorf("unknown command: %s", tp.CommandType_name[int32(commandEvent.Command)])
	}

	if err != nil {
		return fmt.Errorf("failed process message: %w", err)
	}

	return nil
}
