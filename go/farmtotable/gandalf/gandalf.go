package gandalf

import (
	"farmtotable/gandalf/backend"
	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
	"time"
)

type Gandalf struct {
	sqlBackend backend.SqlBackend
	db         *gorm.DB
}

func NewSqliteGandalf() *Gandalf {
	gandalf := &Gandalf{}
	return gandalf
}

func NewPostgresGandalf() *Gandalf {
	gandalf := &Gandalf{}
	return gandalf
}

func (gandalf *Gandalf) RegisterUser(userID string, name string, emailID string, phNum string, address string) error {
	user := &User{
		UserID:  userID,
		Name:    name,
		EmailID: emailID,
		PhNum:   phNum,
		Address: address,
	}
	db := gandalf.sqlBackend.GetDB()
	dbc := db.Create(user)
	if dbc.Error != nil {
		return dbc.Error
	}
	return nil
}

func (gandalf *Gandalf) GetUserByID(userID string) (user User) {
	gandalf.db.Where("user_id = ?", userID).First(&user)
	return
}

func (gandalf *Gandalf) GetUserByEmailID(emailID string) (user User) {
	gandalf.db.Where("email_id = ?", emailID).First(&user)
	return
}

func (gandalf *Gandalf) GetUserByPhNo(phNum string) (user User) {
	gandalf.db.Where("ph_num = ?", phNum).First(&user)
	return
}

func (gandalf *Gandalf) RegisterItem(itemName string, itemDesc string, itemQty uint32, auctionStartTime time.Time, minPrice float32) error {
	var dbc *gorm.DB
	var err error
	err = nil
	item := &Item{
		ItemName:         itemName,
		ItemDescription:  itemDesc,
		ItemQty:          itemQty,
		AuctionStartTime: auctionStartTime,
		MinPrice:         minPrice,
	}
	for ii := 0; ii < 5; ii++ {
		item.ItemID = xid.New().String()
		dbc = gandalf.db.Create(item)
		if dbc.Error != nil {
			// TODO: Ensure that this was a primary key error before retrying.
			// Retry with a new item ID.
			err = dbc.Error
			continue
		} else {
			break
		}
	}
	return err
}

func (gandalf *Gandalf) GetItem(itemID string) (item Item) {
	gandalf.db.Where("item_id = ?", itemID).First(&item)
	return
}

func (gandalf *Gandalf) EditItem(itemID string, itemName string, itemDesc string, itemQty uint32, auctionStartTime time.Time, minPrice float32) error {
	item := Item{
		ItemID:           itemID,
		ItemName:         itemName,
		ItemDescription:  itemDesc,
		ItemQty:          itemQty,
		AuctionStartTime: auctionStartTime,
		MinPrice:         minPrice,
	}
	dbc := gandalf.db.Updates(&item)
	if dbc.Error != nil {
		return dbc.Error
	}
	return nil
}

func (gandalf *Gandalf) AddBid(itemID string, userID string, bidAmount float32, bidQty uint32) error {
	// Registers the user bid.
	bid := Bid{
		ItemID:    itemID,
		UserID:    userID,
		BidAmount: bidAmount,
		BidQty:    bidQty,
	}
	dbc := gandalf.db.Create(&bid)
	if dbc.Error != nil {
		return dbc.Error
	}

	// Update the max bid.
	auction := Auction{}
	tx := gandalf.db.Begin()
	tx.Where("item_id = ?", itemID).First(&auction)
	if auction.MaxBid < bidAmount {
		auction.MaxBid = bidAmount
		dbc := tx.Updates(&auction)
		if dbc.Error != nil {
			tx.Rollback()
			return dbc.Error
		}
		tx.Commit()
	} else {
		tx.Rollback()
	}
	return nil
}

func (gandalf *Gandalf) GetMaxBid(itemIDs []string) ([]Auction, error) {
	var auctions []Auction
	dbc := gandalf.db.Where("item_id IN (?)", itemIDs).Find(&auctions)
	if dbc.Error != nil {
		return auctions, dbc.Error
	}
	return auctions, nil
}

func (gandalf *Gandalf) GetAllAuctions(startIndex uint64, numAuctions uint64) ([]Auction, error) {
	var auctions []Auction
	dbc := gandalf.db.Offset(startIndex).Limit(numAuctions).Find(&auctions)
	if dbc.Error != nil {
		return auctions, dbc.Error
	}
	return auctions, nil
}

func (gandalf *Gandalf) GetUserOrders(userID string) ([]Order, error) {
	var orders []Order
	dbc := gandalf.db.Where("user_id = ?", userID).Find(&orders)
	if dbc.Error != nil {
		return orders, dbc.Error
	}
	return orders, nil
}

func (gandalf *Gandalf) GetUserPaymentPendingOrders(userID string) ([]Order, error) {
	var orders []Order
	dbc := gandalf.db.Where("user_id = ? AND status = ?", userID, KOrderPaymentPending).Find(&orders)
	if dbc.Error != nil {
		return orders, dbc.Error
	}
	return orders, nil

}

func (gandalf *Gandalf) GetUserDeliveryPendingOrders(userID string) ([]Order, error) {
	var orders []Order
	dbc := gandalf.db.Where("user_id = ? AND status = ?", userID, KOrderDeliveryPending).Find(&orders)
	if dbc.Error != nil {
		return orders, dbc.Error
	}
	return orders, nil

}

func (gandalf *Gandalf) GetUserCompletedOrders(userID string) ([]Order, error) {
	var orders []Order
	dbc := gandalf.db.Where("user_id = ? AND status = ?", userID, KOrderDelivered).Find(&orders)
	if dbc.Error != nil {
		return orders, dbc.Error
	}
	return orders, nil
}

func (gandalf *Gandalf) GetOrder(orderID string) (Order, error) {
	var order Order
	dbc := gandalf.db.Where("order_id = ?", orderID, KOrderDelivered).First(&order)
	if dbc.Error != nil {
		return order, dbc.Error
	}
	return order, nil
}

func (gandalf *Gandalf) UpdateOrderStatus(orderID string, status uint32) error {
	var order Order
	dbc := gandalf.db.Model(&order).Where("order_id = ?", status).Update("status", status)
	if dbc.Error != nil {
		return dbc.Error
	}
	return nil
}
