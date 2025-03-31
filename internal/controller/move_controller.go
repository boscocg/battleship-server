package controller

import (
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

	err, game := m.gameService.GetGameFromRedis(input.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
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

	if input.Player == "USER" {
		err, game, cell, isHit = userMove(input, game)
		hits := m.moveService.CountHits(game.HouseGrid)
		if hits == 20 {
			game.Winner = "USER"
		}
	} else if input.Player == "HOUSE" {
		err, game, cell, isHit = houseMove(m, game)
		hits := m.moveService.CountHits(game.UserGrid)
		if hits == 20 {
			game.Winner = "HOUSE"
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

func userMove(input dto.MoveRequest, game dto.Game) (error, dto.Game, int, bool) {
	// Validate the input cell
	if input.Player == "USER" && validateInputCell(input.Cell) == false {
		return fmt.Errorf("Invalid cell"), game, -1, false
	}

	// Check if the cell is valid
	if game.HouseGrid[input.Cell] == dto.Hit || game.HouseGrid[input.Cell] == dto.Miss {
		return fmt.Errorf("This cell has already been chosen"), game, -1, false
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

	return nil, game, input.Cell, isHit
}

func houseMove(m *moveControllerImpl, game dto.Game) (error, dto.Game, int, bool) {
	// Generate a cell for the house move
	err, cell := m.moveService.GenerateHouseMove(game.UserGrid)
	if err != nil {
		return err, game, -1, false
	}

	// Check if the cell is valid
	if game.UserGrid[cell] == dto.Hit || game.UserGrid[cell] == dto.Miss {
		return fmt.Errorf("This cell has already been chosen"), game, -1, false
	}

	// Check if the cell is a hit or miss
	err, isHit := m.cryptoService.IsZero(game.UserGrid[cell])
	if err != nil {
		return err, game, -1, false
	}

	// Update game state
	game.LastMove = "HOUSE"
	game.UpdatedAt = time.Now()
	if isHit {
		game.UserGrid[cell] = dto.Hit
	} else {
		game.UserGrid[cell] = dto.Miss
	}

	return err, game, cell, isHit
}

func validateInputCell(cell int) bool {
	if cell < 0 || cell >= 100 {
		return false
	}
	return true
}
