package gandalf

import (
	"log"
	"os"
	"testing"
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
