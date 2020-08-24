package gandalf

import (
	"farmtotable/gandalf/backend"
	"time"
)

type Gandalf struct {
	backend *backend.Backend
}

func (gandalf *Gandalf) RegisterUser(userID string, name string, emailID string, phNum string, address string) {

}

func (gandalf *Gandalf) GetUserByID(userID string) {

}

func (gandalf *Gandalf) GetUserByEmailID(emailID string) {

}

func (gandalf *Gandalf) GetUserByPhNo(phNum string) {

}

func (gandalf *Gandalf) RegisterItem(itemName string, itemDesc string, itemQty uint32, auctionStartTime time.Time) {

}

func (gandalf *Gandalf) GetItem(itemID string) {

}

func (gandalf *Gandalf) EditItem(itemID string, itemName string, itemDesc string, itemQty uint32, auctionStartTime time.Time) {

}

func (gandalf *Gandalf) AddBid(itemID string, bidAmount float32, bidQty uint32) {

}

func (gandalf *Gandalf) GetMaxBid(itemID string) {

}

func (gandalf *Gandalf) GetAllAuctions() {

}

func (gandalf *Gandalf) GetUserOrders(userID string) {

}

func (gandalf *Gandalf) GetUserPendingOrders(userID string) {

}

func (gandalf *Gandalf) GetUserCompletedOrders(userID string) {

}

func (gandalf *Gandalf) GetOrder(orderID string) {

}

func (gandalf *Gandalf) UpdateOrder(orderID string, status uint32) {

}
