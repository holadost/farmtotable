package aragorn

import (
	"farmtotable/gandalf"
	"time"
)

type BaseAPIResponse struct {
	Status   uint64 `json:"status"`
	ErrorMsg string `json:"error_msg"`
}

type RegistrationStatusRet struct {
	RegistrationStatus bool `json:"registration_status"`
}

// Get User Arg and Ret.
type GetUserArg struct {
	UserID string `json:"user_id"`
}
type GetUserRet struct {
	BaseAPIResponse
	Data gandalf.User `json:"data"`
}

// Register user args and ret.
type RegisterUserArg struct {
	gandalf.User
}

type RegisterUserRet struct {
	BaseAPIResponse
	Data RegistrationStatusRet `json:"data"`
}

// Register Item Args and Ret.
type RegisterItemArg struct {
	SupplierID       string    `json:"supplier_id"`
	ItemName         string    `json:"item_name"`
	ItemDescription  string    `json:"item_description"`
	ItemQty          uint32    `json:"item_qty"`
	AuctionStartDate time.Time `json:"auction_start_date"`
	MinPrice         float32   `json:"min_price"`
	ItemTags         string    `json:"item_tags"`
}

type RegisterItemRet struct {
	BaseAPIResponse
	Data RegistrationStatusRet `json:"data"`
}

// Remove item args and ret.
type RemoveItemArg struct {
	ItemID string `json:"item_id"`
}

type RemoveItemRet struct {
	BaseAPIResponse
	Data RegistrationStatusRet `json:"data"`
}

// Register Item Args and Ret.
type RegisterSupplierArg struct {
	SupplierName        string `json:"supplier_name"`
	SupplierDescription string `json:"supplier_description"`
	SupplierEmailID     string `json:"supplier_email_id"`
	SupplierPhNum       string `json:"supplier_ph_num"`
	SupplierAddress     string `json:"supplier_address"`
	SupplierTags        string `json:"supplier_tags"`
}

type RegisterSupplierRet struct {
	BaseAPIResponse
	Data RegistrationStatusRet `json:"data"`
}

// Get all suppliers args and ret.
type GetAllSuppliersArg struct {
}

type GetAllSuppliersRet struct {
	BaseAPIResponse
	Data []gandalf.Supplier `json:"data"`
}

// Get supplier args and ret.
type GetSupplierArg struct {
	SupplierID string `json:"supplier_id"`
}

type GetSupplierRet struct {
	BaseAPIResponse
	Data gandalf.Supplier `json:"data"`
}

//  Get supplier items args and ret.
type GetSupplierItemsArg struct {
	GetSupplierArg
}

type GetSupplierItemsRet struct {
	BaseAPIResponse
	Data []gandalf.Item `json:"data"`
}

// Register bid args and ret.
type RegisterBidArg struct {
	ItemID    string  `json:"item_id"`
	UserID    string  `json:"user_id"`
	BidAmount float32 `json:"bid_amount"`
	BidQty    uint32  `json:"bid_qty"`
}

type RegisterBidRet struct {
	BaseAPIResponse
	Data RegistrationStatusRet `json:"data"`
}

// Fetch max bids args and ret.
type FetchMaxBidsArg struct {
	ItemIDs []string `json:"item_ids"`
}

type FetchMaxBidsRetData struct {
	ItemID string  `json:"item_id"`
	MaxBid float32 `json:"max_bid"`
}

type FetchMaxBidsRet struct {
	BaseAPIResponse
	Data []FetchMaxBidsRetData `json:"data"`
}

// Fetch all auctions args and ret.
type FetchAllAuctionsArg struct {
	StartID     uint64 `json:"start_id"`
	NumAuctions uint64 `json:"num_auctions"`
}

type FetchAllAuctionsRet struct {
	BaseAPIResponse
	Data []gandalf.Auction `json:"data"`
}

// User Bids
type GetAllUserBidsArg struct {
	GetUserArg
}

type GetUserBidsForItemArg struct {
	UserID string `json:"user_id"`
	ItemID string `json:"item_id"`
}

type GetUserBidsRet struct {
	BaseAPIResponse
	Data []gandalf.Bid `json:"data"`
}

type OrderRet struct {
	gandalf.Order
	ItemName        string `json:"item_name"`
	ItemDescription string `json:"item_description"`
}

type GetOrderArg struct {
	OrderID string `json:"order_id"`
}

type GetOrderRet struct {
	BaseAPIResponse
	Data OrderRet
}

type UpdateOrderArg struct {
	OrderID string `json:"order_id"`
	Status  uint32 `json:"status"`
}

type UpdateOrderRet struct {
	BaseAPIResponse
	Data RegistrationStatusRet `json:"data"`
}

type GetUserOrdersArg struct {
	GetUserArg
}

type GetUserOrdersRet struct {
	BaseAPIResponse
	Data []OrderRet `json:"data"`
}

type GetUserPaymentPendingOrdersArg struct {
	GetUserArg
}

type GetUserPaymentPendingOrdersRet struct {
	GetUserOrdersRet
}

type GetUserDeliveryPendingOrdersArg struct {
	GetUserArg
}

type GetUserDeliveryPendingOrdersRet struct {
	GetUserOrdersRet
}

type GetUserCancelledOrdersArg struct {
	GetUserArg
}

type GetUserCancelledOrdersRet struct {
	GetUserOrdersRet
}

type GetUserCompletedOrdersArg struct {
	GetUserArg
}

type GetUserCompletedOrdersRet struct {
	GetUserOrdersRet
}
