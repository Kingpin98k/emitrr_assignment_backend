package game_handlers

import (
	"context"
	"encoding/json"
	game_models "example/Card-Game-Backend/models"
	utils "example/Card-Game-Backend/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var redisClient = utils.Client()

func GetCurrentGame(c*gin.Context){
	email := c.Param("email")
	// game := Game{EMAIL: email, DECK: "deck"}
	// c.JSON(http.StatusOK, game)

	val, err := redisClient.Get(ctx,email).Result()
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// var game Game
	var game game_models.Game

	err = json.Unmarshal([]byte(val), &game)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	leaderboard, err := redisClient.ZRevRangeWithScores(ctx,"leaderboard",0,2).Result()

	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"game": game, "leaderboard": leaderboard})
}

func StartGame(c*gin.Context){
	email := c.Param("email")
	var highScore int = 0

  val, err := redisClient.Get(ctx,email).Result()

	if(err != nil){
		if(err.Error() != "redis: nil"){
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if val != "" {
		var game game_models.Game
		err = json.Unmarshal([]byte(val), &game)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		highScore = game.CURRENT_HIGH_SCORE
	}

	// game := Game{EMAIL: email, DECK: "deck"}
	deck := utils.GenerateNewDeck()

	game := game_models.Game{EMAIL: email, DECK: deck, POINTER: 0, CURRENT_HIGH_SCORE: highScore}

	c.Set("game", game)

	c.Next()
}

func MovePointer(c*gin.Context){
	email := c.Param("email")

	val, err := redisClient.Get(ctx,email).Result()

	if(err != nil){
		if(err.Error() != "redis: nil"){
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if val != "" {
		var game game_models.Game
		err = json.Unmarshal([]byte(val), &game)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		game.POINTER++
		if(game.POINTER >= 5){
			game.POINTER = 0
			game.CURRENT_HIGH_SCORE++
		}
		c.Set("game", game)
		c.Next()	
	}else{
		c.JSON(http.StatusBadRequest, gin.H{"error": "game not found"})
		return
	}
}


func SaveGame(c*gin.Context){
	var game game_models.Game

	value, exists := c.Get("game")
	if !exists {
		if err:=c.BindJSON(&game); err!=nil{
			fmt.Println("No Body")
			c.JSON(http.StatusBadRequest, gin.H{"error": "game not found"})
		return
		}
	}else{
		game = value.(game_models.Game)
	}


	obj,err := json.Marshal(game)
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = redisClient.Set(ctx,game.EMAIL, obj, 0).Err()
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	z := redis.Z{Score: float64(game.CURRENT_HIGH_SCORE), Member: game.EMAIL}
	err = redisClient.ZAdd(ctx,"leaderboard",&z).Err()
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, game)
}