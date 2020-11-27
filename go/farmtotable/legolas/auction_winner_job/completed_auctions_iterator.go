package auction_winner_job

import (
	"farmtotable/gandalf"
	"time"
)

type CompletedAuctionsIterator struct {
	nextID       uint64
	currBatch    []gandalf.Auction
	scanSize     uint64
	gandalf      *gandalf.Gandalf
	scanComplete bool
}

func NewCompletedAuctionsIterator(gandalf *gandalf.Gandalf, scanSize uint64) *CompletedAuctionsIterator {
	it := CompletedAuctionsIterator{}
	it.scanSize = scanSize
	it.gandalf = gandalf
	return &it
}

func (it *CompletedAuctionsIterator) Next() (gandalf.Auction, bool /* Scan complete */, error /* scan errors */) {
	if it.scanComplete {
		return gandalf.Auction{}, true, nil
	}
	var item gandalf.Auction
	if len(it.currBatch) > 0 {
		item, it.currBatch = it.currBatch[0], it.currBatch[1:]
		return item, false, nil
	}
	auctions, err := it.gandalf.GetAllAuctions(it.nextID, it.scanSize)
	if err != nil {
		return gandalf.Auction{}, false, err
	}
	it.nextID += it.nextID + it.scanSize
	now := time.Now()
	if len(auctions) == 0 {
		it.scanComplete = true
	}
	for _, auction := range auctions {
		duration := time.Second * time.Duration(auction.AuctionDurationSecs)
		if auction.AuctionStartTime.Add(duration).Before(now) {
			// Auction has expired.
			continue
		}
		it.currBatch = append(it.currBatch, auction)
	}
	item, it.currBatch = it.currBatch[0], it.currBatch[1:]
	return item, false, nil
}
