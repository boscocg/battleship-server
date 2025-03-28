package service

import (
	"math/rand"
)

type GameService interface {
	GenerateHouseGrid() ([]string, []int)
}

type gameServiceImpl struct {
}

func NewGameService() *gameServiceImpl {
	return &gameServiceImpl{}
}

func (g *gameServiceImpl) GenerateHouseGrid() ([]string, []int) {
	// Seed the random number generator
	digitToHash := map[int]string{
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

	houseGrid := make([]string, 100)
	publicGrid := make([]int, 100)

	for i := 0; i < 100; i++ {
		randomDigit := rand.Intn(10)
		houseGrid[i] = digitToHash[randomDigit]
		publicGrid[i] = randomDigit
	}

	return houseGrid, publicGrid
}
