package gandalf

import (
	"errors"
	"fmt"
)

const (
	PGHost     = "localhost"
	PGPort     = 5432
	PGUser     = "postgres"
	PGPassword = "farmToTable4u"
	PGDbName   = "farmtotable"
)

const (
	SQLiteDBPath = "/tmp/gandalf.db"
)

/* OrderModel Status Codes */
const (
	KOrderPaymentPending  = 0
	KOrderDeliveryPending = 1
	KOrderComplete        = 2
	KOrderCancelled       = 3
)

const (
	KAuctionKeyPrefix = "auction"
	KBidKeyPrefix     = "bid"
	KKeyDelimiter     = "::"
)

type OrderStatus uint32

func (os OrderStatus) ToString() (OrderStatusStr, error) {
	if os == KOrderPaymentPending {
		return OrderStatusStr("KOrderPaymentPending"), nil
	} else if os == KOrderDeliveryPending {
		return OrderStatusStr("KOrderDeliveryPending"), nil
	} else if os == KOrderComplete {
		return OrderStatusStr("KOrderComplete"), nil
	} else if os == KOrderCancelled {
		return OrderStatusStr("KOrderCancelled"), nil
	} else {
		return OrderStatusStr("KUnsupportedStatus"), errors.New(fmt.Sprintf("unsupported order status: %d", os))
	}
}

type OrderStatusStr string

func (osr OrderStatusStr) ToUint32() (OrderStatus, error) {
	if osr == "KOrderPaymentPending" {
		return OrderStatus(KOrderPaymentPending), nil
	} else if osr == "KOrderDeliveryPending" {
		return OrderStatus(KOrderDeliveryPending), nil
	} else if osr == "KOrderComplete" {
		return OrderStatus(KOrderComplete), nil
	} else if osr == "KOrderCancelled" {
		return OrderStatus(KOrderCancelled), nil
	} else {
		return 9999999, errors.New(fmt.Sprintf("unsupported error status: %s", osr))
	}
}
