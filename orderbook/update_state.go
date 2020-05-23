package orderbook

import (
	"context"
	"github.com/firefly-crm/common/bot"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
)

func (ob orderBook) UpdateOrderEditState(ctx context.Context, orderId uint64, state types.EditState) error {
	hint := ""
	switch state {
	case types.EditStateWaitingItemPrice:
		hint = bot.ReplyEnterItemPrice
	case types.EditStateWaitingItemName:
		hint = bot.ReplyEnterItemName
	case types.EditStateWaitingPaymentAmount:
		hint = bot.ReplyEnterAmount
	case types.EditStateWaitingCustomerInstagram:
		hint = bot.ReplyEnterCustomerInstagram
	case types.EditStateWaitingCustomerEmail:
		hint = bot.ReplyEnterCustomerEmail
	case types.EditStateWaitingCustomerPhone:
		hint = bot.ReplyEnterCustomerPhone
	case types.EditStateWaitingCustomerDescription:
		hint = bot.ReplyEnterCustomerDescription
	case types.EditStateWaitingItemQuantity:
		hint = bot.ReplyEnterItemQty
	case types.EditStateWaitingOrderDueDate:
		hint = bot.ReplyEnterOrderDueDate
	case types.EditStateWaitingOrderDescription:
		hint = bot.ReplyEnterOrderDescription
	case types.EditStateWaitingRefundAmount:
		hint = bot.ReplyEnterAmount
	}

	return ob.storage.UpdateOrderEditState(ctx, orderId, state, hint)
}

func (ob orderBook) UpdateOrderState(ctx context.Context, orderId uint64, state types.OrderState) error {
	return ob.storage.UpdateOrderState(ctx, orderId, state)
}
