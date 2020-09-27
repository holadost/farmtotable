package aragorn

import (
	"context"
	"errors"
	"farmtotable/gandalf"
	"firebase.google.com/go"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Aragorn struct {
	gandalf     *gandalf.Gandalf
	firebaseApp *firebase.App
	authCache   interface{}
}

func NewAragorn() *Aragorn {
	aragorn := &Aragorn{}
	// TODO: Populate the config after we have the file from Raunaq.
	_, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		panic("Unable to initialize firebase app")
	}
	// TODO: Pick the backend type based on env. For now hardcode to sqlite.
	aragorn.gandalf = gandalf.NewSqliteGandalf()
	return aragorn
}

func (aragorn *Aragorn) Run() {
	r := gin.Default()
	// User APIs.
	r.GET("/api/v1/resources/users/fetch/:user_id", aragorn.getUser)
	r.POST("/api/v1/resources/users/register", aragorn.registerUser)

	// Supplier APIs.
	r.GET("/api/v1/resources/suppliers/fetchAll", aragorn.getAllSuppliers)       // Administrator API. Returns all the suppliers.
	r.POST("/api/v1/resources/suppliers/register", aragorn.registerSupplier)     // Administrator API. // Register Supplier.
	r.GET("/api/v1/resources/suppliers/fetch/:supplier_id", aragorn.getSupplier) // Administrator API. Gets the supplier info.

	// Item APIs.
	r.GET("/api/v1/resources/items/fetch", aragorn.getSupplierItems)      // Administrator API. Gets all items by a supplier.
	r.POST("/api/v1/resources/items/register", aragorn.registerItem)      // Administrator API. Registers item.
	r.POST("/api/v1/resources/items/remove/:item_id", aragorn.removeItem) // Administrator API. Removes item

	// Auction APIs.
	r.GET("/api/v1/resources/auctions/fetchallauctions", aragorn.getAllAuctions)     // Returns all the live auctions.
	r.GET("/api/v1/resources/auctions/fetchmaxbids", aragorn.getMaxBids)             // Returns the max bids for all requested items so far.
	r.POST("/api/v1/resources/auctions/registerbid", aragorn.registerBid)            // Registers a new bid by the user.
	r.GET("/api/v1/resources/auctions/fetchuserauctions", aragorn.fetchUserAuctions) // Fetches all items on which the user had previously bid.

	// Order APIs.
	//r.GET("/api/v1/resources/orders/getUserOrders", aragorn.getUserOrders) // Administrator API.
	//r.GET("/api/v1/resources/orders/getPaymentPendingOrders", aragorn.getPaymentPendingOrders) // Administrator API.
	//r.GET("/api/v1/resources/orders/getDeliveryPendingOrders", aragorn.getDeliveryPendingOrders) // Administrator API.
	//r.GET("/api/v1/resources/orders/getOrder", aragorn.getOrder) // Administrator API.
	//r.POST("/api/v1/resources/orders/updateOrder", aragorn.updateOrder) // Administrator API.

	// Start router.
	r.Run()
}

func (aragorn *Aragorn) authenticate(c *gin.Context) (string, error) {
	idToken := c.Request.Header["Authorization"][0]
	client, err := aragorn.firebaseApp.Auth(context.Background())
	if err != nil {
		return "", err
	}
	token, err := client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return "", err
	}
	return token.UID, nil
}

func (aragorn *Aragorn) doesUserExist(userID string) bool {
	user := aragorn.gandalf.GetUserByID(userID)
	if user.UserID != userID {
		return false
	}
	return true
}

func (aragorn *Aragorn) getUser(c *gin.Context) {
	_, err := aragorn.authenticate(c)
	var response APIResponse
	if err != nil {
		response.Status = http.StatusUnauthorized
		response.ErrorMsg = "Invalid user"
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	var user gandalf.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if user.UserID == "" {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = "User ID is incorrect"
		c.JSON(http.StatusBadRequest, response)
		return
	}
	fullUser := aragorn.gandalf.GetUserByID(user.UserID)
	response.Status = http.StatusOK
	response.ErrorMsg = ""
	response.Data = &fullUser
	c.JSON(http.StatusOK, response)
}

func (aragorn *Aragorn) registerUser(c *gin.Context) {
	userID, err := aragorn.authenticate(c)
	var response APIResponse
	if err != nil {
		response.Status = http.StatusUnauthorized
		response.ErrorMsg = "Invalid user"
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	var user gandalf.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if user.UserID != userID {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = "Invalid user ID"
		c.JSON(http.StatusBadRequest, response)
		return
	}
	err = aragorn.gandalf.RegisterUser(user.UserID, user.Name, user.EmailID, user.PhNum, user.Address)
	if err != nil {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = fmt.Sprintf("Error while registering user: %v", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response.Status = http.StatusOK
	response.ErrorMsg = ""
	c.JSON(http.StatusOK, response)
}

func (aragorn *Aragorn) registerSupplier(c *gin.Context) {

}

func (aragorn *Aragorn) getAllSuppliers(c *gin.Context) {

}

func (aragorn *Aragorn) getSupplier(c *gin.Context) {

}

func (aragorn *Aragorn) getSupplierItems(c *gin.Context) {

}

func (aragorn *Aragorn) registerItem(c *gin.Context) {
	userID, err := aragorn.authenticate(c)
	var response APIResponse
	if err != nil {
		response.Status = http.StatusUnauthorized
		response.ErrorMsg = "Invalid user"
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	var item gandalf.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = "Invalid input JSON"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = aragorn.gandalf.RegisterItem(userID, item.ItemName, item.ItemDescription, item.ItemQty, item.AuctionStartTime, item.MinPrice)
	if err != nil {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = fmt.Sprintf("Error while registering item: %v", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response.Status = http.StatusOK
	response.ErrorMsg = ""
	c.JSON(http.StatusOK, response)
}

func (aragorn *Aragorn) removeItem(c *gin.Context) {
	userID, err := aragorn.authenticate(c)
	var response APIResponse
	if err != nil {
		response.Status = http.StatusUnauthorized
		response.ErrorMsg = "Invalid user"
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	itemID := c.Param("item_id")
	item := aragorn.gandalf.GetItem(itemID)
	if item.UserID != userID {
		response.Status = http.StatusUnauthorized
		response.ErrorMsg = "User does not own this item"
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	err = aragorn.gandalf.DeleteItem(itemID)
	if err != nil {
		response.Status = http.StatusBadRequest
		response.ErrorMsg = fmt.Sprintf("Error while removing item: %v", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response.Status = http.StatusOK
	response.ErrorMsg = ""
	c.JSON(http.StatusOK, response)
}

func (aragorn *Aragorn) getAllAuctions(c *gin.Context) {
	userID, err := aragorn.authenticate(c)
	var response APIResponse
	if err != nil {
		response.Status = http.StatusUnauthorized
		response.ErrorMsg = "Invalid user"
		c.JSON(http.StatusUnauthorized, response)
		return
	}
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
	response.Data = &items
	c.JSON(http.StatusOK, response)
}

func (aragorn *Aragorn) getMaxBids(c *gin.Context) {
	auctions, err := aragorn.gandalf.GetMaxBids(itemIDs)
}

func (aragorn *Aragorn) registerBid(c *gin.Context) {
	err := aragorn.gandalf.RegisterBid(itemID, userID, bidAmount, bidQty)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to register bid for item: %s by user: %s. Backend Error: %v", itemID, userID, err))
	}
	return nil
}

func (aragorn *Aragorn) fetchUserAuctions(c *gin.Context) {
	err := aragorn.gandalf.RegisterBid(itemID, userID, bidAmount, bidQty)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to register bid for item: %s by user: %s. Backend Error: %v", itemID, userID, err))
	}
	return nil
}
