package service

import (
	"context"
	"fmt"
	"github.com/firefly-crm/common/logger"
	tp "github.com/firefly-crm/common/messages/telegram"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (s Service) processCallback(ctx context.Context, callback *tp.CallbackEvent) (err error) {
	userId := callback.UserId
	messageId := callback.MessageId
	event := callback.Event

	var markup tg.InlineKeyboardMarkup

	log := logger.FromContext(ctx).
		WithField("user_id", userId).
		WithField("callback", tp.CallbackType_name[int32(event)]).
		WithField("message_id", messageId)

	ctx = logger.ToContext(ctx, log)

	log.Infof("processing callback")

	shouldDelete := false

	err = s.Storage.SetActiveOrderMessageForUser(ctx, userId, messageId)
	if err != nil {
		return fmt.Errorf("failed to set active order msg id: %w", err)
	}

	switch event {
	case tp.CallbackType_ITEMS:
		markup = orderItemsInlineKeyboard()
	case tp.CallbackType_BACK:
		markup, err = startOrderInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	case tp.CallbackType_CANCEL:
		markup, err = startOrderInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
		err = s.processCancelCallback(ctx, messageId)
		if err != nil {
			return fmt.Errorf("failed to process cancel callback: %w", err)
		}
	case tp.CallbackType_RECEIPT_ITEMS_ADD:
		markup = cancelInlineKeyboard()
		err := s.processAddItemCallack(ctx, callback)
		if err != nil {
			return fmt.Errorf("failed to process add item callback: %w", err)
		}
	case tp.CallbackType_RECEIPT_ITEMS_REMOVE:
		markup, err = itemsListInlineKeyboard(ctx, s, messageId, "remove")
		if err != nil {
			return fmt.Errorf("failed to get markup for remove items list: %w", err)
		}
	case tp.CallbackType_RECEIPT_ITEMS_EDIT:
		markup, err = itemsListInlineKeyboard(ctx, s, messageId, "edit")
		if err != nil {
			return fmt.Errorf("failed to get markup for edit items list: %w", err)
		}
	case tp.CallbackType_CUSTOMER:
		markup, err = customerInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get markup for customer action: %w", err)
		}
	case tp.CallbackType_PAYMENTS:
		markup, err = paymentInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get payment inline markup: %w", err)
		}
	case tp.CallbackType_ADD_PAYMENT_TRANSFER:
		err = s.processAddPaymentCallback(ctx, messageId, types.PaymentMethodCard2Card)
		if err != nil {
			return fmt.Errorf("failed to add card payment to order: %w", err)
		}
		markup, err = paymentAmountInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("faield to get payment amount inline markup: %w", err)
		}
	case tp.CallbackType_ADD_PAYMENT_CASH:
		err = s.processAddPaymentCallback(ctx, messageId, types.PaymentMethodCash)
		if err != nil {
			return fmt.Errorf("failed to add cash payment to order: %w", err)
		}
		markup, err = paymentAmountInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("faield to get payment amount inline markup: %w", err)
		}
	case tp.CallbackType_ADD_PAYMENT_LINK:
		err = s.processAddPaymentCallback(ctx, messageId, types.PaymentMethodAcquiring)
		if err != nil {
			return fmt.Errorf("failed to add link payment to order: %w", err)
		}
		markup, err = paymentAmountInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("faield to get payment amount inline markup: %w", err)
		}
	case tp.CallbackType_PAYMENT_AMOUNT_FULL:
		err := s.processPaymentCallback(ctx, messageId, 0)
		if err != nil {
			return fmt.Errorf("failed to process full payment callback: %w", err)
		}
		markup, err = startOrderInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	case tp.CallbackType_PAYMENT_AMOUNT_PARTIAL:
		markup = cancelInlineKeyboard()
		err := s.processPartialPaymentCallback(ctx, callback)
		if err != nil {
			return fmt.Errorf("failed to process partial payment callback: %w", err)
		}
	case tp.CallbackType_PAYMENTS_REFUND:
		var err error
		markup, err = paymentsListInlineKeyboard(ctx, s, messageId, "refund")
		if err != nil {
			return fmt.Errorf("failed to get payments list markup: %w", err)
		}
	case tp.CallbackType_PAYMENT_REFUND_PARTIAL:
		markup = cancelInlineKeyboard()
		err := s.processPartialRefundCallback(ctx, callback)
		if err != nil {
			return fmt.Errorf("failed to process partial refund callback: %w", err)
		}
	case tp.CallbackType_PAYMENT_REFUND_FULL:
		order, err := s.OrderBook.GetOrderByMessageId(ctx, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order: %w", err)
		}
		err = s.processRefundCallback(ctx, order, messageId, 0)
		if err != nil {
			return fmt.Errorf("failed to process refund callback: %w", err)
		}
		markup, err = startOrderInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	case tp.CallbackType_PAYMENTS_REMOVE:
		markup, err = paymentsListInlineKeyboard(ctx, s, messageId, "remove")
		if err != nil {
			return fmt.Errorf("failed to get payments list markup: %w", err)
		}
	case tp.CallbackType_ORDER_ACTIONS:
		markup, err = orderActionsInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order actions markup: %w", err)
		}
	case tp.CallbackType_ORDER_STATE_DONE:
		err := s.processOrderStateCallback(ctx, messageId, types.OrderStateDone)
		if err != nil {
			return fmt.Errorf("failed to process order done callback: %w", err)
		}
		markup, err = startOrderInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	case tp.CallbackType_ORDER_RESTART:
		err := s.processOrderStateCallback(ctx, messageId, types.OrderStateForming)
		if err != nil {
			return fmt.Errorf("failed to process order restart callback: %w", err)
		}
		markup, err = startOrderInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	case tp.CallbackType_ORDER_DELETE:
		err := s.processOrderStateCallback(ctx, messageId, types.OrderStateDeleted)
		if err != nil {
			return fmt.Errorf("failed to process order delete callback: %w", err)
		}
		markup = restoreDeletedOrderInlineKeyboard()
	case tp.CallbackType_ORDER_RESTORE:
		err := s.processOrderStateCallback(ctx, messageId, types.OrderStateForming)
		if err != nil {
			return fmt.Errorf("failed to process order restore callback: %w", err)
		}
		markup, err = startOrderInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	case tp.CallbackType_ORDER_STATE_IN_PROGRESS:
		err := s.processOrderStateCallback(ctx, messageId, types.OrderStateInProgress)
		if err != nil {
			return fmt.Errorf("failed to process order restore callback: %w", err)
		}
		markup, err = startOrderInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	case tp.CallbackType_ORDER_COLLAPSE:
		err := s.processOrderDisplayModeCallback(ctx, messageId, types.DisplayModeCollapsed)
		if err != nil {
			return fmt.Errorf("failed to process order collapse callback: %w", err)
		}
		markup = expandOrderInlineKeyboard()
	case tp.CallbackType_ORDER_EXPAND:
		err := s.processOrderDisplayModeCallback(ctx, messageId, types.DisplayModeFull)
		if err != nil {
			return fmt.Errorf("failed to process order expand callback: %w", err)
		}
		markup, err = startOrderInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	case tp.CallbackType_CUSTOM_ITEM_DELIVERY:
		err := s.processAddKnownItem(ctx, callback, kbDataDelivery)
		if err != nil {
			return fmt.Errorf("failed to process add item callback: %w", err)
		}
		markup = cancelInlineKeyboard()
	case tp.CallbackType_CUSTOM_ITEM_LINGERIE_SET:
		err := s.processAddKnownItem(ctx, callback, kbDataLingerieSet)
		if err != nil {
			return fmt.Errorf("failed to process add item callback: %w", err)
		}
		markup = cancelInlineKeyboard()
	case tp.CallbackType_NOTIFY_READ:
		shouldDelete = true
	case tp.CallbackType_CUSTOMER_EDIT_EMAIL:
		markup = cancelInlineKeyboard()
		err := s.processCustomerEditEmail(ctx, callback)
		if err != nil {
			return fmt.Errorf("failed to process customer edit name callback: %w", err)
		}
	case tp.CallbackType_CUSTOMER_EDIT_INSTAGRAM:
		markup = cancelInlineKeyboard()
		err := s.processCustomerEditInstagram(ctx, callback)
		if err != nil {
			return fmt.Errorf("failed to process customer edit instagram callback: %w", err)
		}
	case tp.CallbackType_CUSTOMER_EDIT_PHONE:
		markup = cancelInlineKeyboard()
		err := s.processCustomerEditPhone(ctx, callback)
		if err != nil {
			return fmt.Errorf("faield to process customer edit phone callback: %w", err)
		}
	case tp.CallbackType_PAYMENT_REMOVE:
		err := s.processPaymentRemove(ctx, callback)
		if err != nil {
			return fmt.Errorf("failed to process remove payment callback: %w", err)
		}
		markup, err = startOrderInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	case tp.CallbackType_PAYMENT_REFUND:
		order, err := s.OrderBook.GetOrderByMessageId(ctx, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order: %w", err)
		}
		err = s.OrderBook.SetActivePaymentId(ctx, order.Id, callback.EntityId)
		if err != nil {
			return fmt.Errorf("failed to set active refund payment: %w", err)
		}
		markup = refundAmountInlineKeyboard()
	case tp.CallbackType_RECEIPT_ITEM_REMOVE:
		err := s.processItemRemove(ctx, callback)
		if err != nil {
			return fmt.Errorf("failed to remove item: %w", err)
		}
	case tp.CallbackType_RECEIPT_ITEM_EDIT:
		markup = editItemActionsInlineKeyboard(callback.EntityId)
	case tp.CallbackType_RECEIPT_ITEM_EDIT_QTY:
		markup = cancelInlineKeyboard()
		err = s.processItemEditQty(ctx, callback)
		if err != nil {
			return fmt.Errorf("failed to process item edit qty callback: %w", err)
		}
	case tp.CallbackType_RECEIPT_ITEM_EDIT_PRICE:
		markup = cancelInlineKeyboard()
		err = s.processItemEditPrice(ctx, callback)
		if err != nil {
			return fmt.Errorf("failed to process item edit price callback: %w", err)
		}
	case tp.CallbackType_RECEIPT_ITEM_EDIT_NAME:
		markup = cancelInlineKeyboard()
		err = s.processItemEditName(ctx, callback)
		if err != nil {
			return fmt.Errorf("failed to process item edit name callback: %w", err)
		}
	default:
		return fmt.Errorf("unknown callback event: %v", event)
	}

	var msg tg.Chattable
	if !shouldDelete {
		msg = tg.NewEditMessageReplyMarkup(int64(userId), int(messageId), markup)
	} else {
		msg = tg.NewDeleteMessage(int64(userId), int(messageId))
	}
	_, err = s.Bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
