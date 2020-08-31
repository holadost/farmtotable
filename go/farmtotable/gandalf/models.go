package gandalf

import (
	"time"
)

type baseModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

/* Database model/schema */
type User struct {
	UserID  string `gorm:"type:varchar(100);PRIMARY_KEY"`
	Name    string `gorm:"type:varchar(255);NOT NULL"`
	EmailID string `gorm:"type:varchar(255);NOT NULL;index"`
	PhNum   string `gorm:"type:varchar(20);NOT NULL;index"`
	Address string `gorm:"NOT NULL"`
}

type Item struct {
	ItemID           string `gorm:"type:varchar(32);PRIMARY_KEY"`
	UserID           string `gorm:"type:varchar(100);index"`
	ItemName         string `gorm:"type:varchar(255);NOT NULL"`
	ItemDescription  string `gorm:"NOT NULL"`
	ItemQty          uint32
	AuctionStartTime time.Time
	MinPrice         float32
}

type Bid struct {
	ItemID    string `gorm:"type:varchar(32);index"`
	UserID    string `gorm:"type:varchar(100);index"`
	BidAmount float32
	BidQty    uint32
}

type Auction struct {
	ItemID              string `gorm:"type:varchar(32);PRIMARY_KEY"`
	ItemQty             uint32
	AuctionStartTime    time.Time
	AuctionDurationSecs uint64
	MaxBid              float32
}

type Order struct {
	OrderID   string `gorm:"type:varchar(32);PRIMARY_KEY"`
	UserID    string `gorm:"type:varchar(100);index"`
	ItemID    string `gorm:"type:varchar(32);index"`
	ItemQty   uint32
	ItemPrice float32
	Status    uint32
}
