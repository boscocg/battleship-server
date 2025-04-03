package config

import (
	"log"
	"os"
	"strconv"
)

// variable or returns the default value 10 if not specified
func GetGridSize() int {
	defaultSize := 10
	gridSizeStr := os.Getenv("GRID_SIZE")

	if gridSizeStr == "" {
		return defaultSize
	}

	gridSize, err := strconv.Atoi(gridSizeStr)
	if err != nil {
		log.Printf("Error converting GRID_SIZE, using default value: %v", err)
		return defaultSize
	}

	// Verification to ensure a reasonable value
	if gridSize < 5 || gridSize > 20 {
		log.Printf("GRID_SIZE outside allowed limits (5-20), using default value 10")
		return defaultSize
	}

	return gridSize
}

// GetTotalCells returns the total number of cells in the grid (gridSize^2)
func GetTotalCells() int {
	gridSize := GetGridSize()
	return gridSize * gridSize
}

func GetTotalHits() int {
	return 20
}
