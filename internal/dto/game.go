package dto

import (
	"time"
)

type Game struct {
	ID                 string     `json:"id" binding:"required" validate:"required"`
	LastMove           string     `json:"last_move" binding:"required" validate:"required"`
	UserGrid           []CellType `json:"user_grid" binding:"required" validate:"required"`
	HouseGrid          []CellType `json:"house_grid" binding:"required" validate:"required"`
	DecryptedHouseGrid []int      `json:"decrypted_house_grid" binding:"required" validate:"required"`
	UpdatedAt          time.Time  `json:"updated_at" binding:"required" validate:"required"`
	CreatedAt          time.Time  `json:"created_at" binding:"required" validate:"required"`
	FinishedAt         time.Time  `json:"finished_at" binding:"required" validate:"required"`
	Winner             string     `json:"winner"`
}

type GameRequest struct {
	UserGrid []CellType `json:"user_grid" binding:"required" validate:"required"`
}

type PublicGame struct {
	ID         string     `json:"id" binding:"required" validate:"required"`
	LastMove   string     `json:"last_move" binding:"required" validate:"required"`
	UserGrid   []CellType `json:"user_grid" binding:"required" validate:"required"`
	HouseGrid  []CellType `json:"house_grid" binding:"required" validate:"required"`
	UpdatedAt  time.Time  `json:"updated_at" binding:"required" validate:"required"`
	CreatedAt  time.Time  `json:"created_at" binding:"required" validate:"required"`
	FinishedAt time.Time  `json:"finished_at" binding:"required" validate:"required"`
	Winner     string     `json:"winner"`
}
