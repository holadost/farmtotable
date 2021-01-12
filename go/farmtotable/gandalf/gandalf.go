package gandalf

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/rs/xid"
	bulk "github.com/sunary/gorm-bulk-insert"
	"os"
	"strings"
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
	user := UserModel{}
	supplier := SupplierModel{}
	item := ItemModel{}
	bid := BidModel{}
	auction := AuctionModel{}
	order := OrderModel{}
	dbc := gandalf.Db.AutoMigrate(&user, &supplier, &item, &auction, &bid, &order)
	if dbc != nil && dbc.Error != nil {
		glog.Fatalf("Unable to create and initialize database due to err: %v", dbc.Error)
		return NewGandalfError(KGandalfBackendError, dbc.Error.Error())
	}
	return nil
}

func (gandalf *Gandalf) Close() {
	gandalf.Db.Close()
}

func (gandalf *Gandalf) RegisterUser(userID string, name string, emailID string, phNum string, address string) error {
	user := &UserModel{
		UserID:  userID,
		Name:    name,
		EmailID: emailID,
		PhNum:   phNum,
		Address: address,
	}
	dbc := gandalf.Db.Create(user)
	if dbc.Error != nil {
		return NewGandalfError(KGandalfBackendError, dbc.Error.Error())
	}
	return nil
}

func (gandalf *Gandalf) GetUserByID(userID string) (user UserModel) {
	gandalf.Db.Where("user_id = ?", userID).First(&user)
	return
}

func (gandalf *Gandalf) GetUserByEmailID(emailID string) (user UserModel) {
	gandalf.Db.Where("email_id = ?", emailID).First(&user)
	return
}

func (gandalf *Gandalf) GetUserByPhNo(phNum string) (user UserModel) {
	gandalf.Db.Where("ph_num = ?", phNum).First(&user)
	return
}

func (gandalf *Gandalf) RegisterSupplier(supplierName string, emailID string, phNum string,
	address string, description string, tags string) error {
	var dbc *gorm.DB
	var err error
	err = nil
	supplier := &SupplierModel{
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
	if err != nil {
		return NewGandalfError(KGandalfBackendError, err.Error())
	}
	return nil
}

func (gandalf *Gandalf) GetSupplierByID(supplierID string) (supplier SupplierModel) {
	gandalf.Db.Where("supplier_id = ?", supplierID).First(&supplier)
	return
}

func (gandalf *Gandalf) GetAllSuppliers() ([]SupplierModel, error) {
	var suppliers []SupplierModel
	dbc := gandalf.Db.Find(&suppliers)
	if dbc.Error != nil {
		return suppliers, NewGandalfError(KGandalfBackendError, dbc.Error.Error())
	}
	return suppliers, nil
}

func (gandalf *Gandalf) RegisterItem(
	supplierID string,
	itemName string,
	itemDesc string,
	itemQty uint32,
	auctionStartTime time.Time,
	minPrice float32,
	auctionDurationSecs uint32,
	imageUrl string,
	minBidQty uint32,
	maxBidQty uint32,
	itemUnit string) error {

	var dbc *gorm.DB
	var err error
	err = nil
	item := &ItemModel{
		SupplierID:          supplierID,
		ItemName:            itemName,
		ItemDescription:     itemDesc,
		ItemQty:             itemQty,
		AuctionStartTime:    auctionStartTime,
		MinPrice:            minPrice,
		AuctionDurationSecs: uint64(auctionDurationSecs),
		ImageURL:            imageUrl,
		MinBidQty:           minBidQty,
		MaxBidQty:           maxBidQty,
		ItemUnit:            itemUnit,
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
	if err != nil {
		return NewGandalfError(KGandalfBackendError, err.Error())
	}
	return nil
}

func (gandalf *Gandalf) GetSupplierItems(supplierID string) ([]ItemModel, error) {
	var items []ItemModel
	dbc := gandalf.Db.Where("supplier_id = ?", supplierID).Find(&items)
	if dbc.Error != nil {
		return items, NewGandalfError(KGandalfBackendError, dbc.Error.Error())
	}
	return items, nil
}

func (gandalf *Gandalf) GetItem(itemID string) (ItemModel, error) {
	var item ItemModel
	dbc := gandalf.Db.Where("item_id = ?", itemID).First(&item)
	if dbc.Error != nil {
		return item, NewGandalfError(KGandalfBackendError, dbc.Error.Error())
	}
	return item, nil
}

func (gandalf *Gandalf) GetItems(itemIDs []string) ([]ItemModel, error) {
	if len(itemIDs) == 0 {
		return []ItemModel{}, nil
	}
	var items []ItemModel
	var args []interface{}
	for _, itemID := range itemIDs {
		args = append(args, itemID)
	}
	// For some reason, gorm WHERE query with IN clause was failing. So we go with a raw query.
	query := "SELECT * FROM item_models WHERE item_id IN (?" + strings.Repeat(",?", len(args)-1) + ")"
	dbc := gandalf.Db.Raw(query, args...).Scan(&items)
	if dbc.Error != nil {
		glog.Errorf("Unable to query items due to error: %v", dbc.Error)
		return items, NewGandalfError(KGandalfBackendError, dbc.Error.Error())
	}
	return items, nil
}

func (gandalf *Gandalf) ScanItems(startIndex uint64, numItems uint64) ([]ItemModel, error) {
	var items []ItemModel
	dbc := gandalf.Db.Offset(startIndex).Limit(numItems).Find(&items)
	if dbc.Error != nil {
		return items, NewGandalfError(KGandalfBackendError, dbc.Error.Error())
	}
	return items, nil
}

func (gandalf *Gandalf) EditItem(itemID string, itemName string, itemDesc string, itemQty uint32, auctionStartTime time.Time, minPrice float32) error {
	item := ItemModel{
		ItemID:           itemID,
		ItemName:         itemName,
		ItemDescription:  itemDesc,
		ItemQty:          itemQty,
		AuctionStartTime: auctionStartTime,
		MinPrice:         minPrice,
	}
	dbc := gandalf.Db.Model(&item).Updates(&item)
	if dbc.Error != nil {
		return NewGandalfError(KGandalfBackendError, dbc.Error.Error())
	}
	return nil
}

/* Updates the various item auction status. This method is a lil ugly. We need separate methods
to update the fields independently. */
func (gandalf *Gandalf) UpdateItemAuctionStatus(itemID string, auctionStartedStatus bool, auctionEndedStatus bool, auctionDecidedStatus bool) error {
	tx := gandalf.Db.Begin()
	var item ItemModel
	dbc := tx.Where("item_id = ?", itemID).First(&item)
	if dbc.Error != nil {
		tx.Rollback()
		return NewGandalfError(KGandalfBackendError, dbc.Error.Error())
	}
	item.AuctionStarted = auctionStartedStatus
	item.AuctionEnded = auctionEndedStatus
	item.AuctionDecided = auctionDecidedStatus
	dbc = tx.Save(&item)
	if dbc.Error != nil {
		tx.Rollback()
		return NewGandalfError(KGandalfBackendError, dbc.Error.Error())
	}
	tx.Commit()
	return nil
}

func (gandalf *Gandalf) DeleteItem(itemID string) error {
	item := ItemModel{
		ItemID: itemID,
	}
	dbc := gandalf.Db.Model(&item).Delete(&item)
	if dbc.Error != nil {
		return NewGandalfError(KGandalfBackendError, dbc.Error.Error())
	}
	return nil
}

func (gandalf *Gandalf) AddAuctions(auctions []AuctionModel) error {
	var insertRecords []interface{}
	for ii := 0; ii < len(auctions); ii++ {
		insertRecords = append(insertRecords, auctions[ii])
	}
	err := bulk.BulkInsertWithTableName(gandalf.Db, "auction_models", insertRecords)
	if err != nil {
		return NewGandalfError(KGandalfBackendError, err.Error())
	}
	return nil
}

func (gandalf *Gandalf) RegisterBid(itemID string, userID string, bidAmount float32, bidQty uint32) error {
	timeout := 5 * time.Second
	timer1 := time.NewTimer(timeout)
	for {
		select {
		case <-timer1.C:
			glog.Errorf("Unable to register bid within timeout: %v secs for "+
				"itemID: %s", timeout, itemID)
			return NewGandalfError(KTimeout, "timed out attempting to register bid")
		default:
			err := gandalf.registerBid(itemID, userID, bidAmount, bidQty)
			if err != nil {
				ge := err.(*GandalfError)
				if ge.errorCode != KGandalfBackendError {
					return ge
				}
				continue
			}
			return nil
		}
	}
}

func (gandalf *Gandalf) registerBid(itemID string, userID string, bidAmount float32, bidQty uint32) error {
	// TODO: This is a massive transaction in the critical path. This will eventually slow down the entire auction
	// TODO: system. We should move the auctions and bids to faster key value store like redis or remote-badger.
	// Check if a bid has already been made by the user.
	tx := gandalf.Db.Begin()
	var currBid BidModel
	dbc := tx.Where("item_id = ? AND user_id = ?", itemID, userID).First(&currBid)
	if dbc.Error != nil {
		if !strings.Contains(dbc.Error.Error(), "record not found") {
			tx.Rollback()
			return NewGandalfError(KInvalidItem, dbc.Error.Error())
		}
	}
	var auction AuctionModel
	dbc = tx.Where("item_id = ?", itemID).First(&auction)
	if dbc.Error != nil {
		glog.Errorf("Cannot register bid for an expired/nonexistent auction: %s, user: %s", itemID, userID)
		tx.Rollback()
		return NewGandalfError(
			KAuctionExpired,
			fmt.Sprintf("cannot register bid for an expired auction: %s", itemID))
	}
	if bidAmount < auction.MinBid {
		tx.Rollback()
		return NewGandalfError(
			KInvalidBidAmount,
			fmt.Sprintf("cannot register bid(%f) which is smaller than min bid price(%f)",
				bidAmount, auction.MinBid))
	}
	if bidQty > auction.ItemQty || bidQty > auction.MaxBidQty {
		tx.Rollback()
		return NewGandalfError(
			KInvalidBidQuantity,
			fmt.Sprintf("cannot register bid with requested quantity: %d as it is > total item qty: %d",
				bidQty, auction.ItemQty))
	}
	if currBid.ItemID == "" {
		dbc := tx.Create(&BidModel{
			ItemID:    itemID,
			UserID:    userID,
			BidAmount: bidAmount,
			BidQty:    bidQty,
		})
		if dbc.Error != nil {
			tx.Rollback()
			return NewGandalfError(KGandalfBackendError, dbc.Error.Error())
		}
	} else {
		// Only update the current bid if the new bid is higher.
		if currBid.BidAmount >= bidAmount {
			tx.Rollback()
			return NewGandalfError(
				KInvalidBidAmount,
				fmt.Sprintf("user current bid(%f) is lower than user previous bid(%f)",
					bidAmount, currBid.BidAmount))
		}
		currBid.BidAmount = bidAmount
		currBid.BidQty = bidQty
		dbc := tx.Model(&currBid).Where("item_id = ? AND user_id = ?", itemID, userID).Updates(&currBid)
		if dbc.Error != nil {
			tx.Rollback()
			return NewGandalfError(KGandalfBackendError, dbc.Error.Error())
		}
	}

	// Update the max bid.
	if auction.MaxBid < bidAmount {
		auction.MaxBid = bidAmount
		dbc := tx.Model(&auction).Where("item_id = ?", itemID).Updates(&auction)
		if dbc.Error != nil {
			tx.Rollback()
			return NewGandalfError(KGandalfBackendError, dbc.Error.Error())
		}
	}
	glog.V(2).Infof("Successfully registered bid for auction: %s by user: %s", itemID, userID)
	tx.Commit()
	return nil
}

/* Returns all the auctions/items that the user has bid on. */
func (gandalf *Gandalf) GetUserBids(userID string) ([]BidModel, error) {
	var bids []BidModel
	dbc := gandalf.Db.Where("user_id = ?", userID).Find(&bids)
	if dbc.Error != nil {
		return bids, NewGandalfError(KGandalfBackendError, dbc.Error.Error())
	}
	return bids, nil
}

/* Returns all the auctions/items that the user has bid on. */
func (gandalf *Gandalf) GetUserBid(userID string, itemID string) (BidModel, error) {
	var bid BidModel
	dbc := gandalf.Db.Where("user_id = ? and item_id = ?", userID, itemID).First(&bid)
	if dbc.Error != nil {
		if strings.Contains(dbc.Error.Error(), "record not found") {
			return bid, nil
		}
		return bid, NewGandalfError(KGandalfBackendError, dbc.Error.Error())
	}
	return bid, nil
}

/* Returns the bids for a given item starting from 'start' row upto numBids rows. */
func (gandalf *Gandalf) ScanItemBids(itemID string, startIndex uint64, numBids uint64) ([]BidModel, error) {
	var bids []BidModel
	dbc := gandalf.Db.Where("item_id = ?", itemID).Offset(startIndex).Limit(numBids).Find(&bids)
	if dbc.Error != nil {
		return bids, NewGandalfError(KGandalfBackendError, dbc.Error.Error())
	}
	return bids, nil
}

/* Returns all the bids for a given item. This method must be used with care as there could potentially be
millions of such records. */
func (gandalf *Gandalf) GetAllItemBids(itemID string) ([]BidModel, error) {
	var bids []BidModel
	dbc := gandalf.Db.Where("item_id = ?", itemID).Find(&bids)
	if dbc.Error != nil {
		return bids, NewGandalfError(KGandalfBackendError, dbc.Error.Error())
	}
	return bids, nil
}

/* Fetches the max bids for the given items. */
func (gandalf *Gandalf) GetMaxBids(itemIDs []string) ([]AuctionModel, error) {
	var auctions []AuctionModel
	dbc := gandalf.Db.Where("item_id IN (?)", itemIDs).Find(&auctions)
	if dbc.Error != nil {
		return auctions, NewGandalfError(KGandalfBackendError, dbc.Error.Error())
	}
	return auctions, nil
}

/* Fetches all the auctions starting from start index upto numAuctions. */
func (gandalf *Gandalf) FetchAuctions(startIndex uint64, numAuctions uint64) ([]AuctionModel, error) {
	var auctions []AuctionModel
	dbc := gandalf.Db.Offset(startIndex).Limit(numAuctions).Find(&auctions)
	if dbc.Error != nil {
		return auctions, NewGandalfError(KGandalfBackendError, dbc.Error.Error())
	}
	return auctions, nil
}

/* Adds the given orders to the backend. */
func (gandalf *Gandalf) AddOrders(orders []OrderModel) error {
	// Adds the given orders to the database.
	var dbc *gorm.DB
	var err error
	for _, order := range orders {
		err = nil
		for ii := 0; ii < 5; ii++ {
			order.OrderID = xid.New().String()
			// All new orders will start from KOrderPaymentPending
			order.Status = KOrderPaymentPending
			dbc = gandalf.Db.Create(&order)
			if dbc.Error != nil {
				// TODO: Check if it is a unique ID error before retrying.
				// Retry with a new order ID.
				err = dbc.Error
				continue
			} else {
				break
			}
		}
		if err != nil {
			glog.Errorf("Unable to add order: %v to backend due to err: %v", order, err)
			return NewGandalfError(KGandalfBackendError, err.Error())
		}
	}
	return nil
}

/* Gets all user orders. */
func (gandalf *Gandalf) GetUserOrders(userID string) ([]OrderModel, error) {
	var orders []OrderModel
	dbc := gandalf.Db.Where("user_id = ?", userID).Order("created_date DESC").Find(&orders)
	if dbc.Error != nil {
		return orders, NewGandalfError(KGandalfBackendError, dbc.Error.Error())
	}
	return orders, nil
}

/*
Gets all orders whose payment is pending. This method must be used with care as it could potentially end up
returning 100,000 of rows.
*/
func (gandalf *Gandalf) ScanPaymentPendingOrders(startIndex uint64, numOrders uint64) ([]OrderModel, error) {
	var orders []OrderModel
	dbc := gandalf.Db.Where(
		"status = ?", KOrderPaymentPending).Offset(startIndex).Limit(numOrders).Find(&orders)
	if dbc.Error != nil {
		return orders, NewGandalfError(KGandalfBackendError, dbc.Error.Error())
	}
	return orders, nil
}

/* Gets all orders whose delivery is pending. */
func (gandalf *Gandalf) ScanDeliveryPendingOrders(startIndex uint64, numOrders uint64) ([]OrderModel, error) {
	var orders []OrderModel
	dbc := gandalf.Db.Where(
		"status = ?", KOrderDeliveryPending).Offset(startIndex).Limit(numOrders).Find(&orders)
	if dbc.Error != nil {
		return orders, NewGandalfError(KGandalfBackendError, dbc.Error.Error())
	}
	return orders, nil
}

/* Get order information */
func (gandalf *Gandalf) GetOrder(orderID string) (OrderModel, error) {
	var order OrderModel
	dbc := gandalf.Db.Where("order_id = ?", orderID).First(&order)
	if dbc.Error != nil {
		return order, NewGandalfError(KGandalfBackendError, dbc.Error.Error())
	}
	return order, nil
}

/* Update order status. */
func (gandalf *Gandalf) UpdateOrderStatus(orderID string, status uint32) error {
	tx := gandalf.Db.Begin()
	var order OrderModel
	dbc := tx.Where("order_id = ?", orderID).First(&order)
	if dbc.Error != nil {
		if !strings.Contains(dbc.Error.Error(), "record not found") {
			tx.Rollback()
			return NewGandalfError(KGandalfBackendError, dbc.Error.Error())
		}
	}
	currTime := time.Now()
	order.Status = status
	order.UpdatedDate = currTime
	if len(order.OrderHistory) > 0 {
		var events []OrderEvent
		err := json.Unmarshal(order.OrderHistory, &events)
		if err != nil {
			glog.Errorf("Unable to update order history due to err: %v\nOrder History: %v", err, order.OrderHistory)
		}
		events = append(events, OrderEvent{
			Date:   currTime,
			Status: order.Status,
			Msg:    "Updated order status",
		})
		order.OrderHistory, err = json.Marshal(events)
		if err != nil {
			glog.Errorf("Unable to marshal order history due to err: %v\nOrder History: %v", err, order.OrderHistory)
		}
	}

	dbc = tx.Save(&order)
	if dbc.Error != nil {
		tx.Rollback()
		return NewGandalfError(KGandalfBackendError, dbc.Error.Error())
	}
	tx.Commit()
	return nil
}
