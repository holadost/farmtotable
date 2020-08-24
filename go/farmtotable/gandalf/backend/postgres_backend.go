package backend

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type PostgresBackend struct {
	db *gorm.DB
}

/* Database model/schema */
type User struct {
	UserID  string `gorm:"type:varchar(100);PRIMARY_KEY"`
	Name    string `gorm:"type:varchar(255);NOT NULL"`
	EmailID string `gorm:"type:varchar(255);NOT NULL"`
	PhNum   string `gorm:"type:varchar(20);NOT NULL"`
	Address string `gorm:"NOT NULL"`
}

type Item struct {
	ItemID           string `gorm:"type:varchar(32);PRIMARY_KEY"`
	ItemName         string `gorm:"type:varchar(255);NOT NULL"`
	ItemDescription  string `gorm:"NOT NULL"`
	ItemQty          uint32
	AuctionStartTime time.Time
	MinPrice         float32
}

type Bid struct {
	ItemID    string `gorm:"type:varchar(32)"`
	UserID    string `gorm:"type:varchar(100)"`
	BidAmount float32
	BidQty    uint32
	CreatedAt time.Time
}

type Auction struct {
	ItemID              string `gorm:"type:varchar(32)"`
	ItemQty             uint32
	AuctionStartTime    time.Time
	AuctionDurationSecs uint64
	MaxBid              float32
}

type Order struct {
	UserID    string
	ItemID    string
	ItemQty   uint32
	ItemPrice float32
	Status    uint32
}

func NewPostgresBackend() *PostgresBackend {
	backend := PostgresBackend{}
	backend.Initialize()
	return &backend
}

func (pgBackend *PostgresBackend) Initialize() {
	args := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := gorm.Open("postgres", args)
	pgBackend.db = db
	if err != nil {
		panic("Unable to initialize backend database")
	}
	pgBackend.db.AutoMigrate(&User{}, &Item{}, &Auction{}, &Bid{}, &Order{})
}

func (pgBackend *PostgresBackend) AddUser(userID string, userName string, emailID string, phNum string, address string) error {
	user := &User{
		UserID:  userID,
		Name:    userName,
		EmailID: emailID,
		PhNum:   phNum,
		Address: address,
	}
	dbc := pgBackend.db.Create(user)
	if dbc.Error != nil {
		return dbc.Error
	}
	return nil
}

func (pgBackend *PostgresBackend) GetUserByID(userID string) (user User) {
	pgBackend.db.Where("user_id = ?", userID).First(&user)
	return
}

func (pgBackend *PostgresBackend) GetUserByEmailID(emailID string) (user User) {
	pgBackend.db.Where("email_id = ?", emailID).First(&user)
	return
}

func (pgBackend *PostgresBackend) GetUserByPhNum(phNum string) (user User) {
	pgBackend.db.Where("ph_num = ?", phNum).First(&user)
	return
}

func (pgBackend *PostgresBackend) AddItem() {

}

func (pgBackend *PostgresBackend) UpdateItem() {

}

func (pgBackend *PostgresBackend) GetItem() {

}
