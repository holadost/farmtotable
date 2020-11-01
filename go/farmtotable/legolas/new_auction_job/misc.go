package new_auction_job

import (
	"fmt"
	"os"
	"strconv"
)

/* Returns the path for this job. */
func getJobPath() string {
	var tmpPath string
	value, exists := os.LookupEnv("FTT_NEW_AUCTIONS_JOB_TMP_PATH")
	if !exists {
		tmpPath = "/tmp/ftt/legolas/new_auctions_job"
	} else {
		tmpPath = value
	}
	err := os.MkdirAll(tmpPath, 0755)
	if err != nil {
		panic("Unable to create directory for legolas job")
	}
	return value
}

/* Returns the path for this job. */
func getScanItemsBatchSize() uint64 {
	var batchSize uint64
	value, exists := os.LookupEnv("FTT_NEW_AUCTIONS_JOB_SCAN_ITEMS_BATCH_SIZE")
	if !exists {
		batchSize = 100
	} else {
		bs, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("Invalid scan items batch size: %v", value))
		}
		if bs <= 0 {
			panic(fmt.Sprintf("scan items batch size: %v must be > 0", bs))
		}
		batchSize = uint64(bs)
	}
	return batchSize
}
