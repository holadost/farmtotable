package gandalf

import (
	"time"
)

/* Database model/schema */
type UserModel struct {
	UserID  string `gorm:"type:varchar(100);PRIMARY_KEY" json:"user_id"`
	Name    string `gorm:"type:varchar(255);NOT NULL" json:"name"`
	EmailID string `gorm:"type:varchar(255);NOT NULL;index" json:"email_id"`
	PhNum   string `gorm:"type:varchar(20);NOT NULL;index" json:"ph_num"`
	Address string `gorm:"NOT NULL" json:"address"`
}

type SupplierModel struct {
	SupplierID          string `gorm:"type:varchar(32);PRIMARY_KEY" json:"supplier_id"`
	SupplierName        string `gorm:"type:varchar(255)" json:"supplier_name"`
	SupplierEmailID     string `gorm:"type:varchar(255)" json:"supplier_email_id"`
	SupplierPhNum       string `gorm:"type:varchar(20)" json:"supplier_ph_num"`
	SupplierAddress     string `gorm:"type:varchar(255)" json:"supplier_address"`
	SupplierDescription string `gorm:"type:varchar(255)" json:"supplier_desc"`
	SupplierTags        string `gorm:"type:varchar(255)" json:"supplier_tags"`
}

type ItemModel struct {
	// Item fields.
	ItemID          string  `gorm:"type:varchar(32);PRIMARY_KEY" json:"item_id"`
	SupplierID      string  `gorm:"type:varchar(100);index" json:"user_id"`
	ItemName        string  `gorm:"type:varchar(255);NOT NULL" json:"item_name"`
	ItemDescription string  `gorm:"type:TEXT;NOT NULL" json:"item_description"`
	ItemUnit        string  `gorm:"NOT NULL" json:"item_unit"`   // The unit of qty. For ex g or Kg.
	MinBidQty       uint32  `gorm:"NOT NULL" json:"min_bid_qty"` // Min bid qty. Like 100g or 1Kg.
	MaxBidQty       uint32  `gorm:"NOT NULL" json:"max_bid_qty"` // Max bid qty. Like 500g or 100Kg.
	ItemQty         uint32  `json:"item_qty"`                    // Total item quantity.
	ImageURL        string  `json:"image_url"`
	MinPrice        float32 `json:"min_price"`
	MaxPrice        float32 `json:"max_price"`
	// Auction fields for an item.
	AuctionStartTime    time.Time `json:"auction_start_time"`
	AuctionDurationSecs uint64    `json:"auction_duration_secs"`
	AuctionStarted      bool      `json:"auction_started"` // A flag to indicate whether the auction for this item has started.
	AuctionEnded        bool      `json:"auction_ended"`   // A flag to indicate whether auction has ended.
	AuctionDecided      bool      `json:"auction_decided"` // A flag to indicate whether the auction has been decided.
}

type BidModel struct {
	ItemID    string  `gorm:"type:varchar(32);index" json:"item_id"`
	UserID    string  `gorm:"type:varchar(100);index" json:"user_id"`
	BidAmount float32 `json:"bid_amount"`
	BidQty    uint32  `json:"bid_qty"`
}

/*
Note: The auction model replicates some of the data in ItemModel so that we can avoid
JOINS in the critical path.
*/
type AuctionModel struct {
	ID                  uint      `gorm:"PRIMARY_KEY;autoIncrement" json:"id"`
	ItemID              string    `gorm:"type:varchar(32)" json:"item_id"`    // Same as ItemModel ItemID.
	ItemName            string    `gorm:"type:varchar(255)" json:"item_name"` // Same as ItemModel ItemName.
	ItemQty             uint32    `json:"item_qty"`                           // Same as the ItemModel ItemQty.
	ImageURL            string    `json:"image_url"`                          // Same as the ItemModel ImageURL
	ItemUnit            string    `gorm:"NOT NULL" json:"item_unit"`          // The unit of qty. For ex g or Kg.
	MinBidQty           uint32    `gorm:"NOT NULL" json:"min_bid_qty"`        // Min bid qty. Like 100g or 1Kg.
	MaxBidQty           uint32    `gorm:"NOT NULL" json:"max_bid_qty"`        // Max bid qty. Like 500g or 100Kg.
	AuctionStartTime    time.Time `json:"auction_start_time"`
	AuctionDurationSecs uint64    `json:"auction_duration_secs"`
	MinBid              float32   `json:"min_bid"` // Same as the ItemModel MinBid.
	MaxBid              float32   `json:"max_bid"`
}

type OrderModel struct {
	ID        uint      `gorm:"PRIMARY_KEY;autoIncrement" json:"id"`
	OrderID   string    `gorm:"type:varchar(32);UNIQUE;index" json:"order_id"`
	UserID    string    `gorm:"type:varchar(100);index" json:"user_id"`
	ItemID    string    `gorm:"type:varchar(32);index" json:"item_id"`
	ItemQty   uint32    `json:"item_qty"`
	ItemPrice float32   `json:"item_price"`
	Status    uint32    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
