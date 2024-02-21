package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLeaderboard(c *gin.Context){
    leaderboard, err := redisClient.ZRevRangeWithScores(ctx,"leaderboard",0,-1).Result()

    if err!=nil{
        fmt.Println(err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK,leaderboard)
    
}