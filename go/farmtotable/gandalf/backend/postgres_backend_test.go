package backend

import (
	"fmt"
	"testing"
)

func TestNewPostgresBackend(t *testing.T) {
	fmt.Println("Testing new postgres backend using gorm")
	NewPostgresBackend()
	fmt.Println("Successfully initialized and created schema")
}

func TestPostgresBackend_Users(t *testing.T) {
	fmt.Println("Testing users")
	backend := NewPostgresBackend()
	backend.AddUser("ssnikhil87", "Nikhil Srivatsan Srinivasan", "nik.Gunner4life@gmail.com",
		"9198029981", "1370, Mills St, Menlo Park, Apt E, California 94025")
	user := backend.GetUserByID("ssnikhil87")
	if user.EmailID != "nik.Gunner4life@gmail.com" {
		t.Fatalf("Didn't find the correct user")
	}
}
