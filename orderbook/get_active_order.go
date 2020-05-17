package orderbook

import (
	"context"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
)

func (ob orderBook) GetActiveOrderForUser(ctx context.Context, userId uint64) (types.Order, error) {
	return ob.storage.GetActiveOrderForUser(ctx, userId)
}
