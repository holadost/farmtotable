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

// Get User Ret.
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
	ItemTags         []string  `json:"item_tags"`
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
	SupplierID string
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
	Data []gandalf.Auction
}