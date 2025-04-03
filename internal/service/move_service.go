package service

import (
	"battledak-server/internal/dto"
	"fmt"
	"math/rand"

	config "battledak-server/configs"
)

type MoveService interface {
	GenerateHouseMove(userGrid []dto.CellType) (error, int)
	CountHits(grid []dto.CellType) int
}

type moveServiceImpl struct {
}

func NewMoveService() *moveServiceImpl {
	return &moveServiceImpl{}
}

func (g *moveServiceImpl) GenerateHouseMove(userGrid []dto.CellType) (error, int) {
	gridSize := config.GetGridSize()
	hits := make([]int, 0)
	availableCells := make([]int, 0)
	for i, cell := range userGrid {
		if cell == dto.Hit {
			hits = append(hits, i)
		}
		if cell != dto.Hit && cell != dto.Miss {
			availableCells = append(availableCells, i)
		}
	}

	if len(hits) > 0 {
		for _, hit := range hits {
			// Get row and column from hit index
			row := hit / gridSize
			col := hit % gridSize

			// Check left
			if col > 0 && userGrid[hit-1] != dto.Hit && userGrid[hit-1] != dto.Miss {
				return nil, hit - 1
			}
			// Check right
			if col < gridSize-1 && userGrid[hit+1] != dto.Hit && userGrid[hit+1] != dto.Miss {
				return nil, hit + 1
			}
			// Check up
			if row > 0 && userGrid[hit-gridSize] != dto.Hit && userGrid[hit-gridSize] != dto.Miss {
				return nil, hit - gridSize
			}
			// Check down
			if row < gridSize-1 && userGrid[hit+gridSize] != dto.Hit && userGrid[hit+gridSize] != dto.Miss {
				return nil, hit + gridSize
			}
		}
	}

	if len(availableCells) > 0 {
		randomIndex := rand.Intn(len(availableCells))
		return nil, availableCells[randomIndex]
	}

	return fmt.Errorf("no more moves available"), -1
}

func (g *moveServiceImpl) CountHits(grid []dto.CellType) int {
	count := 0
	for _, cell := range grid {
		if cell == dto.Hit {
			count++
		}
	}
	return count
}
