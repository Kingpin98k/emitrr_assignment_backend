package main

import (
	game_handlers "example/Card-Game-Backend/handlers"

	"github.com/gin-gonic/gin"
)

func main(){
	router := gin.Default()

	
	// Enable CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	router.GET("/", func(c *gin.Context){
		c.JSON(200, gin.H{
			"message": "Welcome to the Card Game API",
		})
	})
	router.GET("/currentGame/:email",game_handlers.GetCurrentGame)
	router.POST("/startGame/:email",game_handlers.StartGame,game_handlers.SaveGame)
	router.POST("/drawCard/:email",game_handlers.MovePointer,game_handlers.SaveGame)


	router.Run("localhost:8080")
}