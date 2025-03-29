package service

import (
	"battledak-server/internal/dto"
	"fmt"
	"math/rand"
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
			// Check left
			if hit%10 > 0 && userGrid[hit-1] != dto.Hit && userGrid[hit-1] != dto.Miss {
				return nil, hit - 1
			}
			// Check right
			if hit%10 < 9 && userGrid[hit+1] != dto.Hit && userGrid[hit+1] != dto.Miss {
				return nil, hit + 1
			}
			// Check up
			if hit-10 >= 0 && userGrid[hit-10] != dto.Hit && userGrid[hit-10] != dto.Miss {
				return nil, hit - 10
			}
			// Check down
			if hit+10 < len(userGrid) && userGrid[hit+10] != dto.Hit && userGrid[hit+10] != dto.Miss {
				return nil, hit + 10
			}
		}
	}

	if len(availableCells) > 0 {
		randomIndex := rand.Intn(len(availableCells))
		return nil, availableCells[randomIndex]
	}

	return fmt.Errorf("No more moves available"), -1
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
