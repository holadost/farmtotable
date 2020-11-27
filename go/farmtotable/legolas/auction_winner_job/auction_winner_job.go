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
	it      *gandalf.AuctionsScanner
	gandalf *gandalf.Gandalf
}

func newWorker(gandalf *gandalf.Gandalf, it *gandalf.AuctionsScanner) *_Worker {
	worker := _Worker{}
	worker.it = it
	worker.gandalf = gandalf
	return &worker
}

func (worker *_Worker) run() error {
	for {
		auction, status, err := worker.it.Next()
		if err != nil {
			return errors.New(fmt.Sprintf("error while picking winners for auction: %v, error: %v", auction, err))
		}
		if status {
			// Iterator has finished. No more auctions to process.
			return nil
		}
		worker.pickWinners(auction)
	}
}

func (worker *_Worker) pickWinners(auction gandalf.Auction) {
	glog.Infof("Processing winners for auction: %v", auction)
	bidScanner := gandalf.NewItemsBidScanner(worker.gandalf, auction.ItemID, 1024)
	topBids := make([]gandalf.Bid, 0, 2*KNumWinnersPerItem)
	for {
		batch, finished, err := bidScanner.NextN(KNumWinnersPerItem)
		if err != nil {
			glog.Errorf("Unable to process winners for auction: %v due to err: %v", auction, err)
			return
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
	// We now have the winners of the auction in topBids.
}

/* Sorting interface to sort the bids by Bid amount. */
type sortByBidAmount []gandalf.Bid

func (a sortByBidAmount) Len() int           { return len(a) }
func (a sortByBidAmount) Less(i, j int) bool { return a[i].BidAmount < a[j].BidAmount }
func (a sortByBidAmount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
