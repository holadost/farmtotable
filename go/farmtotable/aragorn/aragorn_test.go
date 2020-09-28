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
	baseURL := "http://localhost:8080"
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
	resp, err := http.Post(baseURL+"/api/v1/resources/users/register", "application/json", bytes.NewBuffer(body))
	if err != nil || resp.StatusCode >= 300 {
		t.Fatalf("Unable to register user. Error: %v", err)
	}
	fullBody, err := ioutil.ReadAll(resp.Body)
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
}
