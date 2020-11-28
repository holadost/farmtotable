package aragorn

import (
	"bytes"
	"encoding/json"
	"farmtotable/gandalf"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

func startAragorn() *Aragorn {
	aragorn := NewAragorn()
	glog.Info("Starting aragorn service in a background go routine")
	go aragorn.Run()
	return aragorn
}

func cleanupSqliteDB() {
	SQLiteDBPath := "/tmp/gandalf.db"
	_, err := os.Stat(SQLiteDBPath)
	if os.IsNotExist(err) {
		return
	}
	err = os.Remove(SQLiteDBPath)
	if err != nil {
		log.Fatalf("Unable to delete sqlite db")
	}
}

func TestAragornRun(t *testing.T) {
	cleanupSqliteDB()
	aragorn := startAragorn()
	time.Sleep(100 * time.Millisecond)
	baseURL := "http://localhost:8080/api/v1/resources"
	/*************************** Users *************************/
	glog.Info("Testing aragorn users APIs")
	// Register new user.
	userArg := RegisterUserArg{}
	userArg.UserID = "nikhil.sriniva"
	userArg.EmailID = "nikhil.sriniva@nutanix.com"
	userArg.PhNum = "9198029973"
	userArg.Address = "840 Meridian Way, San Jose 95126"
	userArg.Name = "Nikhil Srivatsan Srinivasan"
	body, err := json.Marshal(&userArg)
	if err != nil {
		t.Fatalf("Failed to marshal JSON. Error: %v", err)
	}
	resp, err := http.Post(baseURL+"/users/register", "application/json", bytes.NewBuffer(body))
	if err != nil || resp.StatusCode >= 300 {
		t.Fatalf("Unable to register user. Error: %v", err)
	}
	fullBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	ret := RegisterUserRet{}
	err = json.Unmarshal(fullBody, &ret)
	if err != nil {
		t.Fatalf("Failed to unmarshall response JSON after registering user. Error: %v", err)
	}
	if ret.Status >= 300 {
		t.Fatalf("Error while registering user. Status: %d", ret.Status)
	}
	if !ret.Data.RegistrationStatus {
		t.Fatalf("Unable to register user even though status code is < 300??. Response: %v", ret)
	}

	// Get user.
	getUserArg := GetUserArg{
		UserID: "nikhil.sriniva",
	}
	body, err = json.Marshal(&getUserArg)
	if err != nil {
		t.Fatalf("Failed to marshal JSON. Error: %v", err)
	}

	resp, err = http.Post(fmt.Sprintf(baseURL+"/users/fetch"), "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Unable to get user. Error: %v", err)
	}
	fullResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body")
	}
	userRet := GetUserRet{}
	err = json.Unmarshal(fullResp, &userRet)
	if userRet.Status != http.StatusOK {
		t.Fatalf("Failed to get user. Error: %d. Message: %s", userRet.Status, userRet.ErrorMsg)
	}
	if userRet.Data.UserID != "nikhil.sriniva" {
		t.Fatalf("Fetched wrong user")
	}
	if userRet.Data.PhNum != "9198029973" {
		t.Fatalf("User ph num is wrong")
	}

	/************************************* Suppliers *******************************************/
	glog.Info("Testing aragorn suppliers APIs")
	// Register supplier
	supplierArg := RegisterSupplierArg{
		SupplierName:        "Supplier 1",
		SupplierTags:        "Tag1, Tag2, Tag3",
		SupplierDescription: "This is a BS supplier",
		SupplierAddress:     "Tera Ghar",
		SupplierPhNum:       "0001112223",
		SupplierEmailID:     "teraghar@meraghar.com",
	}
	body, err = json.Marshal(&supplierArg)
	if err != nil {
		t.Fatalf("Unable to marshal supplier args to JSON. Error: %v", err)
	}
	resp, err = http.Post(baseURL+"/suppliers/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Unable to register supplier. Error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Error while registering supplier. Error Code: %d", resp.StatusCode)
	}
	fullBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unable to read body")
	}
	resp.Body.Close()
	supplierRet := RegisterSupplierRet{}
	err = json.Unmarshal(fullBody, &supplierRet)
	if supplierRet.Status != http.StatusOK {
		t.Fatalf("Unable to register supplier. Error Code: %d, Error Message: %s", supplierRet.Status, supplierRet.ErrorMsg)
	}

	// Get All suppliers.
	allSuppArg := GetAllSuppliersArg{}
	body, err = json.Marshal(allSuppArg)
	if err != nil {
		t.Fatalf("Unable to marshal get all suppliers arg")
	}
	resp, err = http.Post(baseURL+"/suppliers/fetch_all", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Unable to get all suppliers. Error: %v", err)
	}
	fullBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unable to read body")
	}
	resp.Body.Close()
	allSuppRet := GetAllSuppliersRet{}
	err = json.Unmarshal(fullBody, &allSuppRet)
	if err != nil {
		t.Fatalf("Unable to deserialize suppliers ret. Error: %v", err)
	}
	if allSuppRet.Status != http.StatusOK {
		t.Fatalf("Unable to register supplier. Error Code: %d, Error Message: %s", allSuppRet.Status, allSuppRet.ErrorMsg)
	}
	if len(allSuppRet.Data) != 1 {
		t.Fatalf("Unable to get all suppliers as expected")
	}

	// Get Supplier
	suppArg := GetSupplierArg{
		SupplierID: allSuppRet.Data[0].SupplierID,
	}
	body, err = json.Marshal(suppArg)
	if err != nil {
		t.Fatalf("Unable to marshal get all suppliers arg")
	}
	resp, err = http.Post(baseURL+"/suppliers/fetch", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Unable to get all suppliers. Error: %v", err)
	}
	fullBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unable to read body")
	}
	resp.Body.Close()
	suppRet := GetSupplierRet{}
	err = json.Unmarshal(fullBody, &suppRet)
	if err != nil {
		t.Fatalf("Unable to deserialize suppliers ret. Error: %v", err)
	}
	if suppRet.Status != http.StatusOK {
		t.Fatalf("Unable to register supplier. Error Code: %d, Error Message: %s", suppRet.Status, suppRet.ErrorMsg)
	}
	if suppRet.Data.SupplierID != allSuppRet.Data[0].SupplierID {
		t.Fatalf("Failure while fetching required supplier")
	}

	/************************** Items *********************************/
	glog.Info("Testing aragorn items APIs")
	// Register Item
	regItemArg := RegisterItemArg{
		ItemName:         "Item 1",
		ItemQty:          100,
		ItemDescription:  "Some stupid item.",
		ItemTags:         "Tag1, Tag2, Tag3",
		SupplierID:       "Supplier 1",
		AuctionStartDate: time.Now(),
		MinPrice:         5.0,
	}
	body, err = json.Marshal(regItemArg)
	if err != nil {
		t.Fatalf("Unable to marshal register item arg")
	}
	resp, err = http.Post(baseURL+"/items/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Unable to register item. Error: %v", err)
	}
	fullBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unable to read body")
	}
	resp.Body.Close()
	regItemRet := RegisterItemRet{}
	err = json.Unmarshal(fullBody, &regItemRet)
	if err != nil {
		t.Fatalf("Unable to deserialize reg items ret. Error: %v", err)
	}
	if regItemRet.Status != http.StatusOK {
		t.Fatalf("Unable to register item. Error Code: %d, Error Message: %s", regItemRet.Status, regItemRet.ErrorMsg)
	}
	if !regItemRet.Data.RegistrationStatus {
		t.Fatalf("Failure while registering item")
	}

	// Get supplier items.
	getSupplierItemsArg := GetSupplierItemsArg{}
	getSupplierItemsArg.SupplierID = "Supplier 1"
	body, err = json.Marshal(getSupplierItemsArg)
	if err != nil {
		t.Fatalf("Unable to marshal get suppliter items arg")
	}
	resp, err = http.Post(baseURL+"/items/fetch", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Unable to get supplier items. Error: %v", err)
	}
	fullBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unable to read body")
	}
	resp.Body.Close()
	getSuppItemsRet := GetSupplierItemsRet{}
	err = json.Unmarshal(fullBody, &getSuppItemsRet)
	if err != nil {
		t.Fatalf("Unable to deserialize get supp items ret. Error: %v", err)
	}
	if getSuppItemsRet.Status != http.StatusOK {
		t.Fatalf("Unable to get supp items. Error Code: %d, Error Message: %s", getSuppItemsRet.Status,
			getSuppItemsRet.ErrorMsg)
	}
	if len(getSuppItemsRet.Data) != 1 {
		t.Fatalf("Failure while fetching all supplier items")
	}

	// Remove Item
	removeItemArg := RemoveItemArg{}
	removeItemArg.ItemID = getSuppItemsRet.Data[0].ItemID
	body, err = json.Marshal(removeItemArg)
	if err != nil {
		t.Fatalf("Unable to marshal remove items arg")
	}
	resp, err = http.Post(baseURL+"/items/remove", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Unable to remove items. Error: %v", err)
	}
	fullBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unable to read body")
	}
	resp.Body.Close()
	removeItemRet := RemoveItemRet{}
	err = json.Unmarshal(fullBody, &removeItemRet)
	if err != nil {
		t.Fatalf("Unable to deserialize remove items ret. Error: %v", err)
	}
	if getSuppItemsRet.Status != http.StatusOK {
		t.Fatalf("Unable to remove items. Error Code: %d, Error Message: %s", removeItemRet.Status,
			removeItemRet.ErrorMsg)
	}
	if !removeItemRet.Data.RegistrationStatus {
		t.Fatalf("Unable to delete item")
	}

	/********************** Auctions **************************/
	glog.Info("Testing aragorn auctions APIs")
	var auctions []gandalf.Auction
	for ii := 0; ii < 5; ii++ {
		itemName := "Item" + strconv.Itoa(ii)
		itemDesc := itemName + ": Item description"
		err := aragorn.gandalf.RegisterItem("supplier1", itemName, itemDesc, uint32(100*(ii+1)), time.Now(), float32(1.0*ii))
		if err != nil {
			t.Fatalf("Unable to register item")
		}
	}
	items, err := aragorn.gandalf.GetSupplierItems("supplier1")
	if err != nil || len(items) != 5 {
		t.Fatalf("Unable to fetch items for user")
	}
	for ii := 0; ii < 5; ii++ {
		auctions = append(auctions, gandalf.Auction{
			ItemID:              items[ii].ItemID,
			ItemQty:             items[ii].ItemQty,
			AuctionStartTime:    items[ii].AuctionStartTime,
			AuctionDurationSecs: 24 * 86400,
			MaxBid:              items[ii].MinPrice,
		})
	}
	err = aragorn.gandalf.RegisterAuctions(auctions)
	if err != nil {
		t.Fatalf("Unable to register auction")
	}

	// Get All auctions
	glog.Info("Getting all auctions")
	allAucArg := FetchAllAuctionsArg{}
	allAucArg.StartID = 0
	allAucArg.NumAuctions = 5
	body, err = json.Marshal(allAucArg)
	if err != nil {
		t.Fatalf("Unable to marshal get all auctions arg")
	}
	resp, err = http.Post(baseURL+"/auctions/fetch_all", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Unable to get supplier items. Error: %v", err)
	}
	fullBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unable to read body")
	}
	resp.Body.Close()
	allAucRet := FetchAllAuctionsRet{}
	err = json.Unmarshal(fullBody, &allAucRet)
	if err != nil {
		t.Fatalf("Unable to deserialize get supp items ret. Error: %v", err)
	}
	if allAucRet.Status != http.StatusOK {
		t.Fatalf("Unable to get supp items. Error Code: %d, Error Message: %s", allAucRet.Status,
			allAucRet.ErrorMsg)
	}
	if len(allAucRet.Data.Auctions) != 5 {
		t.Fatalf("Failure while fetching all auctions")
	}

	// Register bid
	rbArg := RegisterBidArg{}
	rbArg.ItemID = items[0].ItemID
	rbArg.UserID = "nikhil.sriniva"
	rbArg.BidQty = 10
	rbArg.BidAmount = 1.5

	body, err = json.Marshal(rbArg)
	if err != nil {
		t.Fatalf("Unable to marshal register bid arg")
	}
	resp, err = http.Post(baseURL+"/auctions/register_bid", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Unable to register bid. Error: %v", err)
	}
	fullBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unable to read body")
	}
	resp.Body.Close()
	rbRet := RegisterBidRet{}
	err = json.Unmarshal(fullBody, &rbRet)
	if err != nil {
		t.Fatalf("Unable to deserialize register bid ret. Error: %v", err)
	}
	if allAucRet.Status != http.StatusOK {
		t.Fatalf("Unable to register bid. Error Code: %d, Error Message: %s", rbRet.Status,
			rbRet.ErrorMsg)
	}
	if !rbRet.Data.RegistrationStatus {
		t.Fatalf("Failure while registering bid")
	}

	// Fetch Max bids.
	fmbArg := FetchMaxBidsArg{}
	var tmpItemIds []string
	for ii := 0; ii < 3; ii++ {
		tmpItemIds = append(tmpItemIds, items[ii].ItemID)
	}
	fmbArg.ItemIDs = tmpItemIds
	body, err = json.Marshal(fmbArg)
	if err != nil {
		t.Fatalf("Unable to marshal fetch max bids arg")
	}
	resp, err = http.Post(baseURL+"/auctions/fetch_max_bids", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Unable to fetch max bids. Error: %v", err)
	}
	fullBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unable to read body")
	}
	resp.Body.Close()
	fmbRet := FetchMaxBidsRet{}
	err = json.Unmarshal(fullBody, &fmbRet)
	if err != nil {
		t.Fatalf("Unable to deserialize fetch max bids ret. Error: %v", err)
	}
	if fmbRet.Status != http.StatusOK {
		t.Fatalf("Unable to fetch max bids. Error Code: %d, Error Message: %s", fmbRet.Status,
			fmbRet.ErrorMsg)
	}
	if len(fmbRet.Data) != 3 {
		t.Fatalf("Failure while registering bid")
	}
	for ii := 0; ii < 3; ii++ {
		if ii == 0 {
			if fmbRet.Data[ii].MaxBid != 1.5 {
				t.Fatalf("Incorrect max bid for item 0")
			}
		}
		if ii == 1 {
			if fmbRet.Data[ii].MaxBid != 1.0 {
				t.Fatalf("Incorrect max bid for item 1")
			}
		}
		if ii == 2 {
			if fmbRet.Data[ii].MaxBid != 2.0 {
				t.Fatalf("Incorrect max bid for item 2")
			}
		}
	}

	/********************** Orders **************************/
	glog.Info("Testing aragorn orders APIs")
	numOrders := 5
	items, err = aragorn.gandalf.GetSupplierItems("supplier1")
	if err != nil {
		t.Fatalf("Unable to get supplier items")
	}
	for ii := 0; ii < numOrders; ii++ {
		var order TestOnlyAddOrderArg
		order.ItemPrice = 7.0 * float32(ii+1)
		order.ItemQty = uint32(5 * (ii + 1))
		if ii%2 == 0 {
			order.UserID = "nikhil_0"
			order.ItemID = items[0].ItemID
		} else {
			order.UserID = "nikhil_1"
			order.ItemID = items[1].ItemID
		}
		body, err = json.Marshal(order)
		if err != nil {
			t.Fatalf("Unable to marshal add order arg")
		}
		resp, err = http.Post(baseURL+"/test/orders/test_only_add_order", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Unable to add order. Error: %v", err)
		}
		fullBody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Unable to read add order ret")
		}
		resp.Body.Close()
		var ret TestOnlyAddOrderRet
		err = json.Unmarshal(fullBody, &ret)
		if err != nil {
			t.Fatalf("Unable to deserialize add order ret. Error: %v", err)
		}
		if ret.Status != http.StatusOK {
			t.Fatalf("Unable to add orders due to error: %v", ret.ErrorMsg)
		}
		if !ret.Data.RegistrationStatus {
			t.Fatalf("Unable to add order due to error: %v", ret.ErrorMsg)
		}
	}

	// Get user orders.
	var ordersRet GetOrdersRet
	var userOrdersArg GetUserOrdersArg
	userOrdersArg.UserID = "nikhil_0"

	body, err = json.Marshal(userOrdersArg)
	if err != nil {
		t.Fatalf("Unable to marshal user order arg")
	}
	resp, err = http.Post(baseURL+"/orders/get_user_orders", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Unable to get user orders. Error: %v", err)
	}
	fullBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unable to read add order ret")
	}
	resp.Body.Close()
	err = json.Unmarshal(fullBody, &ordersRet)
	if err != nil {
		t.Fatalf("Unable to deserialize user orders ret. Error: %v", err)
	}
	if ordersRet.Status != http.StatusOK {
		t.Fatalf("Unable to get user orders due to error: %v", ret.ErrorMsg)
	}
	if len(ordersRet.Data.Orders) != 3 {
		t.Fatalf("Expected 3 records for user nikhil_0. Got: %d", len(ordersRet.Data.Orders))
	}

	// Get payment pending orders.
	var ordersArg ScanOrdersArg
	var orderIDs []string
	ordersArg.StartID = 0
	ordersArg.NumOrders = 2
	for ii := 0; ii < 5; ii++ {
		body, err = json.Marshal(ordersArg)
		if err != nil {
			t.Fatalf("Unable to marshal user order arg")
		}
		resp, err = http.Post(baseURL+"/orders/get_payment_pending_orders", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Unable to get user orders. Error: %v", err)
		}
		fullBody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Unable to read add order ret")
		}
		resp.Body.Close()
		err = json.Unmarshal(fullBody, &ordersRet)
		if err != nil {
			t.Fatalf("Unable to deserialize user orders ret. Error: %v", err)
		}
		if ordersRet.Status != http.StatusOK {
			t.Fatalf("Unable to get user orders due to error: %v", ret.ErrorMsg)
		}
		// This assumes that numOrders defined earlier is 5.
		if ii == 0 || ii == 1 {
			if len(ordersRet.Data.Orders) != int(ordersArg.NumOrders) {
				t.Fatalf("Expected %d records. Got: %d",
					ordersArg.NumOrders, len(ordersRet.Data.Orders))
			}
			for _, xx := range ordersRet.Data.Orders {
				orderIDs = append(orderIDs, xx.OrderID)
			}
		} else if ii == 2 {
			if len(ordersRet.Data.Orders) != 1 {
				t.Fatalf("Expected 1 records. Got: %d", len(ordersRet.Data.Orders))
			}
			orderIDs = append(orderIDs, ordersRet.Data.Orders[0].OrderID)
		} else {
			if ordersRet.Data.NextID != -1 {
				t.Fatalf("Scan should have finished but it hasn't")
			}
			break
		}
		// Update the start ID for the next scan.
		ordersArg.StartID = uint64(ordersRet.Data.NextID)
	}

	// Update order status.
	for _, orderID := range orderIDs {
		var updateOrderArg UpdateOrderArg
		var updateOrderRet UpdateOrderRet
		updateOrderArg.OrderID = orderID
		updateOrderArg.Status = "KOrderDeliveryPending"
		body, err = json.Marshal(updateOrderArg)
		if err != nil {
			t.Fatalf("Unable to marshal user order arg")
		}
		resp, err = http.Post(baseURL+"/orders/update_order", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Unable to get user orders. Error: %v", err)
		}
		fullBody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Unable to read add order ret")
		}
		resp.Body.Close()
		err = json.Unmarshal(fullBody, &updateOrderRet)
		if err != nil {
			t.Fatalf("Unable to deserialize update order ret. Error: %v", err)
		}
		if updateOrderRet.Status != http.StatusOK {
			t.Fatalf("Unable to update order due to error: %v", ret.ErrorMsg)
		}
		if !updateOrderRet.Data.RegistrationStatus {
			t.Fatalf("Unable to update order status for order ID: %s", orderID)
		}
	}

	ordersArg.StartID = 0
	ordersArg.NumOrders = 2
	for ii := 0; ii < 5; ii++ {
		body, err = json.Marshal(ordersArg)
		if err != nil {
			t.Fatalf("Unable to marshal scan order arg")
		}
		resp, err = http.Post(baseURL+"/orders/get_delivery_pending_orders", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Unable to get user orders. Error: %v", err)
		}
		fullBody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Unable to read add order ret")
		}
		resp.Body.Close()
		err = json.Unmarshal(fullBody, &ordersRet)
		if err != nil {
			t.Fatalf("Unable to deserialize scan orders ret. Error: %v", err)
		}
		if ordersRet.Status != http.StatusOK {
			t.Fatalf("Unable to scan orders due to error: %v", ret.ErrorMsg)
		}
		// This assumes that numOrders defined earlier is 5.
		if ii == 0 || ii == 1 {
			if len(ordersRet.Data.Orders) != int(ordersArg.NumOrders) {
				t.Fatalf("Expected %d records. Got: %d",
					ordersArg.NumOrders, len(ordersRet.Data.Orders))
			}
			for _, xx := range ordersRet.Data.Orders {
				orderIDs = append(orderIDs, xx.OrderID)
			}
		} else if ii == 2 {
			if len(ordersRet.Data.Orders) != 1 {
				t.Fatalf("Expected 1 records. Got: %d", len(ordersRet.Data.Orders))
			}
			orderIDs = append(orderIDs, ordersRet.Data.Orders[0].OrderID)
		} else {
			if ordersRet.Data.NextID != -1 {
				t.Fatalf("Scan should have finished but it hasn't")
			}
			break
		}
		// Update the start ID for the next scan.
		ordersArg.StartID = uint64(ordersRet.Data.NextID)
	}
}
