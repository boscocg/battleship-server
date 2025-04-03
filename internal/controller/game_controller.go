package controller

import (
	"battledak-server/internal/dto"
	"battledak-server/internal/service"
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
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Game ID is required"})
		return
	}

	game, err := u.gameService.GetGameFromRedis(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	publicGame := u.gameService.MapperGameToPublicGame(game)

	ctx.JSON(http.StatusOK, publicGame)
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
		LastMove:           "HOUSE",
		ID:                 strconv.FormatUint(id, 10),
		UserGrid:           input.UserGrid,
		HouseGrid:          houseGrid,
		UpdatedAt:          time.Now(),
		CreatedAt:          time.Now(),
		DecryptedHouseGrid: decryptedHouseGrid,
	}

	// Check if the game exists in Redis
	val, _ := config.AppConfig.RedisClient.Get(config.AppConfig.Ctx, game.ID).Bytes()
	if val != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Game already exists"})
		return
	}

	err := u.gameService.SetGameToRedis(*game)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	publicGame := u.gameService.MapperGameToPublicGame(*game)

	ctx.JSON(http.StatusOK, publicGame)
}
