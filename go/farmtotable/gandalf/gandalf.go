package gandalf

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/rs/xid"
	bulk "github.com/sunary/gorm-bulk-insert"
	"os"
	"time"
)

type Gandalf struct {
	Db *gorm.DB
}

func NewSqliteGandalf() *Gandalf {
	gandalf := &Gandalf{}
	db, err := gorm.Open("sqlite3", SQLiteDBPath)
	if err != nil {
		panic("Unable to open Sqlite database")
		return nil
	}
	gandalf.Db = db
	gandalf.Initialize()
	return gandalf
}

func NewPostgresGandalf() *Gandalf {
	gandalf := &Gandalf{}
	addrString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s", PGHost, PGPort, PGUser, PGDbName, PGPassword)
	db, err := gorm.Open("postgres", addrString)
	if err != nil {
		panic("Unable to open postgres database")
		return nil
	}
	gandalf.Db = db
	gandalf.Initialize()
	return gandalf
}

func NewGandalf() *Gandalf {
	value, exists := os.LookupEnv("FTT_GANDALF_BACKEND")
	if !exists {
		return NewSqliteGandalf()
	}
	if value == "POSTGRES" {
		return NewPostgresGandalf()
	} else {
		return NewSqliteGandalf()
	}
}

func (gandalf *Gandalf) Initialize() error {
	user := User{}
	supplier := Supplier{}
	item := Item{}
	bid := Bid{}
	auction := Auction{}
	order := Order{}
	dbc := gandalf.Db.AutoMigrate(&user, &supplier, &item, &auction, &bid, &order)
	if dbc != nil && dbc.Error != nil {
		panic("Unable to create database")
		return dbc.Error
	}
	return nil
}

func (gandalf *Gandalf) Close() {
	gandalf.Db.Close()
}

func (gandalf *Gandalf) RegisterUser(userID string, name string, emailID string, phNum string, address string) error {
	user := &User{
		UserID:  userID,
		Name:    name,
		EmailID: emailID,
		PhNum:   phNum,
		Address: address,
	}
	dbc := gandalf.Db.Create(user)
	if dbc.Error != nil {
		return dbc.Error
	}
	return nil
}

func (gandalf *Gandalf) GetUserByID(userID string) (user User) {
	gandalf.Db.Where("user_id = ?", userID).First(&user)
	return
}

func (gandalf *Gandalf) GetUserByEmailID(emailID string) (user User) {
	gandalf.Db.Where("email_id = ?", emailID).First(&user)
	return
}

func (gandalf *Gandalf) GetUserByPhNo(phNum string) (user User) {
	gandalf.Db.Where("ph_num = ?", phNum).First(&user)
	return
}

func (gandalf *Gandalf) RegisterSupplier(supplierName string, emailID string, phNum string,
	address string, description string, tags string) error {
	var dbc *gorm.DB
	var err error
	err = nil
	supplier := &Supplier{
		SupplierName:        supplierName,
		SupplierAddress:     address,
		SupplierEmailID:     emailID,
		SupplierDescription: description,
		SupplierPhNum:       phNum,
		SupplierTags:        tags,
	}
	for ii := 0; ii < 5; ii++ {
		supplier.SupplierID = xid.New().String()
		dbc = gandalf.Db.Create(supplier)
		if dbc.Error != nil {
			// Retry with a new item ID.
			err = dbc.Error
			continue
		} else {
			break
		}
	}
	return err
}

func (gandalf *Gandalf) GetSupplierByID(supplierID string) (supplier Supplier) {
	gandalf.Db.Where("supplier_id = ?", supplierID).First(&supplier)
	return
}

func (gandalf *Gandalf) GetAllSuppliers() ([]Supplier, error) {
	var suppliers []Supplier
	dbc := gandalf.Db.Find(&suppliers)
	if dbc.Error != nil {
		return suppliers, dbc.Error
	}
	return suppliers, nil
}

func (gandalf *Gandalf) RegisterItem(supplierID string, itemName string, itemDesc string, itemQty uint32,
	auctionStartTime time.Time, minPrice float32) error {
	var dbc *gorm.DB
	var err error
	err = nil
	item := &Item{
		SupplierID:       supplierID,
		ItemName:         itemName,
		ItemDescription:  itemDesc,
		ItemQty:          itemQty,
		AuctionStartTime: auctionStartTime,
		MinPrice:         minPrice,
	}
	for ii := 0; ii < 5; ii++ {
		item.ItemID = xid.New().String()
		dbc = gandalf.Db.Create(item)
		if dbc.Error != nil {
			// Retry with a new item ID.
			err = dbc.Error
			continue
		} else {
			break
		}
	}
	return err
}

func (gandalf *Gandalf) GetSupplierItems(supplierID string) ([]Item, error) {
	var items []Item
	dbc := gandalf.Db.Where("supplier_id = ?", supplierID).Find(&items)
	if dbc.Error != nil {
		return items, dbc.Error
	}
	return items, nil
}

func (gandalf *Gandalf) GetItem(itemID string) (item Item) {
	gandalf.Db.Where("item_id = ?", itemID).First(&item)
	return
}

func (gandalf *Gandalf) GetItems(itemIDs []string) (items []Item) {
	gandalf.Db.Where("item_id IN ?", itemIDs).Find(&items)
	return
}

func (gandalf *Gandalf) ScanItems(startIndex uint64, numItems uint64) ([]Item, error) {
	var items []Item
	dbc := gandalf.Db.Offset(startIndex).Limit(numItems).Find(&items)
	if dbc.Error != nil {
		return items, dbc.Error
	}
	return items, nil
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
	dbc := gandalf.Db.Model(&item).Updates(&item)
	if dbc.Error != nil {
		return dbc.Error
	}
	return nil
}

func (gandalf *Gandalf) DeleteItem(itemID string) error {
	item := Item{
		ItemID: itemID,
	}
	dbc := gandalf.Db.Model(&item).Delete(&item)
	if dbc.Error != nil {
		return dbc.Error
	}
	return nil
}

func (gandalf *Gandalf) RegisterAuctions(auctions []Auction) error {
	var insertRecords []interface{}
	for ii := 0; ii < len(auctions); ii++ {
		insertRecords = append(insertRecords, auctions[ii])
	}
	err := bulk.BulkInsertWithTableName(gandalf.Db, "auctions", insertRecords)
	if err != nil {
		return err
	}
	return nil
}

func (gandalf *Gandalf) RegisterBid(itemID string, userID string, bidAmount float32, bidQty uint32) error {
	// Registers the user bid.
	bid := Bid{
		ItemID:    itemID,
		UserID:    userID,
		BidAmount: bidAmount,
		BidQty:    bidQty,
	}
	dbc := gandalf.Db.Create(&bid)
	if dbc.Error != nil {
		return dbc.Error
	}

	// Update the max bid.
	auction := Auction{}
	tx := gandalf.Db.Begin()
	tx.Where("item_id = ?", itemID).First(&auction)
	if auction.MaxBid < bidAmount {
		auction.MaxBid = bidAmount
		dbc := tx.Model(&auction).Updates(&auction)
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

/* Returns all the live auctions that the user has bid on. */
func (gandalf *Gandalf) GetUserBids(userID string) ([]Bid, error) {
	var bids []Bid
	dbc := gandalf.Db.Where("user_id = ?", userID).Find(&bids)
	if dbc.Error != nil {
		return bids, dbc.Error
	}
	return bids, nil
}

/* Returns the bids for a given item starting from 'start' row upto numBids rows. */
func (gandalf *Gandalf) ScanItemBids(itemID string, startIndex uint64, numBids uint64) ([]Bid, error) {
	var bids []Bid
	dbc := gandalf.Db.Where("item_id = ?", itemID).Offset(startIndex).Limit(numBids).Find(&bids)
	if dbc.Error != nil {
		return bids, dbc.Error
	}
	return bids, nil
}

/* Returns all the bids for a given item. This method must be used with care as there could potentially be
millions of such records. */
func (gandalf *Gandalf) GetAllItemBids(itemID string) ([]Bid, error) {
	var bids []Bid
	dbc := gandalf.Db.Where("item_id = ?", itemID).Find(&bids)
	if dbc.Error != nil {
		return bids, dbc.Error
	}
	return bids, nil
}

/* Fetches the max bids for the given items. */
func (gandalf *Gandalf) GetMaxBids(itemIDs []string) ([]Auction, error) {
	var auctions []Auction
	dbc := gandalf.Db.Where("item_id IN (?)", itemIDs).Find(&auctions)
	if dbc.Error != nil {
		return auctions, dbc.Error
	}
	return auctions, nil
}

/* Fetches all the auctions starting from start index upto numAuctions. */
func (gandalf *Gandalf) GetAllAuctions(startIndex uint64, numAuctions uint64) ([]Auction, error) {
	var auctions []Auction
	dbc := gandalf.Db.Offset(startIndex).Limit(numAuctions).Find(&auctions)
	if dbc.Error != nil {
		return auctions, dbc.Error
	}
	return auctions, nil
}

/* Adds the given orders to the backend. */
func (gandalf *Gandalf) AddOrders(orders []Order) error {
	// Adds the given orders to the database.
	return nil
}

/* Gets all user orders. */
func (gandalf *Gandalf) GetUserOrders(userID string) ([]Order, error) {
	var orders []Order
	dbc := gandalf.Db.Where("user_id = ?", userID).Find(&orders)
	if dbc.Error != nil {
		return orders, dbc.Error
	}
	return orders, nil
}

/* Gets all user orders whose payment is pending. */
func (gandalf *Gandalf) GetUserPaymentPendingOrders(userID string) ([]Order, error) {
	var orders []Order
	dbc := gandalf.Db.Where("user_id = ? AND status = ?", userID, KOrderPaymentPending).Find(&orders)
	if dbc.Error != nil {
		return orders, dbc.Error
	}
	return orders, nil
}

/*
Gets all orders whose payment is pending. This method must be used with care as it could potentially end up
returning 100,000 of rows.
*/
func (gandalf *Gandalf) GetAllPaymentPendingOrders() ([]Order, error) {
	var orders []Order
	dbc := gandalf.Db.Where("status = ?", KOrderPaymentPending).Find(&orders)
	if dbc.Error != nil {
		return orders, dbc.Error
	}
	return orders, nil
}

/* Gets all user orders whose delivery is pending. */
func (gandalf *Gandalf) GetUserDeliveryPendingOrders(userID string) ([]Order, error) {
	var orders []Order
	dbc := gandalf.Db.Where("user_id = ? AND status = ?", userID, KOrderDeliveryPending).Find(&orders)
	if dbc.Error != nil {
		return orders, dbc.Error
	}
	return orders, nil
}

/* Gets all orders whose delivery is pending. */
func (gandalf *Gandalf) GetAllDeliveryPendingOrders() ([]Order, error) {
	var orders []Order
	dbc := gandalf.Db.Where("status = ?", KOrderDeliveryPending).Find(&orders)
	if dbc.Error != nil {
		return orders, dbc.Error
	}
	return orders, nil
}

/* Gets all user orders that are complete */
func (gandalf *Gandalf) GetUserCompletedOrders(userID string) ([]Order, error) {
	var orders []Order
	dbc := gandalf.Db.Where("user_id = ? AND status = ?", userID, KOrderComplete).Find(&orders)
	if dbc.Error != nil {
		return orders, dbc.Error
	}
	return orders, nil
}

/* Get order information */
func (gandalf *Gandalf) GetOrder(orderID string) (Order, error) {
	var order Order
	dbc := gandalf.Db.Where("order_id = ?", orderID, KOrderComplete).First(&order)
	if dbc.Error != nil {
		return order, dbc.Error
	}
	return order, nil
}

/* Update order status. */
func (gandalf *Gandalf) UpdateOrderStatus(orderID string, status uint32) error {
	var order Order
	dbc := gandalf.Db.Model(&order).Where("order_id = ?", status).Update("status", status)
	if dbc.Error != nil {
		return dbc.Error
	}
	return nil
}
