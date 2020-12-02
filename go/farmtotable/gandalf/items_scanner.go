package gandalf

import (
	"sync"
)

/* Scans all the bids for a given item. The scanner is thread-safe. */
type ItemsScanner struct {
	nextID       uint64
	currBatch    []ItemModel
	scanSize     uint64
	gandalf      *Gandalf
	scanComplete bool
	mu           sync.Mutex
	scanErr      error
}

func NewItemsScanner(gandalf *Gandalf, scanSize uint64) *ItemsScanner {
	it := ItemsScanner{}
	it.scanSize = scanSize
	it.gandalf = gandalf
	it.nextID = 0
	it.scanComplete = false
	it.scanErr = nil
	it.currBatch = make([]ItemModel, 0, it.scanSize)
	return &it
}

/*
Scans the next batch if currBatch is empty and scan is not complete. This method is not thread-safe.
This method should only be called after the caller holds the scanner mutex.
*/
func (it *ItemsScanner) maybeScanNextBatch() {
	if it.scanComplete || (len(it.currBatch) > 0) {
		return
	}
	items, err := it.gandalf.ScanItems(it.nextID, it.scanSize)
	if err != nil {
		// Error while scanning bids. Mark the scan as complete so that
		// future readers don't get wrong values.
		it.scanComplete = true
		it.scanErr = err
		return
	}
	if len(items) == 0 {
		// There are no more items. Mark the scanner as complete.
		it.scanComplete = true
		return
	}
	for _, item := range items {
		it.currBatch = append(it.currBatch, item)
	}
	it.nextID += it.scanSize
	return
}

func (it *ItemsScanner) Next() (ItemModel, bool /* Scan complete */, error /* scan errors */) {
	it.mu.Lock()
	defer it.mu.Unlock()
	it.maybeScanNextBatch()
	var item ItemModel
	if len(it.currBatch) > 0 {
		item, it.currBatch = it.currBatch[0], it.currBatch[1:]
		return item, it.scanComplete, it.scanErr
	} else {
		return item, it.scanComplete, it.scanErr
	}
}

func (it *ItemsScanner) NextBatch() ([]ItemModel, bool /* Scan complete */, error /* scan errors */) {
	it.mu.Lock()
	defer it.mu.Unlock()
	it.maybeScanNextBatch()
	if it.scanComplete {
		return []ItemModel{}, it.scanComplete, it.scanErr
	}
	items := make([]ItemModel, 0, len(it.currBatch))
	for _, item := range it.currBatch {
		items = append(items, item)
	}
	// Clear the currBatch but keep the underlying memory.
	it.currBatch = it.currBatch[:0]
	return items, it.scanComplete, it.scanErr
}

func (it *ItemsScanner) NextN(n uint) ([]ItemModel, bool /* Scan complete */, error /* scan errors */) {
	it.mu.Lock()
	defer it.mu.Unlock()
	it.maybeScanNextBatch()
	if it.scanComplete {
		return []ItemModel{}, it.scanComplete, it.scanErr
	}
	items := make([]ItemModel, 0, n)
	for ii, item := range it.currBatch {
		if uint(ii) == n {
			break
		}
		items = append(items, item)
	}
	it.currBatch = it.currBatch[n:]
	return items, it.scanComplete, it.scanErr
}
