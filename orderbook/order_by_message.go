package orderbook

import (
	"context"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
)

func (o orderBook) GetOrderByMessageId(ctx context.Context, userId, messageId uint64) (order types.Order, err error) {
	return o.storage.GetOrderByMessageId(ctx, userId, messageId)
}
