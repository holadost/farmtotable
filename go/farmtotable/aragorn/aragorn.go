package aragorn

import (
	"errors"
	"farmtotable/gandalf"
	"farmtotable/util"
	"firebase.google.com/go"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Aragorn struct {
	gandalf     *gandalf.Gandalf
	firebaseApp *firebase.App
	authCache   interface{}
	apiLogger   *zap.Logger
}

func NewAragorn() *Aragorn {
	aragorn := &Aragorn{}
	// TODO: Populate the config after we have the file from Raunaq.
	//_, err := firebase.NewApp(context.Background(), nil)
	//if err != nil {
	//	panic("Unable to initialize firebase app")
	//}
	// TODO: Pick the backend type based on env. For now hardcode to sqlite.
	aragorn.gandalf = gandalf.NewSqliteGandalf()
	aragorn.apiLogger = util.NewJSONLogger()
	return aragorn
}

func NewAragornWithGandalf(g *gandalf.Gandalf) *Aragorn {
	aragorn := &Aragorn{}
	// TODO: Populate the config after we have the file from Raunaq.
	//_, err := firebase.NewApp(context.Background(), nil)
	//if err != nil {
	//	panic("Unable to initialize firebase app")
	//}
	// TODO: Pick the backend type based on env. For now hardcode to sqlite.
	aragorn.gandalf = g
	aragorn.apiLogger = util.NewJSONLogger()
	return aragorn
}

func (aragorn *Aragorn) Run() {
	glog.Info("Starting Aragorn")
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "Cache-Control"},
		MaxAge:       12 * time.Hour,
	}))

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

	// AuctionModel APIs.
	r.POST("/api/v1/resources/auctions/fetch_all", aragorn.GetAllAuctions)  // Returns all the live auctions.
	r.POST("/api/v1/resources/auctions/fetch_max_bids", aragorn.GetMaxBids) // Returns the max bids for all requested items so far.
	r.POST("/api/v1/resources/auctions/register_bid", aragorn.RegisterBid)  // Registers a new bid by the user.

	// OrderModel APIs.
	r.POST("/api/v1/resources/orders/get_order", aragorn.GetOrder)                                   // UserModel and Administrator API.
	r.POST("/api/v1/resources/orders/get_user_orders", aragorn.GetUserOrders)                        // UserModel and Administrator API.
	r.POST("/api/v1/resources/orders/get_payment_pending_orders", aragorn.GetPaymentPendingOrders)   // Administrator API.
	r.POST("/api/v1/resources/orders/get_delivery_pending_orders", aragorn.GetDeliveryPendingOrders) // Administrator API.
	r.POST("/api/v1/resources/orders/update_order", aragorn.UpdateOrder)                             // Administrator API.
	r.POST("/api/v1/resources/orders/purchase", aragorn.PurchaseOrder)                               // UserModel API.
	r.POST("/api/v1/resources/test/orders/test_only_add_order", aragorn.TestOnlyAddOrder)            // Test API.
	r.POST("/api/v1/resources/test/auctions/test_only_add_auctions", aragorn.TestOnlyAddAuctions)    // Test API.

	// Start router.
	r.Run(":8080")
}

//func (aragorn *Aragorn) authenticate(c *gin.Context) (string, error) {
//	idToken := c.Request.Header["Authorization"][0]
//	client, err := aragorn.firebaseApp.Auth(context.Background())
//	if err != nil {
//		return "", err
//	}
//	token, err := client.VerifyIDToken(context.Background(), idToken)
//	if err != nil {
//		return "", err
//	}
//	return token.UID, nil
//}
//
//func (aragorn *Aragorn) doesUserExist(userID string) bool {
//	user := aragorn.gandalf.GetUserByID(userID)
//	if user.UserID != userID {
//		return false
//	}
//	return true
//}

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
	err := aragorn.gandalf.RegisterItem(item.SupplierID, item.ItemName, item.ItemDescription, item.ItemQty,
		item.AuctionStartDate, item.MinPrice, item.AuctionDurationSecs)
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

/* AuctionModel APIs. */
func (aragorn *Aragorn) GetAllAuctions(c *gin.Context) {
	var response FetchAllAuctionsRet
	var fetchAucArg FetchAllAuctionsArg
	if err := c.ShouldBindJSON(&fetchAucArg); err != nil {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, response)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", response.ErrorMsg, err))
		return
	}
	auctions, err := aragorn.gandalf.GetAllAuctions(fetchAucArg.StartID, fetchAucArg.NumAuctions)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.ErrorMsg = "Unable to get all auctions"
		c.JSON(http.StatusInternalServerError, response)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", response.ErrorMsg, err))
		return
	}
	var retAuctions []gandalf.AuctionModel
	if len(auctions) == 0 {
		response.Status = http.StatusOK
		response.ErrorMsg = ""
		response.Data = FetchAllAuctionsRetData{
			Auctions: auctions,
			NextID:   -1,
		}
		c.JSON(http.StatusOK, response)
		return
	}

	var maxID uint
	maxID = 0
	for _, auction := range auctions {
		if auction.ID > maxID {
			maxID = auction.ID
		}
		deadline := auction.AuctionStartTime.Add(time.Second * time.Duration(int64(auction.AuctionDurationSecs)))
		if deadline.Before(time.Now()) {
			// The deadline has expired. Skip this auction.
			continue
		}
		retAuctions = append(retAuctions, auction)
	}

	response.Status = http.StatusOK
	response.ErrorMsg = ""
	response.Data = FetchAllAuctionsRetData{
		Auctions: retAuctions,
		NextID:   int64(maxID + 1),
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
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = "Unable to register bid"
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

func (aragorn *Aragorn) TestOnlyAddAuctions(c *gin.Context) {
	// This method takes all the items in items table and adds them to the auctions table.
	var ret TestOnlyAddAuctionsRet
	suppliers, err := aragorn.gandalf.GetAllSuppliers()
	if err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Unable to get all suppliers to find items to add to auctions"
		c.JSON(http.StatusBadRequest, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	success := false
	var auctions []gandalf.AuctionModel
	if len(suppliers) == 0 {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Did not find any suppliers"
		c.JSON(http.StatusBadRequest, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	for _, supplier := range suppliers {
		items, err := aragorn.gandalf.GetSupplierItems(supplier.SupplierID)
		if err != nil {
			continue
		}
		if len(items) != 0 {
			for _, item := range items {
				var auction gandalf.AuctionModel
				auction.ItemID = item.ItemID
				auction.ItemQty = item.ItemQty
				auction.ItemName = item.ItemName
				auction.AuctionStartTime = item.AuctionStartTime
				auction.AuctionDurationSecs = 86400
				auctions = append(auctions, auction)
			}
			success = true
		}
	}
	if !success {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Did not find any items to add to auctions"
		c.JSON(http.StatusBadRequest, ret)
		aragorn.apiLogger.Error(fmt.Sprintf("%s: error: %v", ret.ErrorMsg, err))
		return
	}
	err = aragorn.gandalf.AddAuctions(auctions)
	if err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = fmt.Sprintf("Failed to register auctions due to error: %v", err)
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
		orderRet.ItemDescription = item.ItemDescription
		orderRet.UserID = orders[ii].UserID
		orderRet.ItemQty = orders[ii].ItemQty
		orderRet.ItemPrice = orders[ii].ItemPrice
		orderRet.Status = orders[ii].Status
		orderItems = append(orderItems, orderRet)
	}
	return orderItems, nil
}
