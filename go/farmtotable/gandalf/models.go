package gandalf

import (
	"time"
)

/* Database model/schema */
type User struct {
	UserID  string `gorm:"type:varchar(100);PRIMARY_KEY" json:"user_id"`
	Name    string `gorm:"type:varchar(255);NOT NULL" json:"name"`
	EmailID string `gorm:"type:varchar(255);NOT NULL;index" json:"email_id"`
	PhNum   string `gorm:"type:varchar(20);NOT NULL;index" json:"ph_num"`
	Address string `gorm:"NOT NULL" json:"address"`
}

type Supplier struct {
	SupplierID          string `gorm:"type:varchar(32);PRIMARY_KEY" json:"supplier_id"`
	SupplierName        string `gorm:"type:varchar(255)" json:"supplier_name"`
	SupplierEmailID     string `gorm:"type:varchar(255)" json:"supplier_email_id"`
	SupplierPhNum       string `gorm:"type:varchar(20)" json:"supplier_ph_num"`
	SupplierAddress     string `gorm:"type:varchar(255)" json:"supplier_address"`
	SupplierDescription string `gorm:"type:varchar(255)" json:"supplier_desc"`
	SupplierTags        string `gorm:"type:varchar(255)" json:"supplier_tags"`
}

type Item struct {
	ItemID           string    `gorm:"type:varchar(32);PRIMARY_KEY" json:"item_id"`
	SupplierID       string    `gorm:"type:varchar(100);index" json:"user_id"`
	ItemName         string    `gorm:"type:varchar(255);NOT NULL" json:"item_name"`
	ItemDescription  string    `gorm:"NOT NULL" json:"item_description"`
	ItemQty          uint32    `json:"item_qty"`
	AuctionStartTime time.Time `json:"auction_start_time"`
	MinPrice         float32   `json:"min_price"`
}

type Bid struct {
	ItemID    string  `gorm:"type:varchar(32);index" json:"item_id"`
	UserID    string  `gorm:"type:varchar(100);index" json:"user_id"`
	BidAmount float32 `json:"bid_amount"`
	BidQty    uint32  `json:"bid_qty"`
}

type Auction struct {
	ItemID              string    `gorm:"type:varchar(32);PRIMARY_KEY" json:"item_id"`
	ItemName            string    `gorm:"type:varchar(255)" json:"item_name"`
	ItemQty             uint32    `json:"item_qty"`
	AuctionStartTime    time.Time `json:"auction_start_time"`
	AuctionDurationSecs uint64    `json:"auction_duration_secs"`
	MaxBid              float32   `json:"max_bid"`
}

type Order struct {
	OrderID   string  `gorm:"type:varchar(32);PRIMARY_KEY" json:"order_id"`
	UserID    string  `gorm:"type:varchar(100);index" json:"user_id"`
	ItemID    string  `gorm:"type:varchar(32);index" json:"item_id"`
	ItemQty   uint32  `json:"item_qty"`
	ItemPrice float32 `json:"item_price"`
	Status    uint32  `json:"status"`
}
