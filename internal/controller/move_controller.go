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
	moveService service.MoveService
	gameService service.GameService
}

func NewMoveController(
	moveService service.MoveService,
	gameService service.GameService,
) *moveControllerImpl {
	return &moveControllerImpl{
		moveService: moveService,
		gameService: gameService,
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

	if input.Player == "USER" {
		err, game = userMove(input, game)
	} else if input.Player == "HOUSE" {
		err, game, cell = houseMove(m, game)
	}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
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
		Succeded:   true,
	}

	if cell != -1 {
		move.Cell = cell
	}

	ctx.JSON(http.StatusOK, move)
}

func validateInputCell(cell int) bool {
	if cell < 0 || cell >= 100 {
		return false
	}
	return true
}

func userMove(input dto.MoveRequest, game dto.Game) (error, dto.Game) {
	// Update game state
	game.LastMove = "USER"
	game.UpdatedAt = time.Now()
	game.HouseGrid[input.Cell] = "HIT"

	if input.Player == "USER" && validateInputCell(input.Cell) == false {
		return fmt.Errorf("Invalid cell"), game
	}

	return nil, game
}

func houseMove(m *moveControllerImpl, game dto.Game) (error, dto.Game, int) {
	err, cell := m.moveService.GenerateHouseMove()
	if err != nil {
		return err, game, -1
	}

	// Update game state
	game.LastMove = "HOUSE"
	game.UpdatedAt = time.Now()
	game.UserGrid[cell] = "HIT"

	return err, game, cell
}
