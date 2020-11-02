package aragorn

import (
	"errors"
	"farmtotable/gandalf"
	"farmtotable/misc"
	"firebase.google.com/go"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Aragorn struct {
	gandalf     *gandalf.Gandalf
	firebaseApp *firebase.App
	authCache   interface{}
	logger      *zap.Logger
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
	aragorn.logger = misc.NewLogger()
	return aragorn
}

func (aragorn *Aragorn) Run() {
	r := gin.Default()
	r.POST("/", aragorn.GetUser)
	// User APIs.
	r.POST("/api/v1/resources/users/fetch", aragorn.GetUser)
	r.POST("/api/v1/resources/users/register", aragorn.RegisterUser)

	// Supplier APIs.
	r.POST("/api/v1/resources/suppliers/fetch_all", aragorn.GetAllSuppliers) // Administrator API. Returns all the suppliers.
	r.POST("/api/v1/resources/suppliers/register", aragorn.RegisterSupplier) // Administrator API. // Register Supplier.
	r.POST("/api/v1/resources/suppliers/fetch", aragorn.GetSupplier)         // Administrator API. Gets the supplier info.

	// Item APIs.
	r.POST("/api/v1/resources/items/fetch", aragorn.GetSupplierItems) // Administrator API. Gets all items by a supplier.
	r.POST("/api/v1/resources/items/register", aragorn.RegisterItem)  // Administrator API. Registers item.
	r.POST("/api/v1/resources/items/remove", aragorn.RemoveItem)      // Administrator API. Removes item

	// Auction APIs.
	r.POST("/api/v1/resources/auctions/fetch_all", aragorn.GetAllAuctions)             // Returns all the live auctions.
	r.POST("/api/v1/resources/auctions/fetch_max_bids", aragorn.GetMaxBids)            // Returns the max bids for all requested items so far.
	r.POST("/api/v1/resources/auctions/register_bid", aragorn.RegisterBid)             // Registers a new bid by the user.
	r.POST("/api/v1/resources/auctions/fetch_all_user_bids", aragorn.FetchAllUserBids) // Fetches all the user bids.
	r.POST("/api/v1/resources/auctions/fetch_user_bids", aragorn.FetchUserBidsForItem) // Fetches user bids for an item

	// Order APIs.
	r.POST("/api/v1/resources/orders/get_user_orders", aragorn.GetUserOrders)                            // User and Administrator API.
	r.POST("/api/v1/resources/orders/get_payment_pending_orders", aragorn.GetUserPaymentPendingOrders)   // User and Administrator API.
	r.POST("/api/v1/resources/orders/get_delivery_pending_orders", aragorn.GetUserDeliveryPendingOrders) // User and Administrator API.
	r.POST("/api/v1/resources/orders/get_order", aragorn.GetOrder)                                       // User and Administrator API.
	r.POST("/api/v1/resources/orders/update_order", aragorn.UpdateOrder)                                 // Administrator API.
	r.POST("/api/v1/resources/orders/purchase", aragorn.PurchaseOrder)                                   // User API.

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
/* User APIs. */
func (aragorn *Aragorn) GetUser(c *gin.Context) {
	var response GetUserRet
	var arg GetUserArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, response)
		aragorn.logger.Error(fmt.Sprintf("Invalid input json while fetching user"))
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
	startTime := time.Now()
	err := aragorn.gandalf.RegisterUser(userArg.UserID, userArg.Name, userArg.EmailID, userArg.PhNum, userArg.Address)
	if err != nil {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = fmt.Sprintf("Error while registering user: %v", err)
		c.JSON(http.StatusBadRequest, response)
		aragorn.logger.Error(response.ErrorMsg)
		return
	}
	fmt.Println(fmt.Sprintf("Elapsed Time: %v", time.Since(startTime)))
	response.Status = http.StatusOK
	response.ErrorMsg = ""
	retData := RegistrationStatusRet{
		RegistrationStatus: true,
	}
	response.Data = retData
	c.JSON(http.StatusOK, response)
}

/* Supplier APIs. */
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
		aragorn.logger.Error(ret.ErrorMsg)
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
		aragorn.logger.Error(ret.ErrorMsg)
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
		aragorn.logger.Error(ret.ErrorMsg)
		return
	}
	ret.Status = http.StatusOK
	ret.ErrorMsg = ""
	ret.Data = supplier
	c.JSON(http.StatusOK, ret)
}

/* Item APIs. */
func (aragorn *Aragorn) GetSupplierItems(c *gin.Context) {
	var ret GetSupplierItemsRet
	var arg GetSupplierItemsArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, ret)
		aragorn.logger.Error(ret.ErrorMsg)
		return
	}
	items, err := aragorn.gandalf.GetSupplierItems(arg.SupplierID)
	if err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = fmt.Sprintf("Error while fetching supplier items for supplier: %s", arg.SupplierID)
		c.JSON(http.StatusBadRequest, ret)
		aragorn.logger.Error(ret.ErrorMsg)
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
		aragorn.logger.Error(response.ErrorMsg)
		return
	}
	err := aragorn.gandalf.RegisterItem(item.SupplierID, item.ItemName, item.ItemDescription, item.ItemQty,
		item.AuctionStartDate, item.MinPrice)
	if err != nil {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = fmt.Sprintf("Error while registering item: %v", err)
		c.JSON(http.StatusBadRequest, response)
		aragorn.logger.Error(response.ErrorMsg)
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
		aragorn.logger.Error(ret.ErrorMsg)
		return
	}
	item := aragorn.gandalf.GetItem(arg.ItemID)
	if item.ItemName == "" {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = fmt.Sprintf("Did not find item to delete")
		c.JSON(http.StatusBadRequest, ret)
		aragorn.logger.Error(ret.ErrorMsg)
		return
	}
	err := aragorn.gandalf.DeleteItem(arg.ItemID)
	if err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = fmt.Sprintf("Error while removing item: %v", err)
		c.JSON(http.StatusBadRequest, ret)
		aragorn.logger.Error(ret.ErrorMsg)
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

/* Auction APIs. */
func (aragorn *Aragorn) GetAllAuctions(c *gin.Context) {
	var response FetchAllAuctionsRet
	var fetchAucArg FetchAllAuctionsArg
	if err := c.ShouldBindJSON(&fetchAucArg); err != nil {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, response)
		aragorn.logger.Error(response.ErrorMsg)
		return
	}
	auctions, err := aragorn.gandalf.GetAllAuctions(fetchAucArg.StartID, fetchAucArg.NumAuctions)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.ErrorMsg = "Unable to get all auctions"
		c.JSON(http.StatusInternalServerError, response)
		aragorn.logger.Error(response.ErrorMsg)
		return
	}
	response.Status = http.StatusOK
	response.ErrorMsg = ""
	response.Data = auctions
	c.JSON(http.StatusOK, response)
}

func (aragorn *Aragorn) GetMaxBids(c *gin.Context) {
	var response FetchMaxBidsRet
	var arg FetchMaxBidsArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, response)
		aragorn.logger.Error(response.ErrorMsg)
		return
	}
	auctions, err := aragorn.gandalf.GetMaxBids(arg.ItemIDs)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.ErrorMsg = "Unable to get max bids for items"
		c.JSON(http.StatusInternalServerError, response)
		aragorn.logger.Error(response.ErrorMsg)
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
		aragorn.logger.Error(ret.ErrorMsg)
		return
	}
	err := aragorn.gandalf.RegisterBid(arg.ItemID, arg.UserID, arg.BidAmount, arg.BidQty)
	if err != nil {
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = "Unable to register bid"
		c.JSON(http.StatusInternalServerError, ret)
		aragorn.logger.Error(ret.ErrorMsg)
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

func (aragorn *Aragorn) FetchUserBidsForItem(c *gin.Context) {
	var ret GetUserBidsRet
	var arg GetUserBidsForItemArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, ret)
		aragorn.logger.Error(ret.ErrorMsg)
		return
	}
	bids, err := aragorn.gandalf.GetUserBids(arg.UserID)
	if err != nil {
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = "Unable to fetch user bids"
		c.JSON(http.StatusInternalServerError, ret)
		aragorn.logger.Error(ret.ErrorMsg)
		return
	}
	var itemBids []gandalf.Bid
	for ii := 0; ii < len(bids); ii++ {
		if bids[ii].ItemID == arg.ItemID {
			itemBids = append(itemBids, bids[ii])
		}
	}
	ret.Status = http.StatusOK
	ret.ErrorMsg = ""
	ret.Data = itemBids
	c.JSON(http.StatusOK, ret)
	return
}

func (aragorn *Aragorn) FetchAllUserBids(c *gin.Context) {
	var ret GetUserBidsRet
	var arg GetAllUserBidsArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, ret)
		aragorn.logger.Error(ret.ErrorMsg)
		return
	}
	bids, err := aragorn.gandalf.GetUserBids(arg.UserID)
	if err != nil {
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = "Unable to fetch all user bids"
		c.JSON(http.StatusInternalServerError, ret)
		aragorn.logger.Error(ret.ErrorMsg)
		return
	}
	ret.Status = http.StatusOK
	ret.ErrorMsg = ""
	ret.Data = bids
	c.JSON(http.StatusOK, ret)
	return
}

/* Order APIs. */
func (aragorn *Aragorn) GetUserOrders(c *gin.Context) {
	var ret GetUserOrdersRet
	var arg GetUserOrdersArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, ret)
		aragorn.logger.Error(ret.ErrorMsg)
		return
	}
	orders, err := aragorn.gandalf.GetUserOrders(arg.UserID)
	if err != nil {
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = "Unable to fetch user orders"
		c.JSON(http.StatusInternalServerError, ret)
		aragorn.logger.Error(ret.ErrorMsg)
		return
	}
	orderRets, err := aragorn.joinOrderWithItemInfo(orders)
	if err != nil {
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = "Unable to fetch user orders"
		c.JSON(http.StatusInternalServerError, ret)
		aragorn.logger.Error(ret.ErrorMsg)
		return
	}
	ret.Status = http.StatusOK
	ret.ErrorMsg = ""
	ret.Data = orderRets
	c.JSON(http.StatusOK, ret)
	return
}

func (aragorn *Aragorn) GetUserPaymentPendingOrders(c *gin.Context) {
	var ret GetUserOrdersRet
	var arg GetUserOrdersArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, ret)
		aragorn.logger.Error(ret.ErrorMsg)
		return
	}
	orders, err := aragorn.gandalf.GetUserPaymentPendingOrders(arg.UserID)
	if err != nil {
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = "Unable to fetch user orders"
		c.JSON(http.StatusInternalServerError, ret)
		aragorn.logger.Error(ret.ErrorMsg)
		return
	}
	orderRets, err := aragorn.joinOrderWithItemInfo(orders)
	if err != nil {
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = "Unable to fetch user orders"
		c.JSON(http.StatusInternalServerError, ret)
		aragorn.logger.Error(ret.ErrorMsg)
		return
	}
	ret.Status = http.StatusOK
	ret.ErrorMsg = ""
	ret.Data = orderRets
	c.JSON(http.StatusOK, ret)
	return
}

func (aragorn *Aragorn) GetUserDeliveryPendingOrders(c *gin.Context) {
	var ret GetUserOrdersRet
	var arg GetUserOrdersArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, ret)
		aragorn.logger.Error(ret.ErrorMsg)
		return
	}
	orders, err := aragorn.gandalf.GetUserDeliveryPendingOrders(arg.UserID)
	if err != nil {
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = "Unable to fetch user delivery pending orders"
		c.JSON(http.StatusInternalServerError, ret)
		aragorn.logger.Error(ret.ErrorMsg)
		return
	}
	orderRets, err := aragorn.joinOrderWithItemInfo(orders)
	if err != nil {
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = "Unable to fetch user delivery pending orders"
		c.JSON(http.StatusInternalServerError, ret)
		aragorn.logger.Error(ret.ErrorMsg)
		return
	}
	ret.Status = http.StatusOK
	ret.ErrorMsg = ""
	ret.Data = orderRets
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
		aragorn.logger.Error(ret.ErrorMsg)
		return
	}
	order, err := aragorn.gandalf.GetOrder(arg.OrderID)
	if err != nil {
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = "Unable to fetch order"
		c.JSON(http.StatusInternalServerError, ret)
		aragorn.logger.Error(ret.ErrorMsg)
		return
	}
	var orders []gandalf.Order
	orders = append(orders, order)
	orderRets, err := aragorn.joinOrderWithItemInfo(orders)
	if err != nil {
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = "Unable to fetch order"
		c.JSON(http.StatusInternalServerError, ret)
		aragorn.logger.Error(ret.ErrorMsg)
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
		aragorn.logger.Error(ret.ErrorMsg)
		return
	}
	err := aragorn.gandalf.UpdateOrderStatus(arg.OrderID, arg.Status)
	if err != nil {
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = "Unable to update order"
		c.JSON(http.StatusInternalServerError, ret)
		aragorn.logger.Error(ret.ErrorMsg)
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

func (aragorn *Aragorn) joinOrderWithItemInfo(orders []gandalf.Order) ([]OrderRet, error) {
	itemIDs := make([]string, 0, len(orders))
	orderItems := make([]OrderRet, 0, len(orders))
	for ii := 0; ii < len(orders); ii++ {
		itemIDs = append(itemIDs, orders[ii].ItemID)
	}
	items := aragorn.gandalf.GetItems(itemIDs)
	if len(items) != len(orders) {
		return orderItems, errors.New("unable to get item info for all orders")
	}
	for ii := 0; ii < len(orders); ii++ {
		var orderRet OrderRet
		orderRet.OrderID = orders[ii].OrderID
		orderRet.ItemID = orders[ii].ItemID
		orderRet.ItemName = items[ii].ItemName
		orderRet.ItemDescription = items[ii].ItemDescription
		orderRet.UserID = orders[ii].UserID
		orderRet.ItemQty = orders[ii].ItemQty
		orderRet.ItemPrice = orders[ii].ItemPrice
		orderRet.Status = orders[ii].Status
		orderItems = append(orderItems, orderRet)
	}
	return orderItems, nil
}
