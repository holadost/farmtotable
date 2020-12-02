package legolas

import (
	"farmtotable/gandalf"
	"farmtotable/legolas/auction_winner_job"
	"farmtotable/legolas/new_auction_job"
	"flag"
	"github.com/golang/glog"
	"time"
)

var (
	newAuctionJobScanIntervalSecs = flag.Uint("new_auction_job_scan_interval_secs", 7200,
		"The scan scheduling interval for new auctions(in seconds)")
	winnerJobScanIntervalSecs = flag.Uint("winner_job_scan_interval_secs", 7200,
		"The scan scheduling interval for picking winners of expired auctions(in seconds)")
)

type Legolas struct {
	gandalf *gandalf.Gandalf
	jobQ    chan LegolasJob
	done    chan bool
}

func NewLegolas() *Legolas {
	var lg Legolas
	lg.jobQ = make(chan LegolasJob)
	lg.done = make(chan bool)
	lg.gandalf = gandalf.NewSqliteGandalf()
	return &lg
}

func NewLegolasWithGandalf(g *gandalf.Gandalf) *Legolas {
	var lg Legolas
	lg.jobQ = make(chan LegolasJob)
	lg.done = make(chan bool)
	lg.gandalf = g
	return &lg
}

func (lg *Legolas) Run() {
	glog.Infof("Legolas initialized")
	lg.scheduleJobs()
	lg.executeJobs()
}

func (lg *Legolas) Shutdown() {
	lg.done <- true
}

func (lg *Legolas) executeJobs() {
	for {
		select {
		case <-lg.done:
			glog.Infof("Legolas received exit notification. Exiting")
			return
		case job := <-lg.jobQ:
			job.Run()
		}
	}
}

func (lg *Legolas) scheduleJobs() {
	lg.scheduleNewAuctionsJob()
	lg.scheduleAuctionWinnerJob()
}

func (lg *Legolas) scheduleNewAuctionsJob() {
	ticker := time.NewTicker(time.Duration(*newAuctionJobScanIntervalSecs) * time.Second)
	go func() {
		for {
			select {
			case <-lg.done:
				return
			case <-ticker.C:
				glog.Infof("Scheduling new auctions job")
				job := new_auction_job.NewPopulateNewAuctionsJob(lg.gandalf, 128, 4)
				lg.jobQ <- job
			}
		}
	}()
}

func (lg *Legolas) scheduleAuctionWinnerJob() {
	ticker := time.NewTicker(time.Duration(*winnerJobScanIntervalSecs) * time.Second)
	go func() {
		for {
			select {
			case <-lg.done:
				return
			case <-ticker.C:
				glog.Infof("Scheduling auction winner job")
				job := auction_winner_job.NewAuctionWinnerJob(lg.gandalf, 4, 128)
				lg.jobQ <- job
			}
		}
	}()
}
