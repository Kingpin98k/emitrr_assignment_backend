package main

import (
	handlers "example/Card-Game-Backend/handlers"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main(){

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
}

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
	router.GET("/currentGame/:email",handlers.GetCurrentGame,handlers.StartGame,handlers.SaveGame)
	router.POST("/startGame/:email",handlers.StartGame,handlers.SaveGame)
	router.POST("/drawCard/:email",handlers.MovePointer,handlers.SaveGame)
	router.GET("/leaderboard",handlers.GetLeaderboard)

	router.Run(":8080")
}