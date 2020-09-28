package aragorn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

func startAragorn() *Aragorn {
	aragorn := NewAragorn()
	fmt.Println("Starting aragorn service in a background go routine")
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
	startAragorn()
	time.Sleep(100 * time.Millisecond)
	baseURL := "http://localhost:8080/api/v1/resources"

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
	resp, err = http.Post(baseURL+"/items/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Unable to register item. Error: %v", err)
	}
	fullBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unable to read body")
	}
	resp.Body.Close()
	getSuppItemsRet := GetSupplierItemsRet{}
	err = json.Unmarshal(fullBody, &getSuppItemsRet)
	if err != nil {
		t.Fatalf("Unable to deserialize reg items ret. Error: %v", err)
	}
	if getSuppItemsRet.Status != http.StatusOK {
		t.Fatalf("Unable to register item. Error Code: %d, Error Message: %s", getSuppItemsRet.Status,
			getSuppItemsRet.ErrorMsg)
	}
	if len(getSuppItemsRet.Data) != 1 {
		t.Fatalf("Failure while fetching all supplier items")
	}
}
