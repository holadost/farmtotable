package gandalf

import (
	"github.com/golang/glog"
	"sync"
)

/* Scans all the bids for a given item. The scanner is thread-safe. */
type ItemBidsScanner struct {
	nextID       uint64
	currBatch    []Bid
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
	it.currBatch = make([]Bid, 0, it.scanSize)
	return &it
}

/*
Scans the next batch if currBatch is empty and scan is not complete. This method is not thread-safe.
This method should only be called after the caller holds the scanner mutex.
*/
func (it *ItemBidsScanner) maybeScanNextBatch() {
	if it.scanComplete {
		return
	}
	if len(it.currBatch) > 0 {
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

func (it *ItemBidsScanner) Next() (Bid, bool /* Scan complete */, error /* scan errors */) {
	it.mu.Lock()
	defer it.mu.Unlock()
	it.maybeScanNextBatch()
	var item Bid
	if it.scanComplete {
		return item, it.scanComplete, it.scanErr
	}
	if len(it.currBatch) == 0 {
		glog.Fatalf("currBatch is empty even though scanner is not complete")
	}
	item, it.currBatch = it.currBatch[0], it.currBatch[1:]
	return item, it.scanComplete, it.scanErr
}

func (it *ItemBidsScanner) NextBatch() ([]Bid, bool /* Scan complete */, error /* scan errors */) {
	it.mu.Lock()
	defer it.mu.Unlock()
	it.maybeScanNextBatch()
	if it.scanComplete {
		return []Bid{}, it.scanComplete, it.scanErr
	}
	if len(it.currBatch) == 0 {
		glog.Fatalf("currBatch is empty even though scanner is not complete")
	}
	items := make([]Bid, 0, len(it.currBatch))
	for _, bid := range it.currBatch {
		items = append(items, bid)
	}
	// Clear the currBatch but keep the underlying memory.
	it.currBatch = it.currBatch[:0]
	return items, it.scanComplete, it.scanErr
}

func (it *ItemBidsScanner) NextN(n uint) ([]Bid, bool /* Scan complete */, error /* scan errors */) {
	it.mu.Lock()
	defer it.mu.Unlock()
	it.maybeScanNextBatch()
	if it.scanComplete {
		return []Bid{}, it.scanComplete, it.scanErr
	}
	if len(it.currBatch) == 0 {
		glog.Fatalf("currBatch is empty even though scanner is not complete")
	}
	items := make([]Bid, 0, n)
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
