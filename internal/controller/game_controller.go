package controller

import (
	"battledak-server/internal/dto"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	config "battledak-server/configs"
)

type GameController interface {
	StartGame(ctx *gin.Context)
	GetGame(ctx *gin.Context)
}

type gameControllerImpl struct {
}

func NewGameController() *gameControllerImpl {
	return &gameControllerImpl{}
}

func (u *gameControllerImpl) GetGame(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	game := &dto.Game{
		ID: id,
	}

	// Check if the game exists in Redis
	err = config.AppConfig.RedisClient.Get(config.AppConfig.Ctx, strconv.FormatUint(id, 10)).Err()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	// Retrieve the game data from Redis
	val, err := config.AppConfig.RedisClient.Get(config.AppConfig.Ctx, strconv.FormatUint(id, 10)).Bytes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve game"})
		return
	}

	if err = json.Unmarshal(val, game); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse game data"})
		return
	}

	ctx.JSON(http.StatusOK, game)
}

func (u *gameControllerImpl) StartGame(ctx *gin.Context) {
	var input dto.GameRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timestamp := uint64(time.Now().UnixNano())
	randomPart := uint64(rand.Intn(1000))
	id := timestamp + randomPart
	game := &dto.PublicGameResponse{
		LastMove: "USER",
		ID:       id,
	}

	// Check if the game exists in Redis
	val, err := config.AppConfig.RedisClient.Get(config.AppConfig.Ctx, strconv.FormatUint(game.ID, 10)).Bytes()
	if val != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Game already exists"})
		return
	}

	gameJSON, err := json.Marshal(game)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process game data"})
		return
	}

	// Set the game in Redis with a 1-hour expiration time
	err = config.AppConfig.RedisClient.Set(config.AppConfig.Ctx, strconv.FormatUint(game.ID, 10), gameJSON, 0).Err()
	if err != nil {
		log.Printf("Error setting game in Redis: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start game"})
		return
	}

	ctx.JSON(http.StatusOK, game)
}
