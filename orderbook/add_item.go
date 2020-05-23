package orderbook

import (
	"context"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
)

//Adds receipt item to existing order and returns it's unique id
func (b orderBook) AddItem(ctx context.Context, orderId uint64, t types.ReceiptItemType) (uint64, error) {
	return b.storage.AddItemToOrder(ctx, orderId, t)
}
