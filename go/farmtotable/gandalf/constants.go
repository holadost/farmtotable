package gandalf

const (
	PGHost     = "localhost"
	PGPort     = 5432
	PGUser     = "postgres"
	PGPassword = "farmToTable4u"
	PGDbName   = "nikhil"
)

const (
	SQLiteDBPath = "/tmp/gandalf.db"
)

const (
	KOrderPaymentPending     = 0
	KOrderPaymentPendingStr  = "KOrderPaymentPending"
	KOrderDeliveryPending    = 1
	KOrderDeliveryPendingStr = "KOrderDeliveryPending"
	KOrderComplete           = 2
	KOrderCompleteStr        = "KOrderComplete"
	KOrderCancelled          = 3
	KOrderCancelledStr       = "KOrderCancelled"
)
