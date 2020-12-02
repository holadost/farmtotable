package gandalf

import (
	"sync"
	"time"
)

const (
	KAllAuctions     = 0
	KLiveAuctions    = 1
	KExpiredAuctions = 2
)

/* Scans all the auctions. The scanner is thread-safe. */
type AuctionsScanner struct {
	nextID       uint64
	currBatch    []AuctionModel
	scanSize     uint64
	gandalf      *Gandalf
	scanComplete bool
	mu           sync.Mutex
	scanErr      error
	auctionType  uint
}

func NewAuctionsScanner(gandalf *Gandalf, scanSize uint64) *AuctionsScanner {
	return newAuctionsScanner(gandalf, scanSize, KAllAuctions)
}

func NewLiveAuctionsScanner(gandalf *Gandalf, scanSize uint64) *AuctionsScanner {
	return newAuctionsScanner(gandalf, scanSize, KLiveAuctions)
}

func NewExpiredAuctionsScanner(gandalf *Gandalf, scanSize uint64) *AuctionsScanner {
	return newAuctionsScanner(gandalf, scanSize, KExpiredAuctions)
}

func newAuctionsScanner(gandalf *Gandalf, scanSize uint64, auctionType uint) *AuctionsScanner {
	it := AuctionsScanner{}
	it.scanSize = scanSize
	it.gandalf = gandalf
	it.scanComplete = false
	it.auctionType = auctionType
	it.scanErr = nil
	return &it
}

/*
Scans the next batch if currBatch is empty and scan is not complete. This method is not thread-safe.
This method should only be called after the caller holds the scanner mutex.
*/
func (it *AuctionsScanner) maybeScanNextBatch() {
	if it.scanComplete {
		return
	}
	if len(it.currBatch) > 0 {
		return
	}
	auctions, err := it.gandalf.GetAllAuctions(it.nextID, it.scanSize)
	if err != nil {
		// There was an error. Mark the scan as complete to avoid using
		// the scanner beyond this.
		it.scanComplete = true
		it.scanErr = err
		return
	}
	if len(auctions) == 0 {
		it.scanComplete = true
		return
	}
	it.nextID += it.nextID + it.scanSize
	now := time.Now()
	for _, auction := range auctions {
		deadline := auction.AuctionStartTime.Add(time.Second * time.Duration(auction.AuctionDurationSecs))
		if it.auctionType == KAllAuctions {
			it.currBatch = append(it.currBatch, auction)
		} else if it.auctionType == KLiveAuctions {
			if deadline.After(now) {
				it.currBatch = append(it.currBatch, auction)
			}
		} else if it.auctionType == KExpiredAuctions {
			if deadline.Before(now) {
				it.currBatch = append(it.currBatch, auction)
			}
		}
	}
}

func (it *AuctionsScanner) Next() (AuctionModel, bool /* Scan complete */, error /* scan errors */) {
	it.mu.Lock()
	defer it.mu.Unlock()
	it.maybeScanNextBatch()
	if it.scanComplete {
		return AuctionModel{}, it.scanComplete, it.scanErr
	}
	var item AuctionModel
	item, it.currBatch = it.currBatch[0], it.currBatch[1:]
	return item, it.scanComplete, it.scanErr
}

func (it *AuctionsScanner) NextBatch() ([]AuctionModel, bool /* Scan complete */, error /* scan errors */) {
	it.mu.Lock()
	defer it.mu.Unlock()
	it.maybeScanNextBatch()
	if it.scanComplete {
		return []AuctionModel{}, it.scanComplete, it.scanErr
	}
	auctions := make([]AuctionModel, 0, len(it.currBatch))
	for _, auction := range it.currBatch {
		auctions = append(auctions, auction)
	}
	return auctions, it.scanComplete, it.scanErr
}
