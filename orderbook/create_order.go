package orderbook

import (
	"context"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
)

func (b orderBook) CreateOrder(ctx context.Context, userId uint64) (types.Order, error) {
	return b.storage.CreateOrder(ctx, userId)
}
