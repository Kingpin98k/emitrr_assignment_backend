package main

import (
	"log"
	"net/http"
	"os"

	handlers "example/Card-Game-Backend/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
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
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the Card Game API",
		})
	})
	router.GET("/currentGame/:email", handlers.GetCurrentGame, handlers.StartGame, handlers.SaveGame)
	router.POST("/startGame/:email", handlers.StartGame, handlers.SaveGame)
	router.POST("/drawCard/:email", handlers.MovePointer, handlers.SaveGame)
	router.GET("/leaderboard", handlers.GetLeaderboard)

	// Get the host and port from environment variables or use default values
	host := envOr("0.0.0.0")
	port := envOr("8080")

	server := &http.Server{
		Addr:    host + ":" + port, // Construct host and port
		Handler: router,
	}

	log.Printf("Server started at http://%s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v", err)
	}
}

// envOr returns the value of the environment variable with the given name
// if it exists, otherwise it returns the provided default value.
func envOr(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return key // Return default value if environment variable is not set
	}
	return val
}
