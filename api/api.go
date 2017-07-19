package api

import (
	// "errors"
	"net/http"
	"strconv"

	"multicard/models"

	"github.com/gin-gonic/gin"
)

func Run() {
	router := gin.Default()

	api := router.Group("multicard/api/v1")
	{
		users := api.Group("/users")
		{
			users.POST("/", createUser)
			users.GET("/:userName", getUserByLogin)
			usersWallets := users.Group("/:userName/wallets") 
			{
				usersWallets.POST("/", createWallet)
				usersWallets.GET("/", fetchAllWallets)
			}			
		}
		wallets := api.Group("/wallets")
		{
			wallets.GET("/:walletId", getWallet)
			wallets.PUT("/:walletId", updateWallet)
			wallets.DELETE("/:walletId", deleteWallet)
			wallets.POST("/:walletId/purchase", purchaseOnWallet)

			walletsCards := wallets.Group("/:walletId/cards")
			{
				walletsCards.POST("/", createCard)
				walletsCards.GET("/", fetchAllCards)
			}
		}
		cards := api.Group("/cards")
		{
			cards.GET("/:cardId", getCard)
			cards.PUT("/:cardId", updateCard)
			cards.DELETE("/:cardId", deleteCard)
		}
	}

	router.Run()
}

// User
// func userIdMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		user := &models.User{ Login: c.Param("userId") }
// 		user.GetByLogin()

// 		if user.ID > 0 {
// 			c.Set("userId", user.ID)
// 		} else {
// 			c.AbortWithError(http.StatusBadRequest, errors.New("User not found!"))
// 		}
// 	}
// }
func createUser(c *gin.Context) {
	var user *models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Unable to read user json: " + err.Error()})
		return
	}

	if err := user.Create(); err != nil {
		c.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict, "message": "Unable to create user: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "User created successfully!", "resourceId": user.ID})
}
func getUserByLogin(c *gin.Context) {
	user := &models.User{ Login: c.Param("userName") }
	if !user.GetByLogin() {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "User not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": user})
}

// Wallet
func createWallet(c *gin.Context) {
	user := &models.User{ Login: c.Param("userName") }
	if !user.GetByLogin() {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "User not found!"})
		return
	}
	// userId, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid user id: " + err.Error()})
	// 	return
	// }

	var wallet *models.Wallet
	if err := c.BindJSON(&wallet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Unable to read json: " + err.Error()})
		return
	}
	wallet.UserId = user.ID
	if err := wallet.Create(); err != nil {
		c.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict, "message": "Unable to create wallet: " + err.Error()})
		return
	}	
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Wallet created successfully!", "resourceId": wallet.ID})
}
func fetchAllWallets(c *gin.Context) {
	user := &models.User{ Login: c.Param("userName") }
	if !user.GetByLogin() {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "User not found!"})
		return
	}
	// userId, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid user id: " + err.Error()})
	// 	return
	// }

	w := &models.Wallet{ UserId: user.ID }
	wallets := w.FechtAllFromUser()

	if len(wallets) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No wallet found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": wallets})
}
func getWallet(c *gin.Context) {
	walletId, err := strconv.ParseUint(c.Param("walletId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid wallet id: " + err.Error()})
		return
	}

	wallet := &models.Wallet{ ID: walletId }
	if !wallet.Get() {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Wallet not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": wallet})
}
func updateWallet(c *gin.Context) {
	walletId, err := strconv.ParseUint(c.Param("walletId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid wallet id: " + err.Error()})
		return
	}

	var wallet *models.Wallet
	if err := c.BindJSON(&wallet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Unable to read json: " + err.Error()})
		return
	}
	wallet.ID = walletId

	if err := wallet.Update(); err != nil {
		c.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict, "message": "Unable to update wallet: " + err.Error()})
		return
	}	
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Wallet updated successfully!"})
}
func deleteWallet(c *gin.Context) {
	walletId, err := strconv.ParseUint(c.Param("walletId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid wallet id: " + err.Error()})
		return
	}

	wallet := &models.Wallet{ ID: walletId }
	if err := wallet.Delete(); err != nil {
		c.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict, "message": "Unable to delete wallet: " + err.Error()})
		return
	}	
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Wallet deleted successfully!"})
}
func purchaseOnWallet(c *gin.Context) {
	walletId, err := strconv.ParseUint(c.Param("walletId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid wallet id: " + err.Error()})
		return
	}	

	wallet := &models.Wallet{ ID: walletId }
	if !wallet.Get() {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Wallet not found!"})
		return
	}

	purchase := struct {
		Value float64 `json:"value"`
	} {}
	if err := c.BindJSON(&purchase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Unable to read json: " + err.Error()})
		return
	}
	if wallet.AvailableCredit < purchase.Value {
		c.JSON(http.StatusForbidden, gin.H{"status": http.StatusForbidden, "message": "Not enough credit to make the purchase"})
		return
	}

	if err := wallet.Purchase(purchase.Value); err != nil {
		c.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict, "message": "Unable to purchase on wallet: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status" : http.StatusOK, "message" : "Purchase was successfull!"})
}

// Card
func createCard(c *gin.Context) {
	walletId, err := strconv.ParseUint(c.Param("walletId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid wallet id: " + err.Error()})
		return
	}

	var card *models.Card
	if err := c.BindJSON(&card); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Unable to read card json: " + err.Error()})
		return
	}
	card.WalletId = walletId

	if err := card.Create(); err != nil {
		c.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict, "message": "Unable to create card: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Card created successfully!", "resourceId": card.ID})
}
func fetchAllCards(c *gin.Context) {
	walletId, err := strconv.ParseUint(c.Param("walletId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid wallet id: " + err.Error()})
		return
	}

	card := &models.Card{ WalletId: walletId }
	cards := card.FechtAllFromWallet()

	if len(cards) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No card found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": cards})
}
func getCard(c *gin.Context) {
	cardId, err := strconv.ParseUint(c.Param("cardId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid card id: " + err.Error()})
		return
	}

	card := &models.Card{ ID: cardId }
	if !card.Get() {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Card not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": card})
}
func updateCard(c *gin.Context) {
	cardId, err := strconv.ParseUint(c.Param("cardId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid card id: " + err.Error()})
		return
	}

	var card *models.Card
	if err := c.BindJSON(&card); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Unable to read json: " + err.Error()})
		return
	}
	card.ID = cardId

	if err := card.Update(); err != nil {
		c.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict, "message": "Unable to update card: " + err.Error()})
		return
	}	
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Card updated successfully!"})
}
func deleteCard(c *gin.Context) {
	cardId, err := strconv.ParseUint(c.Param("cardId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid card id: " + err.Error()})
		return
	}

	card := &models.Card{ ID: cardId }
	if err := card.Delete(); err != nil {
		c.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict, "message": "Unable to delete card: " + err.Error()})
		return
	}	
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Card deleted successfully!"})
}