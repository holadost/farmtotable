package auction_winner_job

/*
This job performs the following tasks:
	1. Finds which auctions have expired and removes them from the main auctions table.
	2. Determines winners for the completed auctions.
	3. Places new orders for the winners.
	4. Notifies the winners via email.
*/

type AuctionWinnerJob struct {
}

func NewAuctionWinnerJob() (awj *AuctionWinnerJob) {
	return
}

/* Starts the job. This method satisfies the LegolasJob interface. */
func (awj *AuctionWinnerJob) Start() {

}

/* Stops the job. This method satisfies the LegolasJob interface. */
func (awj *AuctionWinnerJob) Stop() {

}
