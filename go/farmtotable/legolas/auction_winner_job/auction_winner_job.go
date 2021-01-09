package auction_winner_job

import (
	"errors"
	"farmtotable/gandalf"
	"fmt"
	"github.com/golang/glog"
	"sort"
	"time"
)

/*
This job performs the following tasks:
	1. Finds which auctions have expired.
	2. Determines winners for the completed auctions.
	3. Places new orders for the winners.
	4. Updates the auction status for the items for whom auction has expired.
	5. Notifies the winners via email.
*/

type AuctionWinnerJob struct {
	gandalf        *gandalf.Gandalf
	numWorkers     int
	workerErrs     []error
	doneChan       chan bool
	auctionScanner *gandalf.AuctionsScanner
	scanSize       uint64
}

func NewAuctionWinnerJob(g *gandalf.Gandalf, numWorkers uint, scanSize uint64) *AuctionWinnerJob {
	awj := AuctionWinnerJob{}
	awj.gandalf = g
	awj.numWorkers = int(numWorkers)
	awj.workerErrs = make([]error, awj.numWorkers, awj.numWorkers)
	for ii := 0; ii < awj.numWorkers; ii++ {
		awj.workerErrs[ii] = nil
	}
	awj.doneChan = make(chan bool, awj.numWorkers)
	awj.scanSize = scanSize
	awj.auctionScanner = gandalf.NewExpiredAuctionsScanner(awj.gandalf, awj.scanSize)
	return &awj
}

/* Starts the job. This method satisfies the LegolasJob interface. */
func (awj *AuctionWinnerJob) Run() {
	glog.Infof("Starting all workers to process auction winners")
	for ii := 0; ii < int(awj.numWorkers); ii++ {
		worker := newWorker(awj.gandalf, awj.auctionScanner)
		go func(idx int, worker *_Worker) {
			awj.workerErrs[idx] = worker.run()
			awj.doneChan <- true
		}(ii, worker)
	}
	glog.Infof("Successfully started %d workers. Waiting for workers to finish", awj.numWorkers)
	for ii := 0; ii < awj.numWorkers; ii++ {
		<-awj.doneChan
	}
	for ii := 0; ii < awj.numWorkers; ii++ {
		if awj.workerErrs[ii] != nil {
			glog.Fatalf("Auction Winners Job Worker: %d failed with error: %v", ii, awj.workerErrs[ii])
		}
	}
	glog.Info("All workers finished successfully. Expired auctions successfully processed")
}

/*********************** INTERNAL HELPERS ******************************/
type _Worker struct {
	auctionScanner *gandalf.AuctionsScanner
	gandalf        *gandalf.Gandalf
}

func newWorker(gandalf *gandalf.Gandalf, it *gandalf.AuctionsScanner) *_Worker {
	worker := _Worker{}
	worker.auctionScanner = it
	worker.gandalf = gandalf
	return &worker
}

func (worker *_Worker) run() error {
	for {
		auction, finished, err := worker.auctionScanner.Next()
		if err != nil {
			return errors.New(fmt.Sprintf("error while picking winners for auction: %v, error: %v", auction, err))
		}
		if auction.ItemID == "" && finished {
			glog.Infof("Completed expired auctions scan")
			return nil
		} else if auction.ItemID == "" && !finished {
			glog.Infof("Did not find any expired auctions in current batch. Continuing scan")
			continue
		}
		topBids, err := worker.pickWinners(auction)
		if err != nil {
			return err
		}
		err = worker.placeOrders(auction, topBids)
		if err != nil {
			return err
		}
		err = worker.updateItem(auction)
		if err != nil {
			return err
		}
		err = worker.notifyWinners(auction, topBids)
		if err != nil {
			return err
		}
		if finished {
			return nil
		}
	}
}

func (worker *_Worker) pickWinners(auction gandalf.AuctionModel) ([]gandalf.BidModel, error) {
	glog.Infof("Processing winners for auction: %v", auction)
	bidScanner := gandalf.NewItemsBidScanner(worker.gandalf, auction.ItemID, 1024)
	var topBids []gandalf.BidModel
	for {
		batch, finished, err := bidScanner.NextN(KNumWinnersPerItem)
		if err != nil {
			glog.Errorf("Unable to process winners for auction: %v due to err: %v", auction, err)
			return topBids, err
		}
		if (!finished) && (len(batch) == 0) {
			continue
		}
		if (!finished) || (finished && (len(batch) > 0)) {
			topBids = append(topBids, batch...)
			sort.Sort(sortByBidAmount(topBids))
			if len(topBids) > KNumWinnersPerItem {
				topBids = topBids[0:KNumWinnersPerItem]
			}
			continue
		}
		if topBids == nil {
			return topBids, nil
		}
		if finished {
			if len(topBids) > KNumWinnersPerItem {
				topBids = topBids[0:KNumWinnersPerItem]
			}
			break
		}
	}
	// Filter top bids as we may have empty values in there.
	return topBids, nil
}

func (worker *_Worker) placeOrders(auction gandalf.AuctionModel, topBids []gandalf.BidModel) error {
	glog.Infof("Placing orders for auction: %v", auction)
	item, err := worker.gandalf.GetItem(auction.ItemID)
	if err != nil {
		return err
	}
	if item.ItemID == "" {
		return errors.New(
			fmt.Sprintf(
				"unable to fetch item with item ID: %s from backend", item.ItemID))
	}
	var orders []gandalf.OrderModel
	currTime := time.Now()
	totalQty := item.ItemQty
	for ii := len(topBids) - 1; ii >= 0; ii-- {
		var order gandalf.OrderModel
		order.ItemID = item.ItemID
		order.UserID = topBids[ii].UserID
		order.ItemPrice = topBids[ii].BidAmount
		order.CreatedDate = currTime
		order.UpdatedDate = currTime
		if totalQty >= topBids[ii].BidQty {
			order.ItemQty = topBids[ii].BidQty
			totalQty -= topBids[ii].BidQty
		} else {
			glog.Warningf(
				"Unable to satisfy order from bid: %v due to "+
					"insufficient item quantity(%d)", topBids[ii], totalQty)
			continue
		}
		orders = append(orders, order)
	}
	// TODO: We should ensure that the orders being placed now haven't been placed
	// TODO: already.
	glog.V(1).Infof("Adding %d orders to backend", len(orders))
	err = worker.gandalf.AddOrders(orders)
	if err != nil {
		glog.Errorf("Unable to add orders to backend due to err: %v", err)
		return err
	}
	return nil
}

func (worker *_Worker) updateItem(auction gandalf.AuctionModel) error {
	glog.Infof("Updating auction status for item: %s", auction.ItemID)
	err := worker.gandalf.UpdateItemAuctionStatus(auction.ItemID, true, true, true)
	if err != nil {
		glog.Errorf("Unable to update item auction status for item ID: %s due to err: %v", auction.ItemID, err)
		return err
	}
	return nil
}

func (worker *_Worker) notifyWinners(auction gandalf.AuctionModel, topBids []gandalf.BidModel) error {
	// TODO: Notify the winners via email. We need to figure this out.
	return nil
}

/* Sorting interface to sort the bids by BidModel amount. */
type sortByBidAmount []gandalf.BidModel

func (a sortByBidAmount) Len() int           { return len(a) }
func (a sortByBidAmount) Less(i, j int) bool { return a[i].BidAmount > a[j].BidAmount } // Sort in descending order.
func (a sortByBidAmount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
