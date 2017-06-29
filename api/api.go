package api

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

func Run() {
	router := gin.Default()

	api := router.Group("multicard/api/v1")
	{
		users := api.Group("/users")
		{
			users.POST("/", createUser)
			// users.GET("/:userName", getUser)
			wallets := users.Group("/:userId/wallets")
			{
				wallets.POST("/", createUser)
				wallets.GET("/", createUser)
				wallets.GET("/:walletId", createUser)
				wallets.PUT("/:walletId", createUser)
				wallets.DELETE("/:walletId", createUser)
				wallets.POST("/:walletId/purchase", createUser)

				cards := wallets.Group("/:walletId/cards")
				{
					cards.POST("/", createUser)
					cards.GET("/", createUser)
					cards.GET("/:cardId", createUser)
					cards.PUT("/:cardId", createUser)
					cards.DELETE("/:cardId", createUser)
				}
			}
			// wallets := users.Group("/:userId/wallets")
			// {
			// 	wallets.POST("/", createWallet)
			// 	wallets.GET("/", fetchAllWallets)
			// 	wallets.GET("/:walletId", getWallet)
			// 	wallets.PUT("/:walletId", updateWallet)
			// 	wallets.DELETE("/:walletId", removeWallet)
			// 	wallets.POST("/:walletId/purchase", purchaseOnWallet)

			// 	cards := wallets.Group("/:walletId/cards")
			// 	{
			// 		cards.POST("/", createCard)
			// 		cards.GET("/", fetchAllCards)
			// 		cards.GET("/:cardId", getCard)
			// 		cards.PUT("/:cardId", updateCard)
			// 		cards.DELETE("/:cardId", removeCard)
			// 	}
			// }
		}
		
		// wallet := router.Group("/user/:userId/wallet")
		// {
		// 	wallet.POST("/", createWallet)
		// 	wallet.GET("/", fetchAllWallets)
		// 	wallet.GET("/:walletId", getWallet)
		// 	wallet.PUT("/:walletId", updateWallet)
		// 	wallet.DELETE("/:walletId", removeWallet)
		// }
		// wallet.POST("/", createCard)
		// wallet.GET("/", fetchAllWallets)
		// wallet.GET("/:id", getWallet)
		// wallet.PUT("/:id", updateWallet)
		// wallet.DELETE("/:id", removeWallet)
	}

	// wallet := router.Group("/wallet")
	// {
	// 	wallet.POST("/", createWallet)
	// 	wallet.GET("/", fetchAllWallets)
	// 	wallet.GET("/:id", getWallet)
	// 	wallet.PUT("/:id", updateWallet)
	// 	wallet.DELETE("/:id", removeWallet)
	// }

	// card := router.Group("/cards")
	// {
	// 	card.POST("/", createCard)
	// 	card.GET("/", fetchAllCards)
	// 	card.GET("/:id", getCard)
	// 	card.PUT("/:id", updateCard)
	// 	card.DELETE("/:id", removeCard)
	// }

	router.Run()
}

func createUser(c *gin.Context) {
    //    completed, _ := strconv.Atoi(c.PostForm("completed"))
    //    todo := Todo{Title: c.PostForm("title"), Completed: completed};
    //    db := Database()
    //    db.Save(&todo)
    //    c.JSON(http.StatusCreated, gin.H{"status" : http.StatusCreated, "message" : "User created successfully!", "userId": user.ID})
	c.JSON(http.StatusCreated, gin.H{"status" : http.StatusCreated, "message" : "User created successfully!"})
}

func getUser(c *gin.Context) {
    //    completed, _ := strconv.Atoi(c.PostForm("completed"))
    //    todo := Todo{Title: c.PostForm("title"), Completed: completed};
    //    db := Database()
    //    db.Save(&todo)
    //    c.JSON(http.StatusCreated, gin.H{"status" : http.StatusCreated, "message" : "User created successfully!", "userId": user.ID})
	user := struct {
        UserId uint64
    } { 1 }
	c.JSON(http.StatusOK, gin.H{"status" : http.StatusOK, "data" : user})
}