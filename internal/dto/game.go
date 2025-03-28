package dto

import (
	"time"
)

type Game struct {
	ID                 uint64    `json:"id" binding:"required" validate:"required"`
	LastMove           string    `json:"last_move" binding:"required" validate:"required"`
	UserGrid           []string  `json:"user_grid" binding:"required" validate:"required"`
	HouseGrid          []string  `json:"house_grid" binding:"required" validate:"required"`
	DecryptedHouseGrid []int     `json:"decrypted_house_grid" binding:"required" validate:"required"`
	LastMoveTimestamp  time.Time `json:"updated_at" binding:"required" validate:"required"`
}

type GameRequest struct {
	HouseGrid []string `json:"house_grid" binding:"required" validate:"required"`
}

type PublicGameResponse struct {
	ID                uint64    `json:"id" binding:"required" validate:"required"`
	LastMove          string    `json:"last_move" binding:"required" validate:"required"`
	UserGrid          []string  `json:"user_grid" binding:"required" validate:"required"`
	HouseGrid         []string  `json:"house_grid" binding:"required" validate:"required"`
	LastMoveTimestamp time.Time `json:"updated_at" binding:"required" validate:"required"`
}
