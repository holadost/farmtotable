package gandalf

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"
	"strings"
	"time"
)

type AuctionStore struct {
	rdb *redis.Client
}

var ctx = context.Background()

func NewAuctionsStore() *AuctionStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	var as AuctionStore
	as.rdb = rdb
	return &as
}

func (as *AuctionStore) AddAuctions(auctions []AuctionModel) error {
	glog.Infof("Adding %d auctions to the auctions store", len(auctions))
	for _, auction := range auctions {
		keyStr := GenerateAuctionKey(auction.ItemID)
		val, err := json.Marshal(auction)
		if err != nil {
			glog.Errorf("Unable to marshal auction to JSON object due to err: %v", err)
			return err
		}
		cmd := as.rdb.Set(ctx, keyStr, val, (time.Duration(auction.AuctionDurationSecs) * time.Second))
		if cmd.Err() != nil {
			glog.Errorf("Unable to put auction in auctions store due to err: %v", cmd.Err())
			return cmd.Err()
		}
	}
	return nil
}

func (as *AuctionStore) RemoveAuctions(itemIDs []string) error {
	glog.Infof("Removing items: %v from auctions store", itemIDs)
	var keys []string
	for _, itemID := range itemIDs {
		key := GenerateAuctionKey(itemID)
		keys = append(keys, key)
	}
	cmd := as.rdb.Del(ctx, keys...)
	if cmd.Err() != nil {
		glog.Errorf("Unable to remove auctions due to err: %v", cmd.Err())
		return cmd.Err()
	}
	return nil
}

func (as *AuctionStore) ScanAuctions(cursor uint64, numAuctions uint64) ([]AuctionModel, uint64, error) {
	keys, newCursor, err := as.rdb.Scan(ctx, cursor, KAuctionKeyPrefix, int64(numAuctions)).Result()
	if err != nil {
		glog.Errorf("Unable to scan keys from redis due to err: %v", err)
		return nil, 0, err
	}
	var auctions []AuctionModel
	for _, key := range keys {
		val, err := as.rdb.Get(ctx, key).Result()
		if err != nil {
			glog.Errorf("Unable to get auction value for key: %s due to err: %v", key, err)
			return auctions, 0, err
		}
		var auction AuctionModel
		if val == "" {
			continue
		}

		err = json.Unmarshal([]byte(val), &auction)
		if err != nil {
			glog.Errorf("Failed to deserialize auction entry for key: %s due to err: %v", key, err)
			return auctions, 0, err
		}
		auctions = append(auctions, auction)
	}
	return auctions, newCursor, nil
}

func (as *AuctionStore) GetAuction(itemID string) (AuctionModel, error) {
	val, err := as.rdb.Get(ctx, GenerateAuctionKey(itemID)).Result()
	var auction AuctionModel
	if err != nil {
		glog.Errorf("Unable to scan keys from redis due to err: %v", err)
		return auction, err
	}
	if val == "" {
		return auction, nil
	}
	err = json.Unmarshal([]byte(val), &auction)
	if err != nil {
		glog.Errorf("Failed to deserialize auction entry for item: %s due to err: %v", itemID, err)
		return auction, err
	}
	return auction, nil
}

func (as *AuctionStore) RegisterBid() {

}

func (as *AuctionStore) ScanItemBids(startBidID string, numAuctions uint64) {

}

/* Generates the auction key based on the given item ID. */
func GenerateAuctionKey(itemID string) string {
	return KAuctionKeyPrefix + KKeyDelimiter + itemID
}

func GetItemIDFromAuctionKey(key string) string {
	fields := strings.Split(key, KKeyDelimiter)
	if fields[0] == KAuctionKeyPrefix {
		return fields[1]
	}
	glog.Errorf("Invalid auction key prefix for key: %s. Expected prefix: %s", key, KAuctionKeyPrefix)
	return ""
}

/* Generates the bid key based on the given item ID and user ID. */
func GenerateBidKey(itemID string, userID string) string {
	return KBidKeyPrefix + KKeyDelimiter + itemID + KKeyDelimiter + userID
}

func GetItemAndUserIDFromBidKey(key string) (string, string) {
	fields := strings.Split(key, KKeyDelimiter)
	if fields[0] == KBidKeyPrefix {
		return fields[1], fields[2]
	}
	glog.Errorf("Incorrect bid key prefix for key: %s", key)
	return "", ""
}
