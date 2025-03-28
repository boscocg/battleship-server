package controller

import (
	"battledak-server/internal/dto"
	"battledak-server/internal/service"
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
	gameService service.GameService
}

func NewGameController(
	gameService service.GameService,
) *gameControllerImpl {
	return &gameControllerImpl{
		gameService: gameService,
	}
}

func (u *gameControllerImpl) GetGame(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	game := &dto.PublicGameResponse{
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grid format"})
		return
	}

	houseGrid, decryptedHouseGrid := u.gameService.GenerateHouseGrid()

	if len(input.UserGrid) != 100 || len(houseGrid) != 100 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User grid must be 100 cells"})
		return
	}

	timestamp := uint64(time.Now().UnixNano())
	randomPart := uint64(rand.Intn(1000))
	id := timestamp + randomPart
	game := &dto.Game{
		LastMove:           "USER",
		ID:                 id,
		UserGrid:           input.UserGrid,
		HouseGrid:          houseGrid,
		UpdatedAt:          time.Now(),
		DecryptedHouseGrid: decryptedHouseGrid,
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

	publicGame := &dto.PublicGameResponse{
		ID:        game.ID,
		LastMove:  game.LastMove,
		UserGrid:  game.UserGrid,
		HouseGrid: game.HouseGrid,
		UpdatedAt: game.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, publicGame)
}
