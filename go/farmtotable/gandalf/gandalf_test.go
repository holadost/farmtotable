package gandalf

import (
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

/* Deletes the underlying sqlite database. */
func cleanupSqliteDB() {
	_, err := os.Stat(SQLiteDBPath)
	if os.IsNotExist(err) {
		return
	}
	err = os.Remove(SQLiteDBPath)
	if err != nil {
		log.Fatalf("Unable to delete sqlite db")
	}
}

func TestNewSqliteGandalf(t *testing.T) {
	cleanupSqliteDB()
	gandalf := NewSqliteGandalf()
	defer gandalf.Close()
}

func TestNewPostgresGandalf(t *testing.T) {

}

func TestGandalf_User(t *testing.T) {
	cleanupSqliteDB()
	gandalf := NewSqliteGandalf()
	defer gandalf.Close()
	err := gandalf.RegisterUser("nikhil_srivatsan", "Nikhil Srivtsan", "ssnikhil87@gmail.com", "9198029973", "blahblahblach")
	if err != nil {
		t.Fatalf("Unable to register user")
	}
	err = gandalf.RegisterUser("raunaq_naidu", "Raunaq Naidu", "rbnaidu@gmail.com", "9198029972", "blahblahblach")
	if err != nil {
		t.Fatalf("Unable to register user")
	}
	err = gandalf.RegisterUser("rohit_pagedar", "Rohit Pagedar", "rgpagedar@gmail.com", "9198029971", "blahblahblach")
	if err != nil {
		t.Fatalf("Unable to register user")
	}
	user := gandalf.GetUserByID("nikhil_srivatsan")
	if user.EmailID != "ssnikhil87@gmail.com" {
		t.Fatalf("Failed to insert and fetch records")
	}
	user = gandalf.GetUserByEmailID("rbnaidu@gmail.com")
	if user.UserID != "raunaq_naidu" {
		t.Fatalf("Failed to insert and fetch records")
	}
	user = gandalf.GetUserByPhNo("9198029971")
	if user.UserID != "rohit_pagedar" {
		t.Fatalf("Failed to insert and fetch records")
	}
}

func TestGandalf_Supplier(t *testing.T) {
	cleanupSqliteDB()
	gandalf := NewSqliteGandalf()
	defer gandalf.Close()
	err := gandalf.RegisterSupplier("nikhil_srivatsan", "nikhil.sriniva@nutanix.com",
		"9198029973", "Mera Naam Joker", "This supplier is a god amongst humans",
		"tag1,tag2,tag3")
	if err != nil {
		t.Fatalf("Unable to register supplier")
	}

	err = gandalf.RegisterSupplier("rohit_srinivasan", "rohit.sriniva@nutanix.com",
		"9198029973", "Mera Naam Joker", "This supplier is a god amongst humans",
		"tag1,tag2,tag3")
	if err != nil {
		t.Fatalf("Unable to register supplier")
	}

	suppliers, err := gandalf.GetAllSuppliers()
	if err != nil {
		t.Fatalf("Unable to fetch all suppliers. Error: %v", err)
	}
	if len(suppliers) != 2 {
		t.Fatalf("Did not get the correct number of suppliers")
	}

	supplier := gandalf.GetSupplierByID(suppliers[0].SupplierID)
	if supplier.SupplierName != suppliers[0].SupplierName {
		t.Fatalf("Fetched the wrong record")
	}
}

func TestGandalf_Item(t *testing.T) {
	cleanupSqliteDB()
	gandalf := NewSqliteGandalf()
	defer gandalf.Close()
	err := gandalf.RegisterItem(
		"nikhil_srivatsan", "Item1",
		"This stuff is good",
		100, time.Now(), 100.0,
		100, "http://imageURL",
		100, 5000, "g")
	if err != nil {
		t.Fatalf("Unable to register item 1")
	}
	err = gandalf.RegisterItem("nikhil_srivatsan", "Item2", "This stuff is good 2", 200, time.Now(), 66.66, 100, "http://imageURL", 100, 5000, "g")
	if err != nil {
		t.Fatalf("Unable to register item 2")
	}
	err = gandalf.RegisterItem("nikhil_srivatsan", "Item3", "This stuff is good 3", 300, time.Now(), 33.33, 100, "http://imageURL", 100, 5000, "g")
	if err != nil {
		t.Fatalf("Unable to register item 3")
	}
	items, err := gandalf.GetSupplierItems("nikhil_srivatsan")
	if err != nil {
		t.Fatalf("Unable to get user items")
	}
	if len(items) != 3 {
		t.Fatalf("Did not get the correct number of items")
	}

	item, err := gandalf.GetItem(items[0].ItemID)
	if err != nil {
		t.Fatalf("Unable to get item. Error: %v", err)
	}
	if item.ItemName != "Item1" {
		t.Fatalf("Did not get the correct item")
	}

	err = gandalf.EditItem(item.ItemID, "Item11", item.ItemDescription, 150, item.AuctionStartTime, item.MinPrice)
	if err != nil {
		t.Fatalf("Unable to edit item. Error: %v", err)
	}

	item, err = gandalf.GetItem(items[0].ItemID)
	if err != nil {
		t.Fatalf("Unable to get item. Error: %v", err)
	}
	if item.ItemName != "Item11" || item.ItemQty != 150 {
		t.Fatalf("Did not get the correctly edited item")
	}

	err = gandalf.DeleteItem(items[1].ItemID)
	if err != nil {
		t.Fatalf("Unable to delete item")
	}

	item, err = gandalf.GetItem(items[1].ItemID)
	if item.ItemID != "" {
		t.Fatalf("Deleted item came back")
	}

	items, err = gandalf.GetSupplierItems("nikhil_srivatsan")
	if err != nil {
		t.Fatalf("Unable to get user items")
	}
	if len(items) != 2 {
		t.Fatalf("Did not get the correct number of items")
	}
}

func TestGandalf_Auction(t *testing.T) {
	cleanupSqliteDB()
	gandalf := NewSqliteGandalf()
	defer gandalf.Close()
	var auctions []AuctionModel
	for ii := 0; ii < 5; ii++ {
		itemName := "ItemModel" + strconv.Itoa(ii)
		itemDesc := itemName + ": ItemModel description"
		err := gandalf.RegisterItem(
			"nikhil", itemName, itemDesc, uint32(100*(ii+1)),
			time.Now(), float32(1.0*ii), 100,
			"imageUrl", 100, 5000, "g")
		if err != nil {
			t.Fatalf("Unable to register item")
		}
	}
	items, err := gandalf.GetSupplierItems("nikhil")
	if err != nil || len(items) != 5 {
		t.Fatalf("Unable to fetch items for user")
	}
	for ii := 0; ii < 5; ii++ {
		auctions = append(auctions, AuctionModel{
			ItemID:              items[ii].ItemID,
			ItemQty:             items[ii].ItemQty,
			AuctionStartTime:    items[ii].AuctionStartTime,
			AuctionDurationSecs: 24 * 86400,
			MaxBid:              items[ii].MinPrice,
			MinBidQty:           100,
			MaxBidQty:           1000,
			ItemUnit:            "g",
		})
	}
	err = gandalf.AddAuctions(auctions)
	if err != nil {
		t.Fatalf("Unable to register auction")
	}

	mainAucs, err := gandalf.FetchAuctions(0, 5)
	if err != nil || len(mainAucs) != 5 {
		t.Fatalf("Unable to fetch auctions")
	}

	newMaxBid := items[0].MinPrice + 1.0
	err = gandalf.RegisterBid(items[0].ItemID, "raunaq", newMaxBid, 10)
	if err != nil {
		t.Fatalf("Unable to register bid. Error: %v", err)
	}

	var itemIDs []string
	itemIDs = append(itemIDs, items[0].ItemID)
	mainAucs, err = gandalf.GetMaxBids(itemIDs)
	if err != nil || len(mainAucs) != 1 {
		t.Fatalf("Unable to fetch max bids for required items")
	}
	if mainAucs[0].MaxBid != newMaxBid {
		t.Fatalf("Max bid did not get updated as expected")
	}

	newMaxBid = items[1].MinPrice - 1.0
	err = gandalf.RegisterBid(items[1].ItemID, "raunaq", newMaxBid, 10)
	if err != nil {
		t.Fatalf("Unable to register bid. Error: %v", err)
	}

	itemIDs = itemIDs[:0]
	itemIDs = append(itemIDs, items[1].ItemID)
	mainAucs, err = gandalf.GetMaxBids(itemIDs)
	if err != nil || len(mainAucs) != 1 {
		t.Fatalf("Unable to fetch max bids for required items")
	}
	if mainAucs[0].MaxBid != items[1].MinPrice {
		t.Fatalf("Max bid got updated even though it should not have")
	}
}

func TestGandalf_Order(t *testing.T) {
	cleanupSqliteDB()
	gandalf := NewSqliteGandalf()
	defer gandalf.Close()
	var orders []OrderModel
	numOrders := 5
	for ii := 0; ii < numOrders; ii++ {
		var order OrderModel
		order.ItemQty = 10
		order.ItemPrice = 10.0
		order.UserID = "nikhil"
		order.ItemID = "Item_" + strconv.Itoa(ii)
		orders = append(orders, order)
	}
	if err := gandalf.AddOrders(orders); err != nil {
		t.Fatalf("Unable to add errors due to err: %v", err)
	}
	retOrders, err := gandalf.GetUserOrders("nikhil")
	if err != nil {
		t.Fatalf("Unable to fetch orders for nikhil due to err: %v", err)
	}
	if len(retOrders) != len(orders) {
		t.Fatalf("Did not get all the orders as expected(%d vs %d)",
			len(retOrders), len(orders))
	}
	err = gandalf.UpdateOrderStatus(retOrders[0].OrderID, KOrderDeliveryPending)
	if err != nil {
		t.Fatalf("Unable to update order status for order: %v due to err: %v",
			retOrders[0], err)
	}
	order, err := gandalf.GetOrder(retOrders[0].OrderID)
	if err != nil {
		t.Fatalf("Unable to get order: %s due to err: %v", retOrders[0].OrderID, err)
	}
	if order.Status != KOrderDeliveryPending {
		t.Fatalf("Invalid order status. Expected: %d, got: %d",
			KOrderDeliveryPending, order.Status)
	}
	scanOrders, err := gandalf.ScanPaymentPendingOrders(0, uint64(numOrders))
	if err != nil {
		t.Fatalf("Unable to scan payment pending orders due to err: %v", err)
	}
	if len(scanOrders) != (numOrders - 1) {
		t.Fatalf("Expected %d orders, got %d. Orders: %v", numOrders-1, len(scanOrders), scanOrders)
	}
	scanOrders, err = gandalf.ScanDeliveryPendingOrders(0, uint64(numOrders))
	if err != nil {
		t.Fatalf("Unable to scan delivery pending orders due to err: %v", err)
	}
	if len(scanOrders) != 1 {
		t.Fatalf("Expected 1 order, got %d", len(scanOrders))
	}
}
