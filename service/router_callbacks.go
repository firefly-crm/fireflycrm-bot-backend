package service

import (
	"context"
	"fmt"
	tg "github.com/DarthRamone/telegram-bot-api"
	. "github.com/firefly-crm/common/bot"
	"github.com/firefly-crm/common/logger"
	tp "github.com/firefly-crm/common/messages/telegram"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
)

func (s Service) ProcessCallbackEvent(ctx context.Context, callbackEvent *tp.CallbackEvent) (err error) {
	userId := callbackEvent.UserId
	messageId := callbackEvent.MessageId
	event := callbackEvent.Event

	var markup tg.InlineKeyboardMarkup

	log := logger.FromContext(ctx).
		WithField("user_id", userId).
		WithField("callbackEvent", tp.CallbackType_name[int32(event)]).
		WithField("message_id", messageId)

	ctx = logger.ToContext(ctx, log)

	log.Infof("processing callbackEvent")

	shouldDelete := false

	err = s.Storage.SetActiveOrderMessageForUser(ctx, userId, messageId)
	if err != nil {
		return fmt.Errorf("failed to set active order msg id: %w", err)
	}

	switch event {
	case tp.CallbackType_ITEMS:
		markup = orderItemsInlineKeyboard()
	case tp.CallbackType_BACK:
		markup, err = startOrderInlineKeyboard(ctx, s, userId, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	case tp.CallbackType_CANCEL:
		markup, err = startOrderInlineKeyboard(ctx, s, userId, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
		err = s.processCancelCallback(ctx, userId, messageId)
		if err != nil {
			return fmt.Errorf("failed to process cancel callbackEvent: %w", err)
		}
	case tp.CallbackType_RECEIPT_ITEMS_ADD:
		markup = cancelInlineKeyboard()
		err := s.processAddItemCallack(ctx, callbackEvent)
		if err != nil {
			return fmt.Errorf("failed to process add item callbackEvent: %w", err)
		}
	case tp.CallbackType_RECEIPT_ITEMS_REMOVE:
		markup, err = itemsListInlineKeyboard(ctx, s, userId, messageId, "remove")
		if err != nil {
			return fmt.Errorf("failed to get markup for remove items list: %w", err)
		}
	case tp.CallbackType_RECEIPT_ITEMS_EDIT:
		markup, err = itemsListInlineKeyboard(ctx, s, userId, messageId, "edit")
		if err != nil {
			return fmt.Errorf("failed to get markup for edit items list: %w", err)
		}
	case tp.CallbackType_CUSTOMER:
		markup, err = customerInlineKeyboard(ctx, s, userId, messageId)
		if err != nil {
			return fmt.Errorf("failed to get markup for customer action: %w", err)
		}
	case tp.CallbackType_PAYMENTS:
		markup, err = paymentInlineKeyboard(ctx, s, userId, messageId)
		if err != nil {
			return fmt.Errorf("failed to get payment inline markup: %w", err)
		}
	case tp.CallbackType_ADD_PAYMENT_TRANSFER:
		err = s.processAddPaymentCallback(ctx, userId, messageId, types.PaymentMethodCard2Card)
		if err != nil {
			return fmt.Errorf("failed to add card payment to order: %w", err)
		}
		markup, err = paymentAmountInlineKeyboard(ctx, s, userId, messageId)
		if err != nil {
			return fmt.Errorf("faield to get payment amount inline markup: %w", err)
		}
	case tp.CallbackType_ADD_PAYMENT_CASH:
		err = s.processAddPaymentCallback(ctx, userId, messageId, types.PaymentMethodCash)
		if err != nil {
			return fmt.Errorf("failed to add cash payment to order: %w", err)
		}
		markup, err = paymentAmountInlineKeyboard(ctx, s, userId, messageId)
		if err != nil {
			return fmt.Errorf("faield to get payment amount inline markup: %w", err)
		}
	case tp.CallbackType_ADD_PAYMENT_LINK:
		err = s.processAddPaymentCallback(ctx, userId, messageId, types.PaymentMethodAcquiring)
		if err != nil {
			return fmt.Errorf("failed to add link payment to order: %w", err)
		}
		markup, err = paymentAmountInlineKeyboard(ctx, s, userId, messageId)
		if err != nil {
			return fmt.Errorf("faield to get payment amount inline markup: %w", err)
		}
	case tp.CallbackType_PAYMENT_AMOUNT_FULL:
		err := s.processPaymentCallback(ctx, userId, messageId, 0)
		if err != nil {
			return fmt.Errorf("failed to process full payment callbackEvent: %w", err)
		}
		markup, err = startOrderInlineKeyboard(ctx, s, userId, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	case tp.CallbackType_PAYMENT_AMOUNT_PARTIAL:
		markup = cancelInlineKeyboard()
		err := s.processPartialPaymentCallback(ctx, callbackEvent)
		if err != nil {
			return fmt.Errorf("failed to process partial payment callbackEvent: %w", err)
		}
	case tp.CallbackType_PAYMENTS_REFUND:
		var err error
		markup, err = paymentsListInlineKeyboard(ctx, s, userId, messageId, "refund")
		if err != nil {
			return fmt.Errorf("failed to get payments list markup: %w", err)
		}
	case tp.CallbackType_PAYMENT_REFUND_PARTIAL:
		markup = cancelInlineKeyboard()
		err := s.processPartialRefundCallback(ctx, callbackEvent)
		if err != nil {
			return fmt.Errorf("failed to process partial refund callbackEvent: %w", err)
		}
	case tp.CallbackType_PAYMENT_REFUND_FULL:
		order, err := s.OrderBook.GetOrderByMessageId(ctx, userId, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order: %w", err)
		}
		err = s.processRefundCallback(ctx, order, userId, messageId, 0)
		if err != nil {
			return fmt.Errorf("failed to process refund callbackEvent: %w", err)
		}
		markup, err = startOrderInlineKeyboard(ctx, s, userId, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	case tp.CallbackType_PAYMENTS_REMOVE:
		markup, err = paymentsListInlineKeyboard(ctx, s, userId, messageId, "remove")
		if err != nil {
			return fmt.Errorf("failed to get payments list markup: %w", err)
		}
	case tp.CallbackType_ORDER_ACTIONS:
		markup, err = orderActionsInlineKeyboard(ctx, s, userId, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order actions markup: %w", err)
		}
	case tp.CallbackType_ORDER_STATE_DONE:
		err := s.processOrderStateCallback(ctx, userId, messageId, types.OrderStateDone)
		if err != nil {
			return fmt.Errorf("failed to process order done callbackEvent: %w", err)
		}
		markup, err = startOrderInlineKeyboard(ctx, s, userId, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	case tp.CallbackType_ORDER_RESTART:
		err := s.processOrderStateCallback(ctx, userId, messageId, types.OrderStateForming)
		if err != nil {
			return fmt.Errorf("failed to process order restart callbackEvent: %w", err)
		}
		markup, err = startOrderInlineKeyboard(ctx, s, userId, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	case tp.CallbackType_ORDER_DELETE:
		err := s.processOrderStateCallback(ctx, userId, messageId, types.OrderStateDeleted)
		if err != nil {
			return fmt.Errorf("failed to process order delete callbackEvent: %w", err)
		}
		markup = restoreDeletedOrderInlineKeyboard()
	case tp.CallbackType_ORDER_RESTORE:
		err := s.processOrderStateCallback(ctx, userId, messageId, types.OrderStateForming)
		if err != nil {
			return fmt.Errorf("failed to process order restore callbackEvent: %w", err)
		}
		markup, err = startOrderInlineKeyboard(ctx, s, userId, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	case tp.CallbackType_ORDER_STATE_IN_PROGRESS:
		err := s.processOrderStateCallback(ctx, userId, messageId, types.OrderStateInProgress)
		if err != nil {
			return fmt.Errorf("failed to process order restore callbackEvent: %w", err)
		}
		markup, err = startOrderInlineKeyboard(ctx, s, userId, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	case tp.CallbackType_ORDER_COLLAPSE:
		err := s.processOrderDisplayModeCallback(ctx, userId, messageId, types.DisplayModeCollapsed)
		if err != nil {
			return fmt.Errorf("failed to process order collapse callbackEvent: %w", err)
		}
		markup, err = startOrderInlineKeyboard(ctx, s, userId, messageId)
		if err != nil {
			return fmt.Errorf("failed to get start order inline keyboard: %w", err)
		}
	case tp.CallbackType_ORDER_EXPAND:
		err := s.processOrderDisplayModeCallback(ctx, userId, messageId, types.DisplayModeFull)
		if err != nil {
			return fmt.Errorf("failed to process order expand callbackEvent: %w", err)
		}
		markup, err = startOrderInlineKeyboard(ctx, s, userId, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	case tp.CallbackType_CUSTOM_ITEM_DELIVERY:
		err := s.processAddKnownItem(ctx, callbackEvent, KbDataDelivery)
		if err != nil {
			return fmt.Errorf("failed to process add item callbackEvent: %w", err)
		}
		markup = cancelInlineKeyboard()
	case tp.CallbackType_CUSTOM_ITEM_LINGERIE_SET:
		err := s.processAddKnownItem(ctx, callbackEvent, KbDataLingerieSet)
		if err != nil {
			return fmt.Errorf("failed to process add item callbackEvent: %w", err)
		}
		markup = cancelInlineKeyboard()
	case tp.CallbackType_NOTIFY_READ:
		shouldDelete = true
	case tp.CallbackType_CUSTOMER_EDIT_EMAIL:
		markup = cancelInlineKeyboard()
		err := s.processCustomerEditEmail(ctx, callbackEvent)
		if err != nil {
			return fmt.Errorf("failed to process customer edit name callbackEvent: %w", err)
		}
	case tp.CallbackType_CUSTOMER_EDIT_INSTAGRAM:
		markup = cancelInlineKeyboard()
		err := s.processCustomerEditInstagram(ctx, callbackEvent)
		if err != nil {
			return fmt.Errorf("failed to process customer edit instagram callbackEvent: %w", err)
		}
	case tp.CallbackType_CUSTOMER_EDIT_PHONE:
		markup = cancelInlineKeyboard()
		err := s.processCustomerEditPhone(ctx, callbackEvent)
		if err != nil {
			return fmt.Errorf("faield to process customer edit phone callbackEvent: %w", err)
		}
	case tp.CallbackType_CUSTOMER_EDIT_DESCRIPTION:
		markup = cancelInlineKeyboard()
		err := s.processCustomerEditDescription(ctx, callbackEvent)
		if err != nil {
			return fmt.Errorf("faield to process customer edit description callbackEvent: %w", err)
		}
	case tp.CallbackType_PAYMENT_REMOVE:
		err := s.processPaymentRemove(ctx, callbackEvent)
		if err != nil {
			return fmt.Errorf("failed to process remove payment callbackEvent: %w", err)
		}
		markup, err = startOrderInlineKeyboard(ctx, s, userId, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	case tp.CallbackType_PAYMENT_REFUND:
		order, err := s.OrderBook.GetOrderByMessageId(ctx, userId, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order: %w", err)
		}
		err = s.OrderBook.SetActivePaymentId(ctx, order.Id, callbackEvent.EntityId)
		if err != nil {
			return fmt.Errorf("failed to set active refund payment: %w", err)
		}
		markup = refundAmountInlineKeyboard()
	case tp.CallbackType_RECEIPT_ITEM_REMOVE:
		err := s.processItemRemove(ctx, callbackEvent)
		if err != nil {
			return fmt.Errorf("failed to remove item: %w", err)
		}
	case tp.CallbackType_RECEIPT_ITEM_EDIT:
		markup = editItemActionsInlineKeyboard(callbackEvent.EntityId)
	case tp.CallbackType_RECEIPT_ITEM_EDIT_QTY:
		markup = cancelInlineKeyboard()
		err = s.processItemEditQty(ctx, callbackEvent)
		if err != nil {
			return fmt.Errorf("failed to process item edit qty callbackEvent: %w", err)
		}
	case tp.CallbackType_RECEIPT_ITEM_EDIT_PRICE:
		markup = cancelInlineKeyboard()
		err = s.processItemEditPrice(ctx, callbackEvent)
		if err != nil {
			return fmt.Errorf("failed to process item edit price callbackEvent: %w", err)
		}
	case tp.CallbackType_RECEIPT_ITEM_EDIT_NAME:
		markup = cancelInlineKeyboard()
		err = s.processItemEditName(ctx, callbackEvent)
		if err != nil {
			return fmt.Errorf("failed to process item edit name callbackEvent: %w", err)
		}
	case tp.CallbackType_ORDER_EDIT:
		markup, err = orderEditEntriesInlineKeyboard(ctx, s, userId, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order edit kb markup: %w", err)
		}
	case tp.CallbackType_ORDER_EDIT_DUE_DATE:
		markup = cancelInlineKeyboard()
		err := s.processOrderEditDueDate(ctx, callbackEvent)
		if err != nil {
			return fmt.Errorf("failed to process edit order due date callback: %w", err)
		}
	case tp.CallbackType_ORDER_EDIT_DESCRIPTION:
		markup = cancelInlineKeyboard()
		err := s.processOrderEditDescription(ctx, callbackEvent)
		if err != nil {
			return fmt.Errorf("failed to edit order description: %w", err)
		}
	default:
		return fmt.Errorf("unknown callbackEvent event: %v", event)
	}

	if !shouldDelete {
		err := s.updateOrderMessage(ctx, userId, messageId, &markup)
		if err != nil {
			log.Errorf("failed to update order message: %v", err)
		}
	} else {
		msg := tg.NewDeleteMessage(int64(userId), int(messageId))
		_, err = s.Bot.Send(msg)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
	}

	return nil
}
