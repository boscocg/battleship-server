package service

import (
	"battledak-server/internal/dto"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	config "battledak-server/configs"
)

type GameService interface {
	GenerateHouseGrid() ([]dto.CellType, []int)
	GetGameFromRedis(id uint64) (error, dto.Game)
	SetGameToRedis(game dto.Game) error
	MapperGameToPublicGame(game dto.Game) dto.PublicGame
}

type gameServiceImpl struct {
}

func NewGameService() *gameServiceImpl {
	return &gameServiceImpl{}
}

func (g *gameServiceImpl) GenerateHouseGrid() ([]dto.CellType, []int) {
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

	houseGrid := make([]dto.CellType, 100)
	publicGrid := make([]int, 100)

	// Initialize all cells to empty (water)
	for i := range houseGrid {
		randomDigit := rand.Intn(9) + 1
		houseGrid[i] = digitToHash[randomDigit]
		publicGrid[i] = randomDigit
	}

	// Place ships
	placeShips(houseGrid, publicGrid, 4, 1, digitToHash)
	placeShips(houseGrid, publicGrid, 3, 2, digitToHash)
	placeShips(houseGrid, publicGrid, 2, 3, digitToHash)
	placeShips(houseGrid, publicGrid, 1, 4, digitToHash)

	return houseGrid, publicGrid
}

func placeShips(houseGrid []dto.CellType, publicGrid []int, count, size int, digitToHash map[int]dto.CellType) {
	for range count {
		placed := false
		for !placed {
			// Randomly decide orientation: 0 for horizontal, 1 for vertical
			orientation := rand.Intn(2)
			var startPos int

			if orientation == 0 { // Horizontal
				// Ensure the ship doesn't go off the right edge
				row := rand.Intn(10)
				col := rand.Intn(10 - size + 1)
				startPos = row*10 + col

				// Check if positions are already occupied
				canPlace := true
				for j := range size {
					pos := startPos + j
					if publicGrid[pos] == 0 {
						canPlace = false
						break
					}
				}

				if canPlace {
					for j := range size {
						pos := startPos + j
						houseGrid[pos] = digitToHash[0] // Set to ship (represented by 0)
						publicGrid[pos] = 0
					}
					println("Placed ship at", startPos, "size", size, "orientation", orientation)
					placed = true
				}
			} else { // Vertical
				// Ensure the ship doesn't go off the bottom edge
				row := rand.Intn(10 - size + 1)
				col := rand.Intn(10)
				startPos = row*10 + col

				// Check if positions are already occupied
				canPlace := true
				for j := range size {
					pos := startPos + j*10
					if publicGrid[pos] == 0 {
						canPlace = false
						break
					}
				}

				if canPlace {
					for j := range size {
						pos := startPos + j*10
						houseGrid[pos] = digitToHash[0] // Set to ship (represented by 0)
						publicGrid[pos] = 0
					}
					println("Placed ship at", startPos, "size", size, "orientation", orientation)
					placed = true
				}
			}
		}
	}
}

func (g *gameServiceImpl) GetGameFromRedis(id uint64) (error, dto.Game) {
	game := &dto.Game{
		ID: id,
	}

	// Check if the game exists in Redis
	err := config.AppConfig.RedisClient.Get(config.AppConfig.Ctx, strconv.FormatUint(game.ID, 10)).Err()
	if err != nil {
		return fmt.Errorf("Game not found: %v", err), dto.Game{}
	}

	// Retrieve the game data from Redis
	val, err := config.AppConfig.RedisClient.Get(config.AppConfig.Ctx, strconv.FormatUint(game.ID, 10)).Bytes()
	if err != nil {
		return fmt.Errorf("Failed to retrieve game: %v", err), dto.Game{}
	}

	if err = json.Unmarshal(val, game); err != nil {
		return fmt.Errorf("Failed to parse game data: %v", err), dto.Game{}
	}

	return nil, *game
}

func (g *gameServiceImpl) SetGameToRedis(game dto.Game) error {
	gameJSON, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("Failed to process game data: %v", err)
	}

	// Set the game in Redis with a 1-hour expiration time
	err = config.AppConfig.RedisClient.Set(config.AppConfig.Ctx, strconv.FormatUint(game.ID, 10), gameJSON, 0).Err()
	if err != nil {
		return fmt.Errorf("Error setting game in Redis: %v", err)
	}

	return nil
}

func (g *gameServiceImpl) MapperGameToPublicGame(game dto.Game) dto.PublicGame {
	publicGame := &dto.PublicGame{
		ID:        game.ID,
		LastMove:  game.LastMove,
		UserGrid:  game.UserGrid,
		HouseGrid: game.HouseGrid,
		UpdatedAt: game.UpdatedAt,
		Finished:  game.Finished,
	}

	return *publicGame
}
