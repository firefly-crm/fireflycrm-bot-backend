package types

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	PaymentMethodAcquiring PaymentMethod = iota
	PaymentMethodCard2Card
	PaymentMethodCash
)

const (
	EditStateNone EditState = iota
	EditStateWaitingItemName
	EditStateWaitingItemPrice
	EditStateWaitingItemQuantity
	EditStateWaitingCustomerEmail
	EditStateWaitingCustomerInstagram
	EditStateWaitingCustomerPhone
	EditStateWaitingPaymentAmount
	EditStateWaitingRefundAmount
	EditStateWaitingOrderDueDate
	EditStateWaitingOrderDescription
)

const (
	OrderStateForming OrderState = iota
	OrderStateInProgress
	OrderStateDone
	OrderStateDeleted = 99
)

const (
	DisplayModeFull DisplayMode = iota
	DisplayModeCollapsed
	DisplayModeDeleted
)

const (
	ReceiptItemTypeGoods ReceiptItemType = iota
	ReceiptItemTypeDelivery
)

type (
	OrderState      byte
	EditState       byte
	ReceiptItemType byte
	PaymentMethod   byte
	DisplayMode     byte

	OrderOptions struct {
		Description    string
		PaymentType    PaymentMethod
		CustomerName   string
		CustomerEmail  string
		CustomerPhone  string
		CustomerSocial string
	}

	ReceiptItem struct {
		Id          uint64        `db:"id"`
		Name        string        `db:"name"`
		ItemId      sql.NullInt64 `db:"item_id"`
		OrderId     uint64        `db:"order_id"`
		Price       uint32        `db:"price"`
		Quantity    uint32        `db:"quantity"`
		Initialised bool          `db:"initialised"`
		CreatedAt   time.Time     `db:"created_at"`
		UpdatedAt   time.Time     `db:"updated_at"`
	}

	Customer struct {
		Id        uint64         `db:"id"`
		Email     sql.NullString `db:"email"`
		Phone     sql.NullString `db:"phone"`
		Name      sql.NullString `db:"name"`
		Instagram sql.NullString `db:"instagram"`
		CreatedAt time.Time      `db:"created_at"`
		UpdatedAt time.Time      `db:"updated_at"`
	}

	Order struct {
		Id              uint64        `db:"id"`
		UserOrderId     uint64        `db:"user_order_id"`
		HintMessageId   sql.NullInt64 `db:"hint_message_id"`
		UserId          uint64        `db:"user_id"`
		CustomerId      sql.NullInt64 `db:"customer_id"`
		Description     string        `db:"description"`
		ActiveItemId    sql.NullInt64 `db:"active_item_id"`
		ActivePaymentId sql.NullInt64 `db:"active_payment_id"`
		OrderState      OrderState    `db:"order_state"`
		EditState       EditState     `db:"edit_state"`
		Amount          uint32        `db:"amount"`
		PayedAmount     uint32        `db:"payed_amount"`
		RefundAmount    uint32        `db:"refund_amount"`
		DueDate         sql.NullTime  `db:"due_date"`
		CreatedAt       time.Time     `db:"created_at"`
		UpdatedAt       time.Time     `db:"updated_at"`
		ReceiptItems    []ReceiptItem
		Payments        []Payment
	}

	Payment struct {
		Id            uint64        `db:"id"`
		OrderId       uint64        `db:"order_id"`
		Amount        uint32        `db:"amount"`
		PaymentMethod PaymentMethod `db:"payment_method"`
		PaymentLink   string        `db:"payment_link"`
		BankPaymentId string        `db:"bank_payment_id"`
		Payed         bool          `db:"payed"`
		Refunded      bool          `db:"refunded"`
		RefundAmount  uint32        `db:"refund_amount"`
		Expired       bool          `db:"expired"`
		CreatedAt     time.Time     `db:"created_at"`
		UpdatedAt     time.Time     `db:"updated_at"`
		PayedAt       sql.NullTime  `db:"payed_at"`
	}

	OrderMessage struct {
		Id          uint64      `db:"id"`
		OrderId     uint64      `db:"order_id"`
		UserId      uint64      `db:"user_id"`
		DisplayMode DisplayMode `db:"display_mode"`
	}

	User struct {
		Id                   uint64    `db:"id"`
		IsMerchant           bool      `db:"is_merchant"`
		ActiveOrderId        uint64    `db:"active_order_id"`
		MerchantId           string    `db:"merchant_id"`
		SecretKey            string    `db:"secret_key"`
		ActiveOrderMessageId string    `db:"active_order_msg_id"`
		CreatedAt            time.Time `db:"created_at"`
		UpdatedAt            time.Time `db:"updated_at"`
	}
)

func (o Order) MessageString(c *Customer, mode DisplayMode) string {
	switch mode {
	case DisplayModeFull:
		return o.getFullMessageString(c)
	case DisplayModeCollapsed:
		return o.getCollapsedMessageString(c)
	case DisplayModeDeleted:
		return o.getDeletedMessageString()
	}

	return ""
}

func (o Order) getCollapsedMessageString(c *Customer) string {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		loc = time.Now().Location()
	}
	createdAt := o.CreatedAt.In(loc).Format("01.02.2006")

	var payed string
	if o.Amount <= o.PayedAmount {
		payed = "Оплачен"
	} else {
		payed = fmt.Sprintf("Оплачено %d₽ из %d₽", o.PayedAmount/100, o.Amount/100)
	}

	result := fmt.Sprintf("<b>Заказ: #%d</b> <i>(%s)</i>\n", o.UserOrderId, o.OrderState.MessageString())

	if c != nil {
		if c.Instagram.Valid {
			result += fmt.Sprintf("<b>Клиент:</b> <a href=\"https://instagram.com/%[1]s\">@%[1]s</a>\n", c.Instagram.String)
		} else if c.Phone.Valid {
			result += fmt.Sprintf("<b>Клиент:</b> <a href=\"https://wa.me/%s\">%s</a>\n", c.Phone.String, formatPhone(c.Phone.String))
		}
	}

	if o.Description != "" {
		result += fmt.Sprintf("<i>%s</i>\n", o.Description)
	}

	dueDate := "N/A"
	if o.DueDate.Valid {
		dueDate = o.DueDate.Time.In(loc).Format("01.02.2006")
	}

	result += fmt.Sprintf("<b>Срок сдачи:</b> %s; %s", dueDate, payed)

	return result
}

func (o Order) getDeletedMessageString() string {
	return fmt.Sprintf("<b>Заказ #%d.</b> <i>Удалён</i>.", o.UserOrderId)
}

func (o Order) getFullMessageString(c *Customer) string {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		loc = time.Now().Location()
	}

	createdAt := o.CreatedAt.In(loc).Format("02.01.2006")

	amount := float32(o.Amount) / 100.0

	result := fmt.Sprintf(
		`<b>Заказ #%d</b> <i>(%s)</i>
<b>Создан:</b> %s`, o.UserOrderId, o.OrderState.MessageString(), createdAt)

	if o.DueDate.Valid {
		dueDate := o.DueDate.Time.In(loc).Format("02.01.2006")

		result += fmt.Sprintf(`
<b>Срок сдачи:</b> %s`, dueDate)
	}

	result += fmt.Sprintf(`
<b>Сумма:</b> %.2f₽
`, amount)

	if o.Amount != 0 && o.PayedAmount != 0 {
		if o.PayedAmount >= o.Amount {
			result += "<b>Оплачен:</b> полностью\n"
		} else {
			payedAmount := float32(o.PayedAmount) / 100.0
			result += fmt.Sprintf("<b>Оплачено:</b> %.2f₽\n", payedAmount)
			restAmount := float32(o.Amount-o.PayedAmount) / 100.0
			result += fmt.Sprintf("<b>Остаток:</b> %.2f₽\n", restAmount)
		}
	}

	if o.RefundAmount != 0 {
		if o.RefundAmount >= o.Amount {
			result += "<b>Возврат:</b> в полном объеме\n"
		} else {
			refundAmount := float32(o.RefundAmount) / 100.0
			result += fmt.Sprintf("<b>Возвращено:</b> %.2f₽\n", refundAmount)
		}
	}

	if o.Description != "" {
		result += fmt.Sprintf("\n<i>%s</i>\n", o.Description)
	}

	result += "\n<b>Позиции</b>\n"

	if o.ReceiptItems != nil {
		for _, i := range o.ReceiptItems {
			price := float32(i.Price) / 100.0
			result += fmt.Sprintf("<code>- %s %.2f₽ x%d</code>\n", i.Name, price, i.Quantity)
		}
	}

	result += "\n<b>Клиент</b>"
	if c != nil {

		if c.Name.Valid {
			result += fmt.Sprintf("\n<b>Имя:</b> %s", c.Name.String)
		}

		if c.Email.Valid {
			result += fmt.Sprintf("\n<b>E-Mail:</b> %s", c.Email.String)
		}

		if c.Phone.Valid {
			result += fmt.Sprintf("\n<b>Телефон:</b> <a href=\"https://wa.me/%s\">%s</a>", c.Phone.String, formatPhone(c.Phone.String))
		}

		if c.Instagram.Valid {
			result += fmt.Sprintf("\n<b>Instagram:</b> <a href=\"https://instagram.com/%[1]s\">@%[1]s</a>", c.Instagram.String)
		}

	} else {
		result += "\n‼️ Данные не заполнены"
	}

	result += "\n"

	if o.Payments != nil {
		result += "\n<b>Данные по оплате</b>"
		if len(o.Payments) == 0 {
			result += "\nНе найдено"
		} else {
			sort.Sort(PaymentsByCreatedAt(o.Payments))
			for i, p := range o.Payments {
				result += fmt.Sprintf("\n%s\n", p.MessageString(i+1))
			}
		}
	}

	return result
}

func (p Payment) MessageString(id int) string {
	result := fmt.Sprintf(`<b>Платеж #%d.</b> %s.`, id, p.PaymentMethod.MessageString(p.PaymentLink))

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		loc = time.Now().Location()
	}

	amount := float32(p.Amount) / 100.0
	result += fmt.Sprintf("\n<b>Сумма:</b> %.2f₽", amount)

	if p.PaymentMethod == PaymentMethodAcquiring {
		createdAt := p.CreatedAt.In(loc).Format("02 Jan 2006 15:04")
		result += fmt.Sprintf("\n<b>Создан:</b> %s", createdAt)

		if p.Payed && p.PayedAt.Valid {
			payedAt := p.PayedAt.Time.In(loc).Format("02 Jan 2006 15:04")
			result += fmt.Sprintf("\n<b>Оплачен:</b> %s", payedAt)
		} else {
			result += "\n<b>Оплачен:</b> нет"
		}
	} else {
		payedAt := p.PayedAt.Time.In(loc).Format("02 Jan 2006 15:04")
		result += fmt.Sprintf("\n<b>Оплачен:</b> %s", payedAt)
	}

	if p.RefundAmount != 0 {
		refundAmount := float32(p.RefundAmount) / 100.0
		if p.RefundAmount == p.Amount {
			result += "\n<b>Возвращен:</b> в полном объеме"
		} else {
			result += fmt.Sprintf("\n<b>Возвращено:</b> %.2f₽", refundAmount)
		}
	}

	return result
}

func (p PaymentMethod) MessageString(link string) string {
	switch p {
	case PaymentMethodCard2Card:
		return "Перевод на карту"
	case PaymentMethodAcquiring:
		return fmt.Sprintf("Оплата по <a href=\"%s\">ссылке</a>", link)
	case PaymentMethodCash:
		return "Оплата наличными"
	default:
		return "Неизвестный тип оплаты"
	}
}

func (s OrderState) MessageString() string {
	switch s {
	case OrderStateForming:
		return "Формируется"
	case OrderStateDone:
		return "Завершен"
	case OrderStateDeleted:
		return "Удалён"
	case OrderStateInProgress:
		return "В работе"
	}
	return "Статус неизвестен"
}

func formatPhone(phone string) string {
	digits := strings.Split(phone, "")

	formatted := fmt.Sprintf("+%s(%s)%s-%s-%s",
		digits[0],
		strings.Join(digits[1:4], ""),
		strings.Join(digits[4:7], ""),
		strings.Join(digits[7:9], ""),
		strings.Join(digits[9:], ""))

	return formatted
}

type PaymentsByCreatedAt []Payment

func (p PaymentsByCreatedAt) Len() int {
	return len(p)
}

func (p PaymentsByCreatedAt) Less(i, j int) bool {
	return p[i].CreatedAt.Before(p[j].CreatedAt)
}

func (p PaymentsByCreatedAt) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
