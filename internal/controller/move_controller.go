package controller

import (
	config "battledak-server/configs"
	"battledak-server/internal/dto"
	"battledak-server/internal/service"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type MoveController interface {
	Move(ctx *gin.Context)
}

type moveControllerImpl struct {
	moveService   service.MoveService
	gameService   service.GameService
	cryptoService service.CryptoService
}

func NewMoveController(
	moveService service.MoveService,
	gameService service.GameService,
	cryptoService service.CryptoService,
) *moveControllerImpl {
	return &moveControllerImpl{
		moveService:   moveService,
		gameService:   gameService,
		cryptoService: cryptoService,
	}
}

func (m *moveControllerImpl) Move(ctx *gin.Context) {
	var input dto.MoveRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid move request"})
		return
	}

	game, err := m.gameService.GetGameFromRedis(input.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Game is expired"})
		return
	}

	if !m.gameService.CheckIfIsInTheTimeLimit(game) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Game is expired"})
		return
	}

	if input.Player != "USER" && input.Player != "HOUSE" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player"})
		return
	}

	if input.Player == game.LastMove {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "It's not your turn"})
		return
	}

	cell := -1
	isHit := false
	defaultHits := config.GetTotalHits()

	if input.Player == "USER" {
		game, cell, isHit, err = userMove(input, game)
		hits := m.moveService.CountHits(game.HouseGrid)
		if hits == defaultHits {
			game.Winner = "USER"
			game.FinishedAt = time.Now()
		}
	} else if input.Player == "HOUSE" {
		game, cell, isHit, err = houseMove(m, game)
		hits := m.moveService.CountHits(game.UserGrid)
		if hits == defaultHits {
			game.Winner = "HOUSE"
			game.FinishedAt = time.Now()
		}
	}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	publicGame := m.gameService.MapperGameToPublicGame(game)

	err = m.gameService.SetGameToRedis(game)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	move := &dto.Move{
		PublicGame: publicGame,
		IsHit:      isHit,
		Cell:       cell,
	}

	ctx.JSON(http.StatusOK, move)
}

func userMove(input dto.MoveRequest, game dto.Game) (dto.Game, int, bool, error) {
	gridSize := config.GetGridSize()

	// Validate the input cell
	if input.Player == "USER" && !validateInputCell(input.Cell, gridSize) {
		return game, -1, false, fmt.Errorf("invalid cell")
	}

	// Check if the cell is valid
	if game.HouseGrid[input.Cell] == dto.Hit || game.HouseGrid[input.Cell] == dto.Miss {
		return game, -1, false, fmt.Errorf("this cell has already been chosen")
	}

	// Check if the cell is a hit or miss
	isHit := false
	if game.DecryptedHouseGrid[input.Cell] == 0 {
		isHit = true
	}

	// Update game state
	game.LastMove = "USER"
	game.UpdatedAt = time.Now()
	if isHit {
		game.HouseGrid[input.Cell] = dto.Hit
	} else {
		game.HouseGrid[input.Cell] = dto.Miss
	}

	return game, input.Cell, isHit, nil
}

func houseMove(m *moveControllerImpl, game dto.Game) (dto.Game, int, bool, error) {
	// Generate a cell for the house move
	err, cell := m.moveService.GenerateHouseMove(game.UserGrid)
	if err != nil {
		return game, -1, false, err
	}

	// Check if the cell is valid
	if game.UserGrid[cell] == dto.Hit || game.UserGrid[cell] == dto.Miss {
		return game, -1, false, fmt.Errorf("this cell has already been chosen")
	}

	// Check if the cell is a hit or miss
	err, isHit := m.cryptoService.IsZero(game.UserGrid[cell])
	if err != nil {
		return game, -1, false, err
	}

	// Update game state
	game.LastMove = "HOUSE"
	game.UpdatedAt = time.Now()
	if isHit {
		game.UserGrid[cell] = dto.Hit
	} else {
		game.UserGrid[cell] = dto.Miss
	}

	return game, cell, isHit, err
}

func validateInputCell(cell int, gridSize int) bool {
	totalCells := gridSize * gridSize
	if cell < 0 || cell >= totalCells {
		return false
	}
	return true
}
