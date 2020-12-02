package gandalf

import (
	"sync"
)

/* Scans all the bids for a given item. The scanner is thread-safe. */
type ItemBidsScanner struct {
	nextID       uint64
	currBatch    []BidModel
	scanSize     uint64
	gandalf      *Gandalf
	scanComplete bool
	itemID       string
	mu           sync.Mutex
	scanErr      error
}

func NewItemsBidScanner(gandalf *Gandalf, itemID string, scanSize uint64) *ItemBidsScanner {
	it := ItemBidsScanner{}
	it.scanSize = scanSize
	it.gandalf = gandalf
	it.nextID = 0
	it.scanComplete = false
	it.itemID = itemID
	it.scanErr = nil
	it.currBatch = make([]BidModel, 0, it.scanSize)
	return &it
}

/*
Scans the next batch if currBatch is empty and scan is not complete. This method is not thread-safe.
This method should only be called after the caller holds the scanner mutex.
*/
func (it *ItemBidsScanner) maybeScanNextBatch() {
	if it.scanComplete || (len(it.currBatch) > 0) {
		return
	}
	bids, err := it.gandalf.ScanItemBids(it.itemID, it.nextID, it.scanSize)
	if err != nil {
		// Error while scanning bids. Mark the scan as complete so that
		// future readers don't get wrong values.
		it.scanComplete = true
		it.scanErr = err
		return
	}
	if len(bids) == 0 {
		// There are no more bids. Mark the scanner as complete.
		it.scanComplete = true
		return
	}
	it.nextID += it.nextID + it.scanSize
	return
}

func (it *ItemBidsScanner) Next() (BidModel, bool /* Scan complete */, error /* scan errors */) {
	it.mu.Lock()
	defer it.mu.Unlock()
	it.maybeScanNextBatch()
	var item BidModel
	if it.scanComplete {
		return item, it.scanComplete, it.scanErr
	}
	item, it.currBatch = it.currBatch[0], it.currBatch[1:]
	return item, it.scanComplete, it.scanErr
}

func (it *ItemBidsScanner) NextBatch() ([]BidModel, bool /* Scan complete */, error /* scan errors */) {
	it.mu.Lock()
	defer it.mu.Unlock()
	it.maybeScanNextBatch()
	if it.scanComplete {
		return []BidModel{}, it.scanComplete, it.scanErr
	}
	items := make([]BidModel, 0, len(it.currBatch))
	for _, bid := range it.currBatch {
		items = append(items, bid)
	}
	// Clear the currBatch but keep the underlying memory.
	it.currBatch = it.currBatch[:0]
	return items, it.scanComplete, it.scanErr
}

func (it *ItemBidsScanner) NextN(n uint) ([]BidModel, bool /* Scan complete */, error /* scan errors */) {
	it.mu.Lock()
	defer it.mu.Unlock()
	it.maybeScanNextBatch()
	if it.scanComplete {
		return []BidModel{}, it.scanComplete, it.scanErr
	}
	items := make([]BidModel, 0, n)
	for ii, bid := range it.currBatch {
		if uint(ii) == n {
			break
		}
		items = append(items, bid)
	}
	// Clear the currBatch but keep the underlying memory.
	it.currBatch = it.currBatch[n:]
	return items, it.scanComplete, it.scanErr
}
