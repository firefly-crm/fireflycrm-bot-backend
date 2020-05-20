package orderbook

import "context"

func (ob orderBook) UpdateMessageForOrder(ctx context.Context, userId, orderId, messageId uint64) error {
	return ob.storage.AddOrderMessage(ctx, userId, orderId, messageId)
}
