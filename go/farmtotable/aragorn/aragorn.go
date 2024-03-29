package aragorn

import (
	"errors"
	"farmtotable/gandalf"
	"farmtotable/util"
	"flag"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var (
	fbCredPath = flag.String("fb_cred_path", "",
		"The firebase credentials path")
	skipAuth = flag.Bool("skip_auth", true,
		"Flag to skip auth. Defaults to true for dev and test purposes")
)

type Aragorn struct {
	gandalf   *gandalf.Gandalf
	auth      *Auth
	apiLogger *zap.Logger
}

func NewAragorn() *Aragorn {
	aragorn := &Aragorn{}
	aragorn.auth = NewAuth(*fbCredPath)
	aragorn.gandalf = gandalf.NewSqliteGandalf()
	aragorn.apiLogger = util.NewJSONLogger()
	return aragorn
}

func NewAragornWithGandalf(g *gandalf.Gandalf) *Aragorn {
	aragorn := &Aragorn{}
	aragorn.auth = NewAuth(*fbCredPath)
	aragorn.gandalf = g
	aragorn.apiLogger = util.NewJSONLogger()
	return aragorn
}

func (aragorn *Aragorn) Run() {
	glog.Info("Aragorn initialized")
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "Cache-Control"},
		MaxAge:       12 * time.Hour,
	}))
	r.Use(aragorn.authenticate)

	// UserModel APIs.
	r.POST("/api/v1/resources/users/fetch", aragorn.GetUser)
	r.POST("/api/v1/resources/users/register", aragorn.RegisterUser)

	// SupplierModel APIs.
	r.POST("/api/v1/resources/suppliers/fetch_all", aragorn.GetAllSuppliers) // Administrator API. Returns all the suppliers.
	r.POST("/api/v1/resources/suppliers/register", aragorn.RegisterSupplier) // Administrator API. // Register SupplierModel.
	r.POST("/api/v1/resources/suppliers/fetch", aragorn.GetSupplier)         // Administrator API. Gets the supplier info.

	// ItemModel APIs.
	r.POST("/api/v1/resources/items/fetch", aragorn.GetSupplierItems) // Administrator API. Gets all items by a supplier.
	r.POST("/api/v1/resources/items/register", aragorn.RegisterItem)  // Administrator API. Registers item.
	r.POST("/api/v1/resources/items/remove", aragorn.RemoveItem)      // Administrator API. Removes item
	r.POST("/api/v1/resources/items/get", aragorn.GetItem)            // Gets an item.

	// AuctionModel APIs.
	r.POST("/api/v1/resources/auctions/fetch_all", aragorn.FetchAuctions)   // Returns all the live auctions.
	r.POST("/api/v1/resources/auctions/fetch_max_bids", aragorn.GetMaxBids) // Returns the max bids for all requested items so far.
	r.POST("/api/v1/resources/auctions/register_bid", aragorn.RegisterBid)  // Registers a new bid by the user.
	r.POST("/api/v1/resources/auctions/get_user_bid", aragorn.GetUserBid)   // Gets the bid registered by the user.

	// OrderModel APIs.
	r.POST("/api/v1/resources/orders/get_order", aragorn.GetOrder)                                   // UserModel and Administrator API.
	r.POST("/api/v1/resources/orders/get_user_orders", aragorn.GetUserOrders)                        // UserModel and Administrator API.
	r.POST("/api/v1/resources/orders/get_payment_pending_orders", aragorn.GetPaymentPendingOrders)   // Administrator API.
	r.POST("/api/v1/resources/orders/get_delivery_pending_orders", aragorn.GetDeliveryPendingOrders) // Administrator API.
	r.POST("/api/v1/resources/orders/update_order", aragorn.UpdateOrder)                             // Administrator API.
	r.POST("/api/v1/resources/orders/purchase", aragorn.PurchaseOrder)                               // UserModel API.
	r.POST("/api/v1/resources/test/orders/test_only_add_order", aragorn.TestOnlyAddOrder)            // Test API.

	// Start router.
	r.Run(":8080")
}

func (aragorn *Aragorn) authenticate(c *gin.Context) {
	if *skipAuth {
		c.Next()
		return
	}
	if aragorn.auth.Authenticate(c) != nil {
		c.Abort()
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Writer.Write([]byte("Unauthorized"))
		return
	}
	c.Next()
}

/* UserModel APIs. */
func (aragorn *Aragorn) GetUser(c *gin.Context) {
	var response GetUserRet
	var arg GetUserArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, response)
		aragorn.apiLogger.Error(fmt.Sprintf("Invalid input json while fetching user"))
		return
	}
	fullUser := aragorn.gandalf.GetUserByID(arg.UserID)
	response.Status = http.StatusOK
	response.ErrorMsg = ""
	response.Data = fullUser
	c.JSON(http.StatusOK, response)
}

func (aragorn *Aragorn) RegisterUser(c *gin.Context) {
	var response RegisterUserRet
	var userArg RegisterUserArg
	if err := c.ShouldBindJSON(&userArg); err != nil {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, response)
		return
	}
	err := aragorn.gandalf.RegisterUser(userArg.UserID, userArg.Name, userArg.EmailID, userArg.PhNum, userArg.Address)
	if err != nil {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = fmt.Sprintf("Error while registering user: %v", err)
		c.JSON(http.StatusBadRequest, response)
		aragorn.apiLogger.Error(response.ErrorMsg)
		return
	}
	response.Status = http.StatusOK
	response.ErrorMsg = ""
	retData := RegistrationStatusRet{
		RegistrationStatus: true,
	}
	response.Data = retData
	c.JSON(http.StatusOK, response)
}

/* SupplierModel APIs. */
func (aragorn *Aragorn) RegisterSupplier(c *gin.Context) {
	var ret RegisterSupplierRet
	var arg RegisterSupplierArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, ret)
		return
	}
	err := aragorn.gandalf.RegisterSupplier(arg.SupplierName, arg.SupplierEmailID, arg.SupplierPhNum,
		arg.SupplierAddress, arg.SupplierDescription, arg.SupplierTags)
	if err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = fmt.Sprintf("Error while registering supplier: %v", err)
		c.JSON(http.StatusBadRequest, ret)
		aragorn.apiLogger.Error(ret.ErrorMsg)
		return
	}
	ret.Status = http.StatusOK
	ret.ErrorMsg = ""
	retData := RegistrationStatusRet{}
	retData.RegistrationStatus = true
	ret.Data = retData
	c.JSON(http.StatusOK, ret)

}

func (aragorn *Aragorn) GetAllSuppliers(c *gin.Context) {
	var ret GetAllSuppliersRet
	suppliers, err := aragorn.gandalf.GetAllSuppliers()
	if err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = fmt.Sprintf("Error while registering supplier: %v", err)
		c.JSON(http.StatusBadRequest, ret)
		aragorn.apiLogger.Error(ret.ErrorMsg)
		return
	}
	ret.Status = http.StatusOK
	ret.ErrorMsg = ""
	ret.Data = suppliers
	c.JSON(http.StatusOK, ret)
}

func (aragorn *Aragorn) GetSupplier(c *gin.Context) {
	var ret GetSupplierRet
	var arg GetSupplierArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, ret)
		return
	}
	supplier := aragorn.gandalf.GetSupplierByID(arg.SupplierID)
	if supplier.SupplierName == "" {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = fmt.Sprintf("Error while fetching supplier with ID: %s", arg.SupplierID)
		c.JSON(http.StatusBadRequest, ret)
		aragorn.apiLogger.Error(ret.ErrorMsg)
		return
	}
	ret.Status = http.StatusOK
	ret.ErrorMsg = ""
	ret.Data = supplier
	c.JSON(http.StatusOK, ret)
}

/* ItemModel APIs. */
func (aragorn *Aragorn) GetSupplierItems(c *gin.Context) {
	var ret GetSupplierItemsRet
	var arg GetSupplierItemsArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, ret)
		aragorn.apiLogger.Error(ret.ErrorMsg)
		return
	}
	items, err := aragorn.gandalf.GetSupplierItems(arg.SupplierID)
	if err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = fmt.Sprintf("Error while fetching supplier items for supplier: %s", arg.SupplierID)
		c.JSON(http.StatusBadRequest, ret)
		aragorn.apiLogger.Error(ret.ErrorMsg)
		return
	}
	ret.Status = http.StatusOK
	ret.ErrorMsg = ""
	ret.Data = items
	c.JSON(http.StatusOK, ret)
}

func (aragorn *Aragorn) RegisterItem(c *gin.Context) {
	var response RegisterItemRet
	var item RegisterItemArg
	if err := c.ShouldBindJSON(&item); err != nil {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, response)
		aragorn.apiLogger.Error(response.ErrorMsg)
		return
	}
	err := aragorn.gandalf.RegisterItem(
		item.SupplierID, item.ItemName, item.ItemDescription, item.ItemQty,
		item.AuctionStartDate, item.MinPrice, item.AuctionDurationSecs, item.ImageURL,
		item.MinBidQty, item.MaxBidQty, item.ItemUnit)
	if err != nil {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = fmt.Sprintf("Error while registering item: %v", err)
		c.JSON(http.StatusBadRequest, response)
		aragorn.apiLogger.Error(response.ErrorMsg)
		return
	}
	response.Status = http.StatusOK
	response.ErrorMsg = ""
	retData := RegistrationStatusRet{}
	retData.RegistrationStatus = true
	response.Data = retData
	c.JSON(http.StatusOK, response)
}

func (aragorn *Aragorn) RemoveItem(c *gin.Context) {
	var ret RemoveItemRet
	var arg RemoveItemArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	item, err := aragorn.gandalf.GetItem(arg.ItemID)
	if err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = fmt.Sprintf("Unable to get item from backend")
		c.JSON(http.StatusBadRequest, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	if item.ItemName == "" {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = fmt.Sprintf("Did not find item to delete")
		c.JSON(http.StatusBadRequest, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	// TODO: We also need to remove the item from the auctions table and all associated bids.
	err = aragorn.gandalf.DeleteItem(arg.ItemID)
	if err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = fmt.Sprintf("Error while removing item: %v", err)
		c.JSON(http.StatusBadRequest, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	ret.Status = http.StatusOK
	ret.ErrorMsg = ""
	retData := RegistrationStatusRet{
		RegistrationStatus: true,
	}
	ret.Data = retData
	c.JSON(http.StatusOK, ret)
}

func (aragorn *Aragorn) GetItem(c *gin.Context) {
	var ret GetItemRet
	var arg GetItemArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	item, err := aragorn.gandalf.GetItem(arg.ItemID)
	if err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = fmt.Sprintf("Unable to get item from backend")
		c.JSON(http.StatusBadRequest, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	if item.ItemName == "" {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = fmt.Sprintf("Did not find item: %s", arg.ItemID)
		c.JSON(http.StatusBadRequest, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	ret.Status = http.StatusOK
	ret.ErrorMsg = ""
	ret.Data = item
	c.JSON(http.StatusOK, ret)
}

/* AuctionModel APIs. */
func (aragorn *Aragorn) FetchAuctions(c *gin.Context) {
	var response FetchAllAuctionsRet
	var fetchAucArg FetchAllAuctionsArg
	if err := c.ShouldBindJSON(&fetchAucArg); err != nil {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, response)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", response.ErrorMsg, err))
		return
	}
	var retAuctions []gandalf.AuctionModel
	var gatheredAuctions bool
	var nextID uint64
	startID := fetchAucArg.StartID
	for {
		auctions, err := aragorn.gandalf.FetchAuctions(startID, fetchAucArg.NumAuctions)
		if err != nil {
			response.Status = http.StatusInternalServerError
			response.ErrorMsg = "Unable to get all auctions"
			c.JSON(http.StatusInternalServerError, response)
			aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", response.ErrorMsg, err))
			return
		}
		if len(auctions) == 0 {
			response.Status = http.StatusOK
			response.ErrorMsg = ""
			response.Data = FetchAllAuctionsRetData{
				Auctions: retAuctions,
				NextID:   -1,
			}
			c.JSON(http.StatusOK, response)
			return
		}
		for ii, auction := range auctions {
			deadline := auction.AuctionStartTime.Add(time.Second * time.Duration(int64(auction.AuctionDurationSecs)))
			if deadline.Before(time.Now()) {
				// The deadline has expired. Skip this auction.
				continue
			}
			retAuctions = append(retAuctions, auction)
			if len(retAuctions) == int(fetchAucArg.NumAuctions) {
				// We have gathered the requested number of auctions.
				nextID = startID + uint64(ii) + 1
				gatheredAuctions = true
				break
			}
		}
		startID += fetchAucArg.NumAuctions
		if gatheredAuctions {
			break
		}
	}
	response.Status = http.StatusOK
	response.ErrorMsg = ""
	response.Data = FetchAllAuctionsRetData{
		Auctions: retAuctions,
		NextID:   int64(nextID),
	}
	c.JSON(http.StatusOK, response)
}

func (aragorn *Aragorn) GetMaxBids(c *gin.Context) {
	var response FetchMaxBidsRet
	var arg FetchMaxBidsArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, response)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", response.ErrorMsg, err))
		return
	}
	auctions, err := aragorn.gandalf.GetMaxBids(arg.ItemIDs)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.ErrorMsg = "Unable to get max bids for items"
		c.JSON(http.StatusInternalServerError, response)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", response.ErrorMsg, err))
		return
	}
	results := make([]FetchMaxBidsRetData, 0, len(auctions))
	for ii := 0; ii < len(auctions); ii++ {
		itemBid := FetchMaxBidsRetData{
			ItemID: auctions[ii].ItemID,
			MaxBid: auctions[ii].MaxBid,
		}
		results = append(results, itemBid)
	}
	response.Status = http.StatusOK
	response.ErrorMsg = ""
	response.Data = results
	c.JSON(http.StatusOK, response)
}

func (aragorn *Aragorn) GetUserBid(c *gin.Context) {
	var ret GetUserBidRet
	var arg GetUserBidArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, ret)
		return
	}
	bid, err := aragorn.gandalf.GetUserBid(arg.UserID, arg.ItemID)
	if err != nil {
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = fmt.Sprintf("Unable to get user(%s) bid for item: %s", arg.UserID, arg.ItemID)
		c.JSON(http.StatusInternalServerError, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	ret.Status = http.StatusOK
	ret.ErrorMsg = ""
	if bid.ItemID == "" {
		ret.Data = 0
	} else {
		ret.Data = bid.BidAmount
	}
	c.JSON(http.StatusOK, ret)
}

func (aragorn *Aragorn) RegisterBid(c *gin.Context) {
	var arg RegisterBidArg
	var ret RegisterBidRet
	if err := c.ShouldBindJSON(&arg); err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	err := aragorn.gandalf.RegisterBid(arg.ItemID, arg.UserID, arg.BidAmount, arg.BidQty)
	if err != nil {
		ge := err.(*gandalf.GandalfError)
		glog.Errorf("Received error from gandalf: %d", ge.ErrorCode())
		if ge.ErrorCode() == gandalf.KGandalfBackendError || ge.ErrorCode() == gandalf.KTimeout {
			ret.Status = http.StatusInternalServerError
			ret.ErrorMsg = "Unable to register bid"
			c.JSON(http.StatusInternalServerError, ret)
			aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		} else {
			ret.Status = http.StatusBadRequest
			ret.ErrorMsg = ge.ErrorCodeStr()
			c.JSON(http.StatusInternalServerError, ret)
			aragorn.apiLogger.Error(ge.Error())
		}
		return
	}
	ret.Status = http.StatusOK
	ret.ErrorMsg = ""
	retData := RegistrationStatusRet{
		RegistrationStatus: true,
	}
	ret.Data = retData
	c.JSON(http.StatusOK, ret)
	return
}

/* OrderModel APIs. */
func (aragorn *Aragorn) GetUserOrders(c *gin.Context) {
	var ret GetOrdersRet
	var arg GetUserOrdersArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	orders, err := aragorn.gandalf.GetUserOrders(arg.UserID)
	if err != nil {
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = "Unable to fetch user orders"
		c.JSON(http.StatusInternalServerError, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	orderRets, err := aragorn.joinOrderWithItemInfo(orders)
	if err != nil {
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = "Unable to fetch user orders"
		c.JSON(http.StatusInternalServerError, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	ret.Status = http.StatusOK
	ret.ErrorMsg = ""
	ret.Data = GetOrdersRetData{
		Orders: orderRets,
		NextID: 0,
	}
	c.JSON(http.StatusOK, ret)
	return
}

func (aragorn *Aragorn) GetPaymentPendingOrders(c *gin.Context) {
	var ret GetOrdersRet
	var arg ScanOrdersArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	orders, err := aragorn.gandalf.ScanPaymentPendingOrders(arg.StartID, arg.NumOrders)
	if err != nil {
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = "Unable to fetch payment pending orders"
		c.JSON(http.StatusInternalServerError, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	var orderRets []OrderRet
	if len(orders) == 0 {
		ret.Status = http.StatusOK
		ret.ErrorMsg = ""
		ret.Data = GetOrdersRetData{
			Orders: orderRets,
			NextID: -1,
		}
		c.JSON(http.StatusOK, ret)
		return
	}
	orderRets, err = aragorn.joinOrderWithItemInfo(orders)
	if err != nil {
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = "Unable to fetch payment pending orders"
		c.JSON(http.StatusInternalServerError, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	ret.Status = http.StatusOK
	ret.ErrorMsg = ""
	ret.Data = GetOrdersRetData{
		Orders: orderRets,
		NextID: int64(orders[len(orders)-1].ID),
	}
	c.JSON(http.StatusOK, ret)
	return
}

func (aragorn *Aragorn) GetDeliveryPendingOrders(c *gin.Context) {
	var ret GetOrdersRet
	var arg ScanOrdersArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	orders, err := aragorn.gandalf.ScanDeliveryPendingOrders(arg.StartID, arg.NumOrders)
	if err != nil {
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = "Unable to fetch delivery pending orders"
		c.JSON(http.StatusInternalServerError, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}

	var orderRets []OrderRet
	if len(orders) == 0 {
		ret.Status = http.StatusOK
		ret.ErrorMsg = ""
		ret.Data = GetOrdersRetData{
			Orders: orderRets,
			NextID: -1,
		}
		c.JSON(http.StatusOK, ret)
		return
	}

	orderRets, err = aragorn.joinOrderWithItemInfo(orders)
	if err != nil {
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = "Unable to fetch delivery pending orders"
		c.JSON(http.StatusInternalServerError, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	ret.Status = http.StatusOK
	ret.ErrorMsg = ""
	ret.Data = GetOrdersRetData{
		Orders: orderRets,
		NextID: int64(orders[len(orders)-1].ID),
	}
	c.JSON(http.StatusOK, ret)
	return
}

func (aragorn *Aragorn) GetOrder(c *gin.Context) {
	var ret GetOrderRet
	var arg GetOrderArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	order, err := aragorn.gandalf.GetOrder(arg.OrderID)
	if err != nil {
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = "Unable to fetch order"
		c.JSON(http.StatusInternalServerError, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	var orders []gandalf.OrderModel
	orders = append(orders, order)
	orderRets, err := aragorn.joinOrderWithItemInfo(orders)
	if err != nil {
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = "Unable to fetch order"
		c.JSON(http.StatusInternalServerError, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	ret.Status = http.StatusOK
	ret.ErrorMsg = ""
	ret.Data = orderRets[0]
	c.JSON(http.StatusOK, ret)
	return
}

func (aragorn *Aragorn) UpdateOrder(c *gin.Context) {
	var ret UpdateOrderRet
	var arg UpdateOrderArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	status := gandalf.OrderStatusStr(arg.Status)
	statusCode, err := status.ToUint32()
	if err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Invalid input JSON. OrderModel status not supported"
		c.JSON(http.StatusBadRequest, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	err = aragorn.gandalf.UpdateOrderStatus(arg.OrderID, uint32(statusCode))
	if err != nil {
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = "Unable to update order"
		c.JSON(http.StatusInternalServerError, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("Unable to update order with order ID: %s due to error: %v",
			arg.OrderID, err))
		return
	}
	ret.Status = http.StatusOK
	ret.ErrorMsg = ""
	retData := RegistrationStatusRet{
		RegistrationStatus: true,
	}
	ret.Data = retData
	c.JSON(http.StatusOK, ret)
	return
}

func (aragorn *Aragorn) PurchaseOrder(c *gin.Context) {
	// TODO: Still needs to be implemented
	return
}

func (aragorn *Aragorn) TestOnlyAddOrder(c *gin.Context) {
	// This is a test API. This must not be used for any other reason.
	var ret TestOnlyAddOrderRet
	var arg TestOnlyAddOrderArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	var orders []gandalf.OrderModel
	var order gandalf.OrderModel
	order.ItemID = arg.ItemID
	order.UserID = arg.UserID
	order.ItemPrice = arg.ItemPrice
	order.ItemQty = arg.ItemQty
	orders = append(orders, order)
	err := aragorn.gandalf.AddOrders(orders)
	if err != nil {
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = "Unable to add order"
		c.JSON(http.StatusInternalServerError, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	ret.Status = http.StatusOK
	ret.ErrorMsg = ""
	retData := RegistrationStatusRet{
		RegistrationStatus: true,
	}
	ret.Data = retData
	c.JSON(http.StatusOK, ret)
	return
}

func (aragorn *Aragorn) joinOrderWithItemInfo(orders []gandalf.OrderModel) ([]OrderRet, error) {
	itemIDs := make([]string, 0, len(orders))
	orderItems := make([]OrderRet, 0, len(orders))
	for ii := 0; ii < len(orders); ii++ {
		itemIDs = append(itemIDs, orders[ii].ItemID)
	}
	items, err := aragorn.gandalf.GetItems(itemIDs)
	if err != nil {
		glog.Errorf("Unable to get item IDs for orders due to error: %v", err)
	}
	itemSet := make(map[string]gandalf.ItemModel)
	for _, item := range items {
		itemSet[item.ItemID] = item
	}
	for ii := 0; ii < len(orders); ii++ {
		var orderRet OrderRet
		orderRet.OrderID = orders[ii].OrderID
		orderRet.ItemID = orders[ii].ItemID
		item, err := itemSet[orders[ii].ItemID]
		if !err {
			glog.Errorf("Unable to get item information for item: %s", orders[ii].ItemID)
			return orderItems, errors.New("unable to find item for order")
		}
		orderRet.ItemName = item.ItemName
		orderRet.ImageURL = item.ImageURL
		orderRet.UserID = orders[ii].UserID
		orderRet.ItemQty = orders[ii].ItemQty
		orderRet.ItemPrice = orders[ii].ItemPrice
		orderRet.Status = orders[ii].Status
		orderRet.DeliveryPrice = orders[ii].DeliveryPrice
		orderRet.TaxPrice = orders[ii].TaxPrice
		orderRet.TotalPrice = orders[ii].TotalPrice
		orderRet.OrderHistory = orders[ii].OrderHistory
		orderItems = append(orderItems, orderRet)
	}
	return orderItems, nil
}
