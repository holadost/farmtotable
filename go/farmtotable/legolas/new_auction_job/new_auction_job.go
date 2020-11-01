package new_auction_job

import (
	"farmtotable/gandalf"
	"time"
)

/*
This job performs the following tasks:
	1. Scans all the items and sees which items can be auctioned
	2. Adds all items that can be auctioned to the main auctions db.
*/

type PopulateNewAuctionsJob struct {
	gandalf *gandalf.Gandalf
	currIndex uint64
	batchSize uint64
	jobPath string
}

func NewPopulateNewAuctionsJob(g *gandalf.Gandalf) *PopulateNewAuctionsJob {
	pnaj := PopulateNewAuctionsJob{}
	if g == nil {
		panic("Invalid gandalf object")
	}
	pnaj.gandalf = g
	pnaj.currIndex = 0
	pnaj.batchSize = getScanItemsBatchSize()
	pnaj.jobPath = getJobPath()
	return &pnaj
}

/* Starts the job. This method satisfies the LegolasJob interface. */
func (pnaj *PopulateNewAuctionsJob) Start() {

}

/* Stops the job. This method satisfies the LegolasJob interface. */
func (pnaj *PopulateNewAuctionsJob) Stop() {

}

func (pnaj *PopulateNewAuctionsJob) populateAuctionableItems() {
	for {
		items := pnaj.fetchNextBatch()
		if len(items) == 0 {
			break
		}
		for _, item := range items {
			currTime := time.Now()
			if currTime >
		}
	}
}

func (pnaj *PopulateNewAuctionsJob) fetchNextBatch() []gandalf.Item {
	var items []gandalf.Item
	return items
}