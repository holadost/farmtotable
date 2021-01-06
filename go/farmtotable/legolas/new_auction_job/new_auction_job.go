package new_auction_job

import (
	"errors"
	"farmtotable/gandalf"
	"fmt"
	"github.com/golang/glog"
	"time"
)

/*
This job performs the following tasks:
	1. Scans all the items and sees which items can be auctioned
	2. Adds all items that can be auctioned to the main auctions db.
*/

type PopulateNewAuctionsJob struct {
	gandalf    *gandalf.Gandalf
	batchSize  uint64
	it         *gandalf.ItemsScanner
	numWorkers int
	workerErrs []error
	doneChan   chan bool
}

func NewPopulateNewAuctionsJob(g *gandalf.Gandalf, scanSize uint64, numWorkers uint) *PopulateNewAuctionsJob {
	pnaj := PopulateNewAuctionsJob{}
	if g == nil {
		panic("Invalid gandalf object")
	}
	pnaj.gandalf = g
	pnaj.it = gandalf.NewItemsScanner(g, scanSize)
	pnaj.numWorkers = int(numWorkers)
	pnaj.workerErrs = make([]error, pnaj.numWorkers, pnaj.numWorkers)
	for ii := 0; ii < int(pnaj.numWorkers); ii++ {
		pnaj.workerErrs[ii] = nil
	}
	pnaj.doneChan = make(chan bool, pnaj.numWorkers)
	return &pnaj
}

/* Runs the job. This method satisfies the LegolasJob interface. */
func (pnaj *PopulateNewAuctionsJob) Run() {
	glog.Infof("Starting all workers to populate new auctions")
	for ii := 0; ii < int(pnaj.numWorkers); ii++ {
		worker := newWorker(pnaj.gandalf, pnaj.it)
		go func(idx int, worker *_Worker) {
			pnaj.workerErrs[idx] = worker.run()
			pnaj.doneChan <- true
		}(ii, worker)
	}
	glog.Infof("Successfully started %d workers. Waiting for workers to finish", pnaj.numWorkers)
	for ii := 0; ii < pnaj.numWorkers; ii++ {
		<-pnaj.doneChan
	}
	for ii := 0; ii < pnaj.numWorkers; ii++ {
		if pnaj.workerErrs[ii] != nil {
			glog.Fatalf("Worker: %d failed with error: %v", ii, pnaj.workerErrs[ii])
		}
	}
	glog.Info("All workers finished successfully. New auctions populated")
}

/*********************** INTERNAL HELPERS ******************************/
type _Worker struct {
	itemScanner *gandalf.ItemsScanner
	gandalf     *gandalf.Gandalf
}

func newWorker(gandalf *gandalf.Gandalf, it *gandalf.ItemsScanner) *_Worker {
	worker := _Worker{}
	worker.itemScanner = it
	worker.gandalf = gandalf
	return &worker
}

func (worker *_Worker) run() error {
	for {
		items, status, err := worker.itemScanner.NextBatch()
		if err != nil {
			return errors.New(fmt.Sprintf("error while scanning items, error: %v", err))
		}
		if len(items) != 0 {
			var chosenItems []gandalf.AuctionModel
			for _, item := range items {
				if item.AuctionStarted || item.AuctionEnded {
					// We have already handled this auction. Move on.
					continue
				}
				now := time.Now()
				start := item.AuctionStartTime
				if start.Before(now) {
					// This auction can now be started.
					var auction gandalf.AuctionModel
					auction.AuctionStartTime = item.AuctionStartTime
					auction.AuctionDurationSecs = item.AuctionDurationSecs
					auction.ItemQty = item.ItemQty
					auction.ItemID = item.ItemID
					auction.ItemName = item.ItemName
					auction.MaxBid = item.MinPrice
					auction.MinBid = item.MinPrice
					auction.ImageURL = item.ImageURL
					chosenItems = append(chosenItems, auction)
				}
			}
			err := worker.gandalf.AddAuctions(chosenItems)
			if err != nil {
				glog.Errorf("Unable to add the following items to the auctions table due to err: %v\n%v",
					err, chosenItems)
				return err
			}
			glog.V(1).Infof("Successfully added auctions: %v", chosenItems)
			// TODO: Ugly! Bulk update all of these in the future.
			for _, auction := range chosenItems {
				err = worker.gandalf.UpdateItemAuctionStatus(auction.ItemID, true, false, false)
				if err != nil {
					glog.Errorf("Unable to update item: %s auction status due to err: %v", auction.ItemID, err)
					return err
				}
			}
		}
		if status {
			// Scan has completed. Return no error
			return nil
		}
	}
}
