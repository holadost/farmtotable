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

func TestGandalf_Item(t *testing.T) {
	cleanupSqliteDB()
	gandalf := NewSqliteGandalf()
	defer gandalf.Close()
	err := gandalf.RegisterItem("nikhil_srivatsan", "Item1", "This stuff is good", 100, time.Now(), 100.0)
	if err != nil {
		t.Fatalf("Unable to register item 1")
	}
	err = gandalf.RegisterItem("nikhil_srivatsan", "Item2", "This stuff is good 2", 200, time.Now(), 66.66)
	if err != nil {
		t.Fatalf("Unable to register item 2")
	}
	err = gandalf.RegisterItem("nikhil_srivatsan", "Item3", "This stuff is good 3", 300, time.Now(), 33.33)
	if err != nil {
		t.Fatalf("Unable to register item 3")
	}
	items, err := gandalf.GetUserItems("nikhil_srivatsan")
	if err != nil {
		t.Fatalf("Unable to get user items")
	}
	if len(items) != 3 {
		t.Fatalf("Did not get the correct number of items")
	}

	item := gandalf.GetItem(items[0].ItemID)
	if item.ItemName != "Item1" {
		t.Fatalf("Did not get the correct item")
	}

	err = gandalf.EditItem(item.ItemID, "Item11", item.ItemDescription, 150, item.AuctionStartTime, item.MinPrice)
	if err != nil {
		t.Fatalf("Unable to edit item. Error: %v", err)
	}

	item = gandalf.GetItem(items[0].ItemID)
	if item.ItemName != "Item11" || item.ItemQty != 150 {
		t.Fatalf("Did not get the correctly edited item")
	}

	err = gandalf.DeleteItem(items[1].ItemID)
	if err != nil {
		t.Fatalf("Unable to delete item")
	}

	item = gandalf.GetItem(items[1].ItemID)
	if item.ItemID != "" {
		t.Fatalf("Deleted item came back")
	}

	items, err = gandalf.GetUserItems("nikhil_srivatsan")
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
	var auctions []Auction
	for ii := 0; ii < 5; ii++ {
		itemName := "Item" + strconv.Itoa(ii)
		itemDesc := itemName + ": Item description"
		err := gandalf.RegisterItem("nikhil", itemName, itemDesc, uint32(100*(ii+1)), time.Now(), float32(1.0*ii))
		if err != nil {
			t.Fatalf("Unable to register item")
		}
	}
	items, err := gandalf.GetUserItems("nikhil")
	if err != nil || len(items) != 5 {
		t.Fatalf("Unable to fetch items for user")
	}
	for ii := 0; ii < 5; ii++ {
		auctions = append(auctions, Auction{
			ItemID:              items[ii].ItemID,
			ItemQty:             items[ii].ItemQty,
			AuctionStartTime:    items[ii].AuctionStartTime,
			AuctionDurationSecs: 24 * 86400,
			MaxBid:              items[ii].MinPrice,
		})
	}
	err = gandalf.RegisterAuctions(auctions)
	if err != nil {
		t.Fatalf("Unable to register auction")
	}

	mainAucs, err := gandalf.GetAllAuctions(0, 5)
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
