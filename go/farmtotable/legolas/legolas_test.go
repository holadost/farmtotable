package legolas

import (
	"farmtotable/gandalf"
	"farmtotable/legolas/auction_winner_job"
	"farmtotable/legolas/new_auction_job"
	"github.com/golang/glog"
	"log"
	"math"
	"os"
	"strconv"
	"testing"
	"time"
)

/* Deletes the underlying sqlite database. */
func cleanupDB() {
	_, err := os.Stat(gandalf.SQLiteDBPath)
	if os.IsNotExist(err) {
		return
	}
	err = os.Remove(gandalf.SQLiteDBPath)
	if err != nil {
		log.Fatalf("Unable to delete sqlite db")
	}
}

func prepareDB(t *testing.T, gnd *gandalf.Gandalf, userID string, numItems int, numBids int, auctionDurationSecs uint32) {
	now := time.Now()
	// Add user.
	err := gnd.RegisterUser(userID, "Nikhil Srivatsan", "kjahd@lkaj.com",
		"9873981799", "khadkjhadkha")
	if err != nil {
		t.Fatalf("Unable to register user")
	}

	// Add supplier.
	err = gnd.RegisterSupplier(
		"supplier", "supplier@supplier.com",
		"987134987139", "kjahdkjhadkjhadkh",
		"kjhadkjhadkjhadkjhakjdhakjhdkjahdkjhadkjhak",
		"klahjd,kahd,alkhd")
	if err != nil {
		t.Fatalf("Unable to register supplier")
	}
	suppliers, err := gnd.GetAllSuppliers()
	if err != nil {
		t.Fatalf("Unable to fetch suppliers")
	}
	if len(suppliers) != 1 {
		t.Fatalf("Did not find the expected number of suppliers. Expected: 1, got: %d",
			len(suppliers))
	}

	// Add items.
	supplier := suppliers[0]
	for ii := 0; ii < numItems; ii++ {
		itemName := "Item_" + strconv.Itoa(ii)
		if ii%2 == 0 {
			// These items are ready to be auctioned.
			err = gnd.RegisterItem(supplier.SupplierID, itemName,
				"lkjadlkjadl", uint32((ii+1)*100), now, float32((ii+1)*2),
				auctionDurationSecs)
		} else {
			// These items are not yet ready to be auctioned.
			err = gnd.RegisterItem(supplier.SupplierID, itemName,
				"lkjadlkjadl", uint32((ii+1)*100),
				now.Add(time.Duration(time.Second*900)),
				float32((ii+1)*10), auctionDurationSecs)
		}
		if err != nil {
			t.Fatalf("Unable to register item: %s", itemName)
		}
	}

	// Add bids.
	addBids(t, gnd, userID, uint(numBids), numItems*2)
}

func addBids(t *testing.T, gnd *gandalf.Gandalf, userID string, numBids uint, minBidPrice int) {
	suppliers, err := gnd.GetAllSuppliers()
	if err != nil {
		t.Fatalf("Unable to fetch suppliers")
	}

	// Add items.
	supplier := suppliers[0]
	items, err := gnd.GetSupplierItems(supplier.SupplierID)
	if err != nil {
		t.Fatalf("Unable to get supplier items due to err: %v", err)
	}
	nextID := 0
	for ii := 0; ii < int(numBids); ii++ {
		itemID := items[nextID].ItemID
		nextID += 2
		if nextID >= len(items) {
			nextID = 0
		}
		err = gnd.RegisterBid(itemID, userID, float32(minBidPrice+ii), 3)
		if err != nil {
			t.Fatalf("Unable to register bid due to err: %v", err)
		}
	}
}

func TestNewAuctionsJob(t *testing.T) {
	cleanupDB()
	gnd := gandalf.NewSqliteGandalf()
	defer gnd.Close()
	numItems := 15
	expectedAucs := int(math.Ceil(float64(numItems) / 2))
	numBids := 15
	prepareDB(t, gnd, "nikhil", numItems, 0, 300)
	naj := new_auction_job.NewPopulateNewAuctionsJob(gnd, 5, 2)
	naj.Run()
	auctionScanner := gandalf.NewAuctionsScanner(gnd, 5)
	suppliers, err := gnd.GetAllSuppliers()
	if err != nil {
		t.Fatalf("Unable to fetch suppliers")
	}
	// Add bids.
	supplier := suppliers[0]
	items, err := gnd.GetSupplierItems(supplier.SupplierID)
	if err != nil {
		t.Fatalf("Unable to get supplier items due to err: %v", err)
	}
	nextID := 0
	for ii := 0; ii < numBids; ii++ {
		itemID := items[nextID].ItemID
		nextID += 2
		if nextID >= len(items) {
			nextID = 0
		}
		err = gnd.RegisterBid(itemID, "nikhil", float32(30+ii), 3)
		if err != nil {
			t.Fatalf("Unable to register bid due to err: %v", err)
		}
	}
	// Testing auction scanner.
	var auctions []gandalf.AuctionModel
	counter := 0
	for {
		var scannedAuctions []gandalf.AuctionModel
		var finished bool
		var err error
		if counter%2 == 0 {
			scannedAuctions, finished, err = auctionScanner.NextBatch()
		} else if counter%2 == 1 {
			for jj := 0; jj < 2; jj++ {
				var scannedAuction gandalf.AuctionModel
				scannedAuction, finished, err = auctionScanner.Next()
				if err != nil {
					t.Fatalf("Unable to get next auction due to err: %v", err)
				}
				scannedAuctions = append(scannedAuctions, scannedAuction)
			}
		}
		if err != nil {
			t.Fatalf("Unable to scan auctions due to err: %v", err)
		}
		for _, auction := range scannedAuctions {
			if auction.ItemID != "" {
				auctions = append(auctions, auction)
			}
		}
		if finished {

			if len(auctions) != expectedAucs {
				t.Fatalf("Did not scan all auctions. Expected: %d, got: %d", expectedAucs, len(auctions))
			}
			break
		}
		counter += 1
	}
}

func TestAuctionWinnerJob(t *testing.T) {
	cleanupDB()
	gnd := gandalf.NewSqliteGandalf()
	defer gnd.Close()
	numItems := 15
	numBids := 15
	durationSecs := uint32(30)
	prepareDB(t, gnd, "nikhil", numItems, 0, durationSecs)
	naj := new_auction_job.NewPopulateNewAuctionsJob(gnd, 5, 2)
	naj.Run()
	addBids(t, gnd, "nikhil", uint(numBids), numItems*2)
	awj := auction_winner_job.NewAuctionWinnerJob(gnd, 2, 4)
	awj.Run()
	glog.Infof("Sleeping for %d seconds to allow auctions to expire", durationSecs)
	time.Sleep(time.Second * time.Duration(durationSecs))
	// Pick the winners.
	awj = auction_winner_job.NewAuctionWinnerJob(gnd, 2, 4)
	awj.Run()

	// Run job again and ensure that orders are not placed again.
	awj = auction_winner_job.NewAuctionWinnerJob(gnd, 2, 4)
	awj.Run()
}
