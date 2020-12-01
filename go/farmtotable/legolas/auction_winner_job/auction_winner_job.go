package auction_winner_job

import (
	"errors"
	"farmtotable/gandalf"
	"fmt"
	"github.com/golang/glog"
	"sort"
)

/*
This job performs the following tasks:
	1. Finds which auctions have expired.
	2. Determines winners for the completed auctions.
	3. Places new orders for the winners.
	4. Notifies the winners via email.
*/

type AuctionWinnerJob struct {
	gandalf    *gandalf.Gandalf
	numWorkers uint
	workerPool []*_Worker
}

func NewAuctionWinnerJob(gandalf *gandalf.Gandalf, numWorkers uint) *AuctionWinnerJob {
	awj := AuctionWinnerJob{}
	awj.gandalf = gandalf
	awj.numWorkers = numWorkers
	return &awj
}

/* Starts the job. This method satisfies the LegolasJob interface. */
func (awj *AuctionWinnerJob) Start() {

}

/* Stops the job. This method satisfies the LegolasJob interface. */
func (awj *AuctionWinnerJob) Stop() {

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
		auction, status, err := worker.auctionScanner.Next()
		if err != nil {
			return errors.New(fmt.Sprintf("error while picking winners for auction: %v, error: %v", auction, err))
		}
		if status {
			// Iterator has finished. No more auctions to process.
			return nil
		}
		topBids, err := worker.pickWinners(auction)
		if err != nil {
			return err
		}
		err = worker.placeOrders(auction, topBids)

	}
}

func (worker *_Worker) pickWinners(auction gandalf.AuctionModel) ([]gandalf.BidModel, error) {
	glog.Infof("Processing winners for auction: %v", auction)
	bidScanner := gandalf.NewItemsBidScanner(worker.gandalf, auction.ItemID, 1024)
	topBids := make([]gandalf.BidModel, 0, 2*KNumWinnersPerItem)
	for {
		batch, finished, err := bidScanner.NextN(KNumWinnersPerItem)
		if err != nil {
			glog.Errorf("Unable to process winners for auction: %v due to err: %v", auction, err)
			return topBids, err
		}
		if (!finished) || (finished && (len(batch) > 0)) {
			topBids = append(topBids, batch...)
			sort.Sort(sortByBidAmount(topBids))
			topBids = topBids[0:KNumWinnersPerItem]
			continue
		}
		if finished {
			topBids = topBids[0:KNumWinnersPerItem]
			break
		}
	}
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
	totalQty := item.ItemQty
	for ii := len(topBids) - 1; ii >= 0; ii-- {
		var order gandalf.OrderModel
		order.ItemID = item.ItemID
		order.UserID = topBids[ii].UserID
		order.ItemPrice = topBids[ii].BidAmount
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
	glog.V(1).Infof("Adding orders to backend: %v", orders)
	err = worker.gandalf.AddOrders(orders)
	if err != nil {
		glog.Errorf("Unable to add orders to backend")
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
func (a sortByBidAmount) Less(i, j int) bool { return a[i].BidAmount < a[j].BidAmount }
func (a sortByBidAmount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
