package aragorn

import (
	"errors"
	"farmtotable/gandalf"
	"fmt"
)

type Aragorn struct {
	gandalf *gandalf.Gandalf
}

func NewAragorn() *Aragorn {
	aragorn := &Aragorn{}
	// TODO: Pick the backend type based on env. For now hardcode to sqlite.
	aragorn.gandalf = gandalf.NewSqliteGandalf()
	return aragorn
}

func (aragorn *Aragorn) Authenticate(token string) error {
	return nil
}

func (aragorn *Aragorn) GetUser(userID string) (gandalf.User, error) {
	user := aragorn.gandalf.GetUserByID(userID)
	return user, nil
}

func (aragorn *Aragorn) RegisterUser(user gandalf.User) error {
	err := aragorn.gandalf.RegisterUser(user.UserID, user.Name, user.EmailID, user.PhNum, user.Address)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to register user. Backend Error: %v", err))
	}
	return nil
}

func (aragorn *Aragorn) RegisterItem(item gandalf.Item) error {
	err := aragorn.gandalf.RegisterItem(item.UserID, item.ItemName, item.ItemDescription, item.ItemQty, item.AuctionStartTime, item.MinPrice)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to register item. Backend Error: %v", err))
	}
	return nil
}

func (aragorn *Aragorn) EditItem(item gandalf.Item) error {
	oldItem := aragorn.gandalf.GetItem(item.ItemID)
	if oldItem.ItemID == "" {
		return errors.New(fmt.Sprintf("Did not find item with ID: %s", item.ItemID))
	}
	if item.MinPrice == 0.0 {
		item.MinPrice = oldItem.MinPrice
	}
	if item.ItemQty == 0 {
		item.ItemQty = oldItem.ItemQty
	}
	if item.ItemDescription == "" {
		item.ItemDescription = oldItem.ItemDescription
	}
	if item.ItemName == "" {
		item.ItemName = oldItem.ItemName
	}
	item.UserID = oldItem.UserID
	err := aragorn.gandalf.EditItem(item.ItemID, item.ItemName, item.ItemDescription, item.ItemQty, item.AuctionStartTime, item.MinPrice)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to edit item: %s. Backend Error: %v", item.ItemID, err))
	}
	return nil
}

func (aragorn *Aragorn) RemoveItem(itemID string) error {
	err := aragorn.gandalf.DeleteItem(itemID)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to remove item: %s. Backend error: %v", itemID, err))
	}
	return nil
}

func (aragorn *Aragorn) GetUserItems(userID string) ([]gandalf.Item, error) {
	items, err := aragorn.gandalf.GetUserItems(userID)
	if err != nil {
		return items, errors.New(fmt.Sprintf("Unable to fetch user items for user: %s. Backend Error: %v", userID, err))
	}
	return items, nil
}

func (aragorn *Aragorn) GetAllAuctions(startIdx uint64, numAuctions uint64) ([]gandalf.Auction, error) {
	auctions, err := aragorn.gandalf.GetAllAuctions(startIdx, numAuctions)
	if err != nil {
		return auctions, errors.New(fmt.Sprintf("Unable to fetch auctions. Backend Error: %v", err))
	}
	return auctions, nil
}

func (aragorn *Aragorn) GetMaxBids(itemIDs []string) (map[string]float32, error) {
	auctions, err := aragorn.gandalf.GetMaxBids(itemIDs)
	maxBidMap := make(map[string]float32)
	if err != nil {
		return maxBidMap, errors.New(fmt.Sprintf("Unable to fetch max bids for items. Backend error: %v", err))
	}
	for ii := 0; ii < len(auctions); ii++ {
		maxBidMap[auctions[ii].ItemID] = auctions[ii].MaxBid
	}
	return maxBidMap, nil
}

func (aragorn *Aragorn) RegisterBid(itemID string, userID string, bidAmount float32, bidQty uint32) error {
	err := aragorn.gandalf.RegisterBid(itemID, userID, bidAmount, bidQty)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to register bid for item: %s by user: %s. Backend Error: %v", itemID, userID, err))
	}
	return nil
}
