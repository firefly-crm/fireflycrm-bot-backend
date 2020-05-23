package service

import (
	"context"
	"fmt"
	"github.com/firefly-crm/common/bot"
	"github.com/firefly-crm/fireflycrm-bot-backend/types"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"sort"
)

func merchantStandByKeyboardMarkup() tg.ReplyKeyboardMarkup {
	createOrderButton := tg.NewKeyboardButton(bot.KbCreateOrder)
	getOrdersButton := tg.NewKeyboardButton(bot.KbActiveOrders)
	rows := [][]tg.KeyboardButton{
		{createOrderButton},
		{getOrdersButton},
	}
	return tg.NewReplyKeyboard(rows...)
}

func startOrderInlineKeyboard(ctx context.Context, s Service, userId, messageId uint64) (tg.InlineKeyboardMarkup, error) {
	var markup tg.InlineKeyboardMarkup

	order, err := s.OrderBook.GetOrderByMessageId(ctx, userId, messageId)
	if err != nil {
		return markup, fmt.Errorf("failed to get order for markup: %w", err)
	}

	message, err := s.Storage.GetOrderMessage(ctx, userId, messageId)
	if err != nil {
		return markup, fmt.Errorf("failed to get order message: %w", err)
	}

	if message.DisplayMode == types.DisplayModeCollapsed {
		return startOrderCollapsedInlineKeyboard(), nil
	}

	collapseButton := tg.NewInlineKeyboardButtonData(bot.KbOrderCollapsePictogram, bot.KbDataOrderCollapse)
	customerButton := tg.NewInlineKeyboardButtonData(bot.KbCustomerPictogram, bot.KbDataCustomer)
	paymentButton := tg.NewInlineKeyboardButtonData(bot.KbPaymentPictogram, bot.KbDataPayment)
	actionsButton := tg.NewInlineKeyboardButtonData(bot.KbOrderActionsPictogram, bot.KbDataOrderActions)
	itemsButton := tg.NewInlineKeyboardButtonData(bot.KbItemsPictogram, bot.KbDataItems)

	if order.OrderState != types.OrderStateDone {
		var row1 []tg.InlineKeyboardButton
		if !order.CustomerId.Valid {
			row1 = []tg.InlineKeyboardButton{itemsButton, customerButton, actionsButton, collapseButton}
		} else {
			row1 = []tg.InlineKeyboardButton{itemsButton, customerButton, paymentButton, actionsButton, collapseButton}
		}
		markup = tg.NewInlineKeyboardMarkup(row1)
	} else {
		row1 := []tg.InlineKeyboardButton{customerButton, paymentButton, actionsButton, collapseButton}
		markup = tg.NewInlineKeyboardMarkup(row1)
	}

	return markup, nil
}
func startOrderCollapsedInlineKeyboard() tg.InlineKeyboardMarkup {
	var markup tg.InlineKeyboardMarkup

	expand := tg.NewInlineKeyboardButtonData(bot.KbOrderExpandPictogram, bot.KbDataOrderExpand)
	forming := tg.NewInlineKeyboardButtonData(bot.KbOrderRestorePictogram, bot.KbDataOrderRestore)
	inProgress := tg.NewInlineKeyboardButtonData(bot.KbOrderInProgressPictogram, bot.KbDataOrderInProgress)
	done := tg.NewInlineKeyboardButtonData(bot.KbOrderDonePictogram, bot.KbDataOrderDone)

	actionsButton := tg.NewInlineKeyboardButtonData(bot.KbOrderActionsPictogram, bot.KbDataOrderActions)

	row1 := []tg.InlineKeyboardButton{forming, inProgress, done, actionsButton, expand}
	markup = tg.NewInlineKeyboardMarkup(row1)

	return markup
}

func restoreDeletedOrderInlineKeyboard() tg.InlineKeyboardMarkup {
	restoreButton := tg.NewInlineKeyboardButtonData(bot.KbOrderRestore, bot.KbDataOrderRestore)
	return tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(restoreButton))
}

func orderEditEntriesInlineKeyboard(ctx context.Context, s Service, userId, messageId uint64) (tg.InlineKeyboardMarkup, error) {
	var markup tg.InlineKeyboardMarkup

	order, err := s.OrderBook.GetOrderByMessageId(ctx, userId, messageId)
	if err != nil {
		return markup, fmt.Errorf("failed to get order for markup: %w", err)
	}

	dueDataData := fmt.Sprintf("%s_%d", bot.KbDataOrderEditDate, order.Id)
	descriptionData := fmt.Sprintf("%s_%d", bot.KbDataOrderEditDescription, order.Id)

	editDateButton := tg.NewInlineKeyboardButtonData(bot.KbOrderEditDueDate, dueDataData)
	editDescription := tg.NewInlineKeyboardButtonData(bot.KbOrderEditDescription, descriptionData)
	backButton := tg.NewInlineKeyboardButtonData(bot.KbBack, bot.KbDataBack)

	markup = tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(editDateButton),
		tg.NewInlineKeyboardRow(editDescription),
		tg.NewInlineKeyboardRow(backButton))
	return markup, nil
}

func orderActionsInlineKeyboard(ctx context.Context, s Service, userId, messageId uint64) (tg.InlineKeyboardMarkup, error) {
	var markup tg.InlineKeyboardMarkup
	var rows [][]tg.InlineKeyboardButton

	order, err := s.OrderBook.GetOrderByMessageId(ctx, userId, messageId)
	if err != nil {
		return markup, fmt.Errorf("failed to get order for markup: %w", err)
	}

	doneButton := tg.NewInlineKeyboardButtonData(bot.KbOrderDone, bot.KbDataOrderDone)
	inProgressButton := tg.NewInlineKeyboardButtonData(bot.KbOrderInProgress, bot.KbDataOrderInProgress)
	restartButton := tg.NewInlineKeyboardButtonData(bot.KbOrderRestart, bot.KbDataOrderRestart)
	editButton := tg.NewInlineKeyboardButtonData(bot.KbOrderEdit, bot.KbDataOrderEdit)
	deleteButton := tg.NewInlineKeyboardButtonData(bot.KbOrderDelete, bot.KbDataOrderDelete)
	backButton := tg.NewInlineKeyboardButtonData(bot.KbBack, bot.KbDataBack)

	if order.OrderState == types.OrderStateForming && order.Amount > 0 {
		row := []tg.InlineKeyboardButton{inProgressButton}
		if order.PayedAmount == order.Amount {
			row = append(row, doneButton)
		}
		rows = append(rows, row)
	}

	if order.OrderState == types.OrderStateInProgress {
		if order.PayedAmount == order.Amount {
			rows = append(rows, []tg.InlineKeyboardButton{doneButton})
		}
	}

	if order.OrderState == types.OrderStateDone {
		rows = append(rows, []tg.InlineKeyboardButton{restartButton})
	}

	rows = append(rows, []tg.InlineKeyboardButton{editButton, deleteButton})
	rows = append(rows, []tg.InlineKeyboardButton{backButton})

	return tg.NewInlineKeyboardMarkup(rows...), nil
}

func orderItemsInlineKeyboard() tg.InlineKeyboardMarkup {
	addItemButton := tg.NewInlineKeyboardButtonData(bot.KbAddItemPictogram, bot.KbDataAddItem)
	editItemButton := tg.NewInlineKeyboardButtonData(bot.KbEditItemPictogram, bot.KbDataEditItem)
	removeItemButton := tg.NewInlineKeyboardButtonData(bot.KbRemovePictogram, bot.KbDataRemoveItem)
	backButton := tg.NewInlineKeyboardButtonData(bot.KbBack, bot.KbDataBack)
	deliveryButton := tg.NewInlineKeyboardButtonData(bot.KbDeliveryPictogram, bot.KbDataDelivery)
	lingerieButton := tg.NewInlineKeyboardButtonData(bot.KbLingerieSetPictogram, bot.KbDataLingerieSet)
	row1 := []tg.InlineKeyboardButton{addItemButton, deliveryButton, lingerieButton, editItemButton, removeItemButton}
	row2 := []tg.InlineKeyboardButton{backButton}
	return tg.NewInlineKeyboardMarkup(row1, row2)
}

func customerInlineKeyboard(ctx context.Context, s Service, userId, messageId uint64) (tg.InlineKeyboardMarkup, error) {
	var markup tg.InlineKeyboardMarkup

	order, err := s.OrderBook.GetOrderByMessageId(ctx, userId, messageId)
	if err != nil {
		return markup, fmt.Errorf("failed to get order for markup: %w", err)
	}

	//nameButton := tg.NewInlineKeyboardButtonData(bot.KbName, fmt.Sprintf("customer_edit_name_%d", order.Id))
	emailButton := tg.NewInlineKeyboardButtonData(bot.KbCustomerEmail, fmt.Sprintf("customer_edit_email_%d", order.Id))
	instaButton := tg.NewInlineKeyboardButtonData(bot.KbCustomerInstagram, fmt.Sprintf("customer_edit_instagram_%d", order.Id))
	phoneButton := tg.NewInlineKeyboardButtonData(bot.KbCustomerPhone, fmt.Sprintf("customer_edit_phone_%d", order.Id))
	noteButton := tg.NewInlineKeyboardButtonData(bot.KbCustomerDescription, fmt.Sprintf("customer_edit_description_%d", order.Id))
	backButton := tg.NewInlineKeyboardButtonData(bot.KbBack, bot.KbDataBack)

	row2 := tg.NewInlineKeyboardRow(phoneButton)
	if order.CustomerId.Valid {
		row2 = append(row2, noteButton)
	}

	markups := [][]tg.InlineKeyboardButton{
		{instaButton, emailButton},
		row2,
		{backButton},
	}

	markup = tg.NewInlineKeyboardMarkup(markups...)
	return markup, nil
}

func itemsListInlineKeyboard(ctx context.Context, s Service, userId, messageId uint64, action string) (tg.InlineKeyboardMarkup, error) {
	var markup tg.InlineKeyboardMarkup

	order, err := s.OrderBook.GetOrderByMessageId(ctx, userId, messageId)
	if err != nil {
		return markup, fmt.Errorf("failed to get items markup: %w", err)
	}

	markups := make([][]tg.InlineKeyboardButton, 0)
	for _, i := range order.ReceiptItems {
		button := tg.NewInlineKeyboardButtonData(i.Name, fmt.Sprintf("item_%s_%d", action, i.Id))
		markups = append(markups, []tg.InlineKeyboardButton{button})
	}
	backButton := tg.NewInlineKeyboardButtonData(bot.KbBack, bot.KbDataBack)
	markups = append(markups, []tg.InlineKeyboardButton{backButton})

	return tg.NewInlineKeyboardMarkup(markups...), nil
}

func paymentsListInlineKeyboard(ctx context.Context, s Service, userId, messageId uint64, action string) (tg.InlineKeyboardMarkup, error) {
	var markup tg.InlineKeyboardMarkup

	order, err := s.OrderBook.GetOrderByMessageId(ctx, userId, messageId)
	if err != nil {
		return markup, fmt.Errorf("failed to get payments markup: %w", err)
	}

	sort.Sort(types.PaymentsByCreatedAt(order.Payments))

	markups := make([][]tg.InlineKeyboardButton, 0)
	for i, p := range order.Payments {
		if p.PaymentMethod == types.PaymentMethodAcquiring && p.Payed {
			continue
		}

		if action == "refund" && !p.Payed {
			continue
		}

		if action == "refund" && p.RefundAmount == p.Amount {
			continue
		}

		name := fmt.Sprintf("Платеж #%d", i+1)
		button := tg.NewInlineKeyboardButtonData(name, fmt.Sprintf("payment_%s_%d", action, p.Id))
		markups = append(markups, []tg.InlineKeyboardButton{button})
	}
	backButton := tg.NewInlineKeyboardButtonData(bot.KbBack, bot.KbDataBack)
	markups = append(markups, []tg.InlineKeyboardButton{backButton})

	return tg.NewInlineKeyboardMarkup(markups...), nil
}

func editItemActionsInlineKeyboard(itemId uint64) tg.InlineKeyboardMarkup {
	nameButton := tg.NewInlineKeyboardButtonData(bot.KbName, fmt.Sprintf("item_edit_name_%d", itemId))
	qtyButton := tg.NewInlineKeyboardButtonData(bot.KbQty, fmt.Sprintf("item_edit_qty_%d", itemId))
	priceButton := tg.NewInlineKeyboardButtonData(bot.KbPrice, fmt.Sprintf("item_edit_price_%d", itemId))
	row := []tg.InlineKeyboardButton{nameButton, qtyButton, priceButton}
	backButton := tg.NewInlineKeyboardButtonData(bot.KbBack, bot.KbDataBack)
	row1 := []tg.InlineKeyboardButton{backButton}
	return tg.NewInlineKeyboardMarkup(row, row1)
}

func cancelInlineKeyboard() tg.InlineKeyboardMarkup {
	cancelButton := tg.NewInlineKeyboardButtonData(bot.KbCancel, bot.KbDataCancel)
	row1 := []tg.InlineKeyboardButton{cancelButton}
	return tg.NewInlineKeyboardMarkup(row1)
}

func paymentInlineKeyboard(ctx context.Context, s Service, userId, messageId uint64) (tg.InlineKeyboardMarkup, error) {
	var markup tg.InlineKeyboardMarkup

	order, err := s.OrderBook.GetOrderByMessageId(ctx, userId, messageId)
	if err != nil {
		return markup, fmt.Errorf("failed to get payments markup: %w", err)
	}

	customer, err := s.Users.GetCustomer(ctx, uint64(order.CustomerId.Int64))
	if err != nil {
		return markup, fmt.Errorf("failed to get customer: %w", err)
	}

	linkButton := tg.NewInlineKeyboardButtonData(bot.KbPaymentLink, bot.KbDataPaymentLink)
	cardButton := tg.NewInlineKeyboardButtonData(bot.KbPaymentCard, bot.KbDataPaymentCard)
	cashButton := tg.NewInlineKeyboardButtonData(bot.KbPaymentCash, bot.KbDataPaymentCash)
	deleteButton := tg.NewInlineKeyboardButtonData(bot.KbRemove, bot.KbDataRemovePayment)
	refundButton := tg.NewInlineKeyboardButtonData(bot.KbRefundPayment, bot.KbDataRefundPayment)
	cancelButton := tg.NewInlineKeyboardButtonData(bot.KbCancel, bot.KbDataCancel)

	if order.Amount == 0 {
		return tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(cancelButton)), nil
	}

	if order.PayedAmount >= order.Amount {
		var row1 []tg.InlineKeyboardButton
		if order.OrderState == types.OrderStateDone {
			row1 = []tg.InlineKeyboardButton{refundButton}
		} else {
			row1 = []tg.InlineKeyboardButton{deleteButton, refundButton}
		}
		row2 := []tg.InlineKeyboardButton{cancelButton}
		markup = tg.NewInlineKeyboardMarkup(row1, row2)
	} else {
		var row1 []tg.InlineKeyboardButton
		if order.PayedAmount == 0 && customer.Email.Valid {
			row1 = []tg.InlineKeyboardButton{linkButton, cardButton, cashButton}
		} else {
			row1 = []tg.InlineKeyboardButton{cardButton, cashButton}
		}
		row2 := []tg.InlineKeyboardButton{deleteButton, refundButton}
		row3 := []tg.InlineKeyboardButton{cancelButton}
		markup = tg.NewInlineKeyboardMarkup(row1, row2, row3)
	}

	return markup, nil
}

func notifyReadInlineKeyboard() tg.InlineKeyboardMarkup {
	okButton := tg.NewInlineKeyboardButtonData(bot.KbOk, bot.KbDataNotifyRead)
	return tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(okButton))
}

func paymentAmountInlineKeyboard(ctx context.Context, s Service, userId, messageId uint64) (tg.InlineKeyboardMarkup, error) {
	var markup tg.InlineKeyboardMarkup

	order, err := s.OrderBook.GetOrderByMessageId(ctx, userId, messageId)
	if err != nil {
		return markup, fmt.Errorf("failed to get payments markup: %w", err)
	}

	payment, err := s.OrderBook.GetPayment(ctx, uint64(order.ActivePaymentId.Int64))
	if err != nil {
		return markup, fmt.Errorf("failed to get payment: %w", err)
	}

	fullButton := tg.NewInlineKeyboardButtonData(bot.KbFullPayment, bot.KbDataFullPayment)
	partialButton := tg.NewInlineKeyboardButtonData(bot.KbPartialPayment, bot.KbDataPartialPayment)
	cancelButton := tg.NewInlineKeyboardButtonData(bot.KbCancel, bot.KbDataCancel)

	var markups [][]tg.InlineKeyboardButton
	if payment.PaymentMethod == types.PaymentMethodAcquiring {
		markups = [][]tg.InlineKeyboardButton{
			{fullButton},
			{cancelButton},
		}
	} else {
		markups = [][]tg.InlineKeyboardButton{
			{fullButton},
			{partialButton},
			{cancelButton},
		}
	}

	return tg.NewInlineKeyboardMarkup(markups...), nil
}

func refundAmountInlineKeyboard() tg.InlineKeyboardMarkup {
	fullButton := tg.NewInlineKeyboardButtonData(bot.KbFullRefund, bot.KbDataFullRefund)
	partialButton := tg.NewInlineKeyboardButtonData(bot.KbPartialRefund, bot.KbDataPartialRefund)
	cancelButton := tg.NewInlineKeyboardButtonData(bot.KbCancel, bot.KbDataCancel)

	markups := [][]tg.InlineKeyboardButton{
		{fullButton},
		{partialButton},
		{cancelButton},
	}
	return tg.NewInlineKeyboardMarkup(markups...)
}
