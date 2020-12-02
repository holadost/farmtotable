package legolas

import (
	"farmtotable/gandalf"
	"log"
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

func prepareDB(t *testing.T, gnd *gandalf.Gandalf, numItems int) {
	now := time.Now()
	err := gnd.RegisterUser("nikhil", "Nikhil Srivatsan", "kjahd@lkaj.com",
		"9873981799", "khadkjhadkha")
	if err != nil {
		t.Fatalf("Unable to register user")
	}
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
	supplier := suppliers[0]
	for ii := 0; ii < numItems; ii++ {
		itemName := "Item_" + strconv.Itoa(ii)
		if ii%2 == 0 {
			// These items are ready to be auctioned.
			err = gnd.RegisterItem(supplier.SupplierID, itemName,
				"lkjadlkjadl", uint32((ii+1)*100), now, float32((ii+1)*10))
		} else {
			// These items are not yet ready to be auctioned.
			err = gnd.RegisterItem(supplier.SupplierID, itemName,
				"lkjadlkjadl", uint32((ii+1)*100),
				now.Add(time.Duration(time.Second*900)),
				float32((ii+1)*10))
		}
		if err != nil {
			t.Fatalf("Unable to register item: %s", itemName)
		}
	}

}

func TestScanners(t *testing.T) {
	cleanupDB()
	gnd := gandalf.NewSqliteGandalf()
	numItems := 5
	prepareDB(t, gnd, numItems)
	itemScanner := gandalf.NewItemsScanner(gnd, 3)
	var items []gandalf.ItemModel
	for {
		item, finished, err := itemScanner.Next()
		if err != nil {
			t.Fatalf("Unable to scan items due to err: %v", err)
		}
		if finished {
			if len(items) != numItems {
				t.Fatalf("Did not scan all items")
			}
			break
		}
		items = append(items, item)
	}
}
