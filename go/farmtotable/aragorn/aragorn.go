package aragorn

import (
	"farmtotable/gandalf"
	"firebase.google.com/go"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Aragorn struct {
	gandalf     *gandalf.Gandalf
	firebaseApp *firebase.App
	authCache   interface{}
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
	return aragorn
}

func (aragorn *Aragorn) Run() {
	r := gin.Default()
	// User APIs.
	r.POST("/api/v1/resources/users/fetch", aragorn.getUser)
	r.POST("/api/v1/resources/users/register", aragorn.registerUser)

	// Supplier APIs.
	r.POST("/api/v1/resources/suppliers/fetch_all", aragorn.getAllSuppliers) // Administrator API. Returns all the suppliers.
	r.POST("/api/v1/resources/suppliers/register", aragorn.registerSupplier) // Administrator API. // Register Supplier.
	r.POST("/api/v1/resources/suppliers/fetch", aragorn.getSupplier)         // Administrator API. Gets the supplier info.

	// Item APIs.
	r.POST("/api/v1/resources/items/fetch", aragorn.getSupplierItems) // Administrator API. Gets all items by a supplier.
	r.POST("/api/v1/resources/items/register", aragorn.registerItem)  // Administrator API. Registers item.
	r.POST("/api/v1/resources/items/remove", aragorn.removeItem)      // Administrator API. Removes item

	// Auction APIs.
	r.POST("/api/v1/resources/auctions/fetch_all", aragorn.getAllAuctions)      // Returns all the live auctions.
	r.POST("/api/v1/resources/auctions/fetch_max_bids", aragorn.getMaxBids)     // Returns the max bids for all requested items so far.
	r.POST("/api/v1/resources/auctions/register_bid", aragorn.registerBid)      // Registers a new bid by the user.
	r.POST("/api/v1/resources/auctions/fetch_user_bids", aragorn.fetchUserBids) // Registers a new bid by the user.

	// Order APIs.
	//r.GET("/api/v1/resources/orders/getUserOrders", aragorn.getUserOrders) // Administrator API.
	//r.GET("/api/v1/resources/orders/getPaymentPendingOrders", aragorn.getPaymentPendingOrders) // Administrator API.
	//r.GET("/api/v1/resources/orders/getDeliveryPendingOrders", aragorn.getDeliveryPendingOrders) // Administrator API.
	//r.GET("/api/v1/resources/orders/getOrder", aragorn.getOrder) // Administrator API.
	//r.POST("/api/v1/resources/orders/updateOrder", aragorn.updateOrder) // Administrator API.

	// Start router.
	r.Run("localhost:8080")
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
func (aragorn *Aragorn) getUser(c *gin.Context) {
	var response GetUserRet
	var arg GetUserArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, response)
		return
	}
	fullUser := aragorn.gandalf.GetUserByID(arg.UserID)
	response.Status = http.StatusOK
	response.ErrorMsg = ""
	response.Data = fullUser
	c.JSON(http.StatusOK, response)
}

func (aragorn *Aragorn) registerUser(c *gin.Context) {
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
func (aragorn *Aragorn) registerSupplier(c *gin.Context) {
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
		return
	}
	ret.Status = http.StatusOK
	ret.ErrorMsg = ""
	retData := RegistrationStatusRet{}
	retData.RegistrationStatus = true
	ret.Data = retData
	c.JSON(http.StatusOK, ret)

}

func (aragorn *Aragorn) getAllSuppliers(c *gin.Context) {
	var ret GetAllSuppliersRet
	suppliers, err := aragorn.gandalf.GetAllSuppliers()
	if err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = fmt.Sprintf("Error while registering supplier: %v", err)
		c.JSON(http.StatusBadRequest, ret)
		return
	}
	ret.Status = http.StatusOK
	ret.ErrorMsg = ""
	ret.Data = suppliers
	c.JSON(http.StatusOK, ret)
}

func (aragorn *Aragorn) getSupplier(c *gin.Context) {
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
		return
	}
	ret.Status = http.StatusOK
	ret.ErrorMsg = ""
	ret.Data = supplier
	c.JSON(http.StatusOK, ret)
}

/* Item APIs. */
func (aragorn *Aragorn) getSupplierItems(c *gin.Context) {
	var ret GetSupplierItemsRet
	var arg GetSupplierItemsArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, ret)
		return
	}
	items, err := aragorn.gandalf.GetSupplierItems(arg.SupplierID)
	if err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = fmt.Sprintf("Error while fetching supplier items for supplier: %s", arg.SupplierID)
		c.JSON(http.StatusBadRequest, ret)
		return
	}
	ret.Status = http.StatusOK
	ret.ErrorMsg = ""
	ret.Data = items
	c.JSON(http.StatusOK, ret)
}

func (aragorn *Aragorn) registerItem(c *gin.Context) {
	var response RegisterItemRet
	var item RegisterItemArg
	if err := c.ShouldBindJSON(&item); err != nil {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, response)
		return
	}
	err := aragorn.gandalf.RegisterItem(item.SupplierID, item.ItemName, item.ItemDescription, item.ItemQty,
		item.AuctionStartDate, item.MinPrice)
	if err != nil {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = fmt.Sprintf("Error while registering item: %v", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response.Status = http.StatusOK
	response.ErrorMsg = ""
	retData := RegistrationStatusRet{}
	retData.RegistrationStatus = true
	response.Data = retData
	c.JSON(http.StatusOK, response)
}

func (aragorn *Aragorn) removeItem(c *gin.Context) {
	var ret RemoveItemRet
	var arg RemoveItemArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, ret)
		return
	}
	item := aragorn.gandalf.GetItem(arg.ItemID)
	if item.ItemName == "" {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = fmt.Sprintf("Did not find item to delete")
		c.JSON(http.StatusBadRequest, ret)
		return
	}
	err := aragorn.gandalf.DeleteItem(arg.ItemID)
	if err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = fmt.Sprintf("Error while removing item: %v", err)
		c.JSON(http.StatusBadRequest, ret)
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
func (aragorn *Aragorn) getAllAuctions(c *gin.Context) {
	var response FetchAllAuctionsRet
	var fetchAucArg FetchAllAuctionsArg
	if err := c.ShouldBindJSON(&fetchAucArg); err != nil {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, response)
		return
	}
	auctions, err := aragorn.gandalf.GetAllAuctions(fetchAucArg.StartID, fetchAucArg.NumAuctions)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.ErrorMsg = "Unable to fetch user items"
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response.Status = http.StatusOK
	response.ErrorMsg = ""
	response.Data = auctions
	c.JSON(http.StatusOK, response)
}

func (aragorn *Aragorn) getMaxBids(c *gin.Context) {
	var response FetchMaxBidsRet
	var arg FetchMaxBidsArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, response)
		return
	}
	auctions, err := aragorn.gandalf.GetMaxBids(arg.ItemIDs)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.ErrorMsg = "Unable to fetch user items"
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	var results []FetchMaxBidsRetData
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

func (aragorn *Aragorn) registerBid(c *gin.Context) {
	var arg RegisterBidArg
	var ret RegisterBidRet
	if err := c.ShouldBindJSON(&arg); err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, ret)
		return
	}
	err := aragorn.gandalf.RegisterBid(arg.ItemID, arg.UserID, arg.BidAmount, arg.BidQty)
	if err != nil {
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = "Unable to register bid"
		c.JSON(http.StatusInternalServerError, ret)
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

func (aragorn *Aragorn) fetchUserBids(c *gin.Context) {
	var ret GetUserBidsRet
	var arg GetUserBidsArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		ret.Status = http.StatusBadRequest
		ret.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, ret)
		return
	}
	auctions, err := aragorn.gandalf.GetUserAuctions(arg.UserID)
	if err != nil {
		ret.Status = http.StatusInternalServerError
		ret.ErrorMsg = "Unable to fetch user bids"
		c.JSON(http.StatusInternalServerError, ret)
		return
	}
	ret.Status = http.StatusOK
	ret.ErrorMsg = ""
	ret.Data = auctions
	c.JSON(http.StatusOK, ret)
	return
}

/* Order APIs. */
