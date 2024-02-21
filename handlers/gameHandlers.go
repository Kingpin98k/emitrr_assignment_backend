package handlers

import (
	"context"
	"encoding/json"
	game_models "example/Card-Game-Backend/models"
	utils "example/Card-Game-Backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var redisClient = utils.Client()

func GetCurrentGame(c *gin.Context) {
	email := c.Param("email")

	val, err := redisClient.Get(ctx, email).Result()

	if err != nil {
			if err.Error() == "redis: nil" {
					// Key not found in Redis, continue with other handlers
					c.Next()
					return
			} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
			}
	}

	var game game_models.Game
	err = json.Unmarshal([]byte(val), &game)
	if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}

	var gameResponse game_models.GameResponse	

	gameResponse.EMAIL = game.EMAIL
	gameResponse.POINTER = game.POINTER
	gameResponse.CURRENT_HIGH_SCORE = game.CURRENT_HIGH_SCORE
	if(game.POINTER == -1 || game.POINTER == -2){
		gameResponse.TOP_CARD = "start"
		gameResponse.DEFFUSECARDS = 0
		gameResponse.BOMBCARDS = 0
	}else{
		gameResponse.TOP_CARD = strings.Split(game.DECK,",")[game.POINTER]
	}
	gameResponse.BOMBCARDS = game.BOMBCARDS
	gameResponse.DEFFUSECARDS = game.DEFFUSECARDS


	c.JSON(http.StatusOK, gameResponse)
	c.Abort()
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

	game := game_models.Game{EMAIL: email, DECK: deck, POINTER: -1, CURRENT_HIGH_SCORE: highScore,BOMBCARDS: 0,DEFFUSECARDS: 0}

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
		if(game.POINTER != -2){

			game.POINTER++
			if(game.POINTER >= 5){
				game.POINTER = -2
				game.BOMBCARDS = 0
				game.DEFFUSECARDS = 0
				game.CURRENT_HIGH_SCORE++
				}else{
					if(strings.Split(game.DECK,",")[game.POINTER] == "bomb"){
						game.BOMBCARDS++
						}else if(strings.Split(game.DECK,",")[game.POINTER] == "defuse"){
							game.DEFFUSECARDS++
						}
					}
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "game not found"})
		return
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

	var gameResponse game_models.GameResponse	

	gameResponse.EMAIL = game.EMAIL
	gameResponse.POINTER = game.POINTER
	gameResponse.CURRENT_HIGH_SCORE = game.CURRENT_HIGH_SCORE
	if(game.POINTER == -1 || game.POINTER == -2){
		gameResponse.TOP_CARD = "start"
		gameResponse.DEFFUSECARDS = 0
		gameResponse.BOMBCARDS = 0
	}else{
		gameResponse.TOP_CARD = strings.Split(game.DECK,",")[game.POINTER]
	}
	gameResponse.BOMBCARDS = game.BOMBCARDS
	gameResponse.DEFFUSECARDS = game.DEFFUSECARDS
	
	c.JSON(http.StatusOK, gameResponse)
}