package gandalf

import (
	"fmt"
	"github.com/golang/glog"
)

// Various gandalf related errors.
const (
	KLowerBidAmount = iota
	KInvalidBidAmount
	KInvalidBidQuantity
	KTimeout
	KInvalidItem
	KInvalidUser
	KInvalidSupplier
	KAuctionExpired
	KGandalfBackendError
)

type GandalfError struct {
	errorCode uint
	errorMsg  string
}

func NewGandalfError(errorCode uint, errorMsg string) *GandalfError {
	return &GandalfError{
		errorCode: errorCode,
		errorMsg:  errorMsg,
	}
}

func (ge *GandalfError) Error() string {
	return fmt.Sprintf(
		"GandalfError(%s): %s",
		ge.ErrorCodeStr(), ge.errorMsg)
}

func (ge *GandalfError) ErrorMsg() string {
	return ge.errorMsg
}

func (ge *GandalfError) ErrorCode() uint {
	return ge.errorCode
}

func (ge *GandalfError) ErrorCodeStr() string {
	if ge.errorCode == KLowerBidAmount {
		return "KLowerBidAmount"
	} else if ge.errorCode == KInvalidBidAmount {
		return "KInvalidBidAmount"
	} else if ge.errorCode == KTimeout {
		return "KTimeout"
	} else if ge.errorCode == KInvalidItem {
		return "KInvalidItem"
	} else if ge.errorCode == KInvalidUser {
		return "KInvalidUser"
	} else if ge.errorCode == KInvalidSupplier {
		return "KInvalidSupplier"
	} else if ge.errorCode == KAuctionExpired {
		return "KAuctionExpired"
	} else if ge.errorCode == KInvalidBidQuantity {
		return "KInvalidBidQuantity"
	} else if ge.errorCode == KGandalfBackendError {
		return "KGandalfBackendError"
	} else {
		glog.Fatalf("Invalid error code: %d", ge.errorCode)
	}
	return ""
}
