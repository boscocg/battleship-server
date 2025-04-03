package service

import (
	"battledak-server/internal/dto"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	config "battledak-server/configs"
)

type GameService interface {
	GenerateHouseGrid(gridSize int) ([]dto.CellType, []int)
	GetGameFromRedis(id string) (dto.Game, error)
	SetGameToRedis(game dto.Game) error
	MapperGameToPublicGame(game dto.Game) dto.PublicGame
	CheckIfIsInTheTimeLimit(game dto.Game) bool
}

type gameServiceImpl struct {
}

func NewGameService() *gameServiceImpl {
	return &gameServiceImpl{}
}

func (g *gameServiceImpl) GenerateHouseGrid(gridSize int) ([]dto.CellType, []int) {
	// Seed the random number generator
	digitToHash := map[int]dto.CellType{
		0: "5feceb66ffc86f38d952786c6d696c79c2dbc239dd4e91b46729d73a27fb57e9",
		1: "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b",
		2: "d4735e3a265e16eee03f59718b9b5d03019c07d8b6c51f90da3a666eec13ab35",
		3: "4e1cde028b7ef7a85e01634c9b9f12374d02fc3d9d0c5fd84af0c8ecc04c2864",
		4: "4b227777d4dd1fc61c6f884f48641d02b4d121d3fd328cb08b5531fcacdabf8a",
		5: "ef2d127de37b942baad06145e54b0c619a1f22327b2ebbcfbec78f5564afe39d",
		6: "e7f6c011776e8db7cd330b54174fd76f7d0216b612387a5ffcfb81e6f0919683",
		7: "7d1e0a93f246dbcee814aed2f70e47f5b3360a8ddb4f5422f4dda5bf83744dfc",
		8: "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92",
		9: "21049d1e9d599bf59ef8364908a6938e6ad9c587c2c2c3065b4ac29c558659ca",
	}

	totalCells := gridSize * gridSize
	houseGrid := make([]dto.CellType, totalCells)
	publicGrid := make([]int, totalCells)

	// Initialize all cells to empty (water)
	for i := range houseGrid {
		randomDigit := rand.Intn(9) + 1
		houseGrid[i] = digitToHash[randomDigit]
		publicGrid[i] = randomDigit
	}

	// Place ships
	placeShips(houseGrid, publicGrid, 4, 1, digitToHash, gridSize)
	placeShips(houseGrid, publicGrid, 3, 2, digitToHash, gridSize)
	placeShips(houseGrid, publicGrid, 2, 3, digitToHash, gridSize)
	placeShips(houseGrid, publicGrid, 1, 4, digitToHash, gridSize)

	return houseGrid, publicGrid
}

func placeShips(houseGrid []dto.CellType, publicGrid []int, count, size int, digitToHash map[int]dto.CellType, gridSize int) {
	for range count {
		placed := false
		for !placed {
			// Randomly decide orientation: 0 for horizontal, 1 for vertical
			orientation := rand.Intn(2)
			var startPos int

			if orientation == 0 { // Horizontal
				// Ensure the ship doesn't go off the right edge
				row := rand.Intn(gridSize)
				col := rand.Intn(gridSize - size + 1)
				startPos = row*gridSize + col

				// Check if positions are already occupied or adjacent to other ships
				canPlace := true
				// Check the ship positions and their surroundings
				for j := range size {
					pos := startPos + j

					// Check the ship position itself
					if publicGrid[pos] == 0 {
						canPlace = false
						break
					}

					// Define surrounding positions to check (apenas horizontais e verticais, não diagonais)
					surroundingOffsets := []int{
						-gridSize, // top
						-1, 1,     // left, right
						gridSize, // bottom
					}

					for _, offset := range surroundingOffsets {
						adjPos := pos + offset

						// Make sure we don't go out of bounds and check if adjacent cell is a ship
						// Verifica apenas ortogonalmente (sem diagonal): cada offset já está verificando apenas um lado
						if adjPos >= 0 && adjPos < gridSize*gridSize && // Within grid bounds
							publicGrid[adjPos] == 0 { // It's a ship
							canPlace = false
							break
						}
					}

					if !canPlace {
						break
					}
				}

				if canPlace {
					for j := range size {
						pos := startPos + j
						houseGrid[pos] = digitToHash[0] // Set to ship (represented by 0)
						publicGrid[pos] = 0
					}
					placed = true
				}
			} else { // Vertical
				// Ensure the ship doesn't go off the bottom edge
				row := rand.Intn(gridSize - size + 1)
				col := rand.Intn(gridSize)
				startPos = row*gridSize + col

				// Check if positions are already occupied or adjacent to other ships
				canPlace := true
				// Check the ship positions and their surroundings
				for j := range size {
					pos := startPos + j*gridSize

					// Check the ship position itself
					if publicGrid[pos] == 0 {
						canPlace = false
						break
					}

					// Define surrounding positions to check (apenas horizontais e verticais, não diagonais)
					surroundingOffsets := []int{
						-gridSize, // top
						-1, 1,     // left, right
						gridSize, // bottom
					}

					for _, offset := range surroundingOffsets {
						adjPos := pos + offset

						// Make sure we don't go out of bounds and check if adjacent cell is a ship
						// Verifica apenas ortogonalmente (sem diagonal): cada offset já está verificando apenas um lado
						if adjPos >= 0 && adjPos < gridSize*gridSize && // Within grid bounds
							publicGrid[adjPos] == 0 { // It's a ship
							canPlace = false
							break
						}
					}

					if !canPlace {
						break
					}
				}

				if canPlace {
					for j := range size {
						pos := startPos + j*gridSize
						houseGrid[pos] = digitToHash[0] // Set to ship (represented by 0)
						publicGrid[pos] = 0
					}
					placed = true
				}
			}
		}
	}
}

func (g *gameServiceImpl) GetGameFromRedis(id string) (dto.Game, error) {
	game := &dto.Game{
		ID: id,
	}

	// Check if the game exists in Redis
	err := config.AppConfig.RedisClient.Get(config.AppConfig.Ctx, game.ID).Err()
	if err != nil {
		return dto.Game{}, fmt.Errorf("game not found: %v", err)
	}

	// Retrieve the game data from Redis
	val, err := config.AppConfig.RedisClient.Get(config.AppConfig.Ctx, game.ID).Bytes()
	if err != nil {
		return dto.Game{}, fmt.Errorf("gailed to retrieve game: %v", err)
	}

	if err = json.Unmarshal(val, game); err != nil {
		return dto.Game{}, fmt.Errorf("failed to parse game data: %v", err)
	}

	return *game, nil
}

func (g *gameServiceImpl) SetGameToRedis(game dto.Game) error {
	gameJSON, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("failed to process game data: %v", err)
	}

	timeLimit := config.GetTimeLimit()
	err = config.AppConfig.RedisClient.Set(config.AppConfig.Ctx, game.ID, gameJSON, timeLimit).Err()
	if err != nil {
		return fmt.Errorf("error setting game in Redis: %v", err)
	}

	return nil
}

func (g *gameServiceImpl) MapperGameToPublicGame(game dto.Game) dto.PublicGame {
	publicGame := &dto.PublicGame{
		ID:         game.ID,
		LastMove:   game.LastMove,
		UserGrid:   game.UserGrid,
		HouseGrid:  game.HouseGrid,
		UpdatedAt:  game.UpdatedAt,
		CreatedAt:  game.CreatedAt,
		FinishedAt: game.FinishedAt,
		Winner:     game.Winner,
	}

	return *publicGame
}

func (g *gameServiceImpl) CheckIfIsInTheTimeLimit(game dto.Game) bool {
	timeLimit := config.GetTimeLimit()
	return time.Since(game.CreatedAt) <= timeLimit
}
