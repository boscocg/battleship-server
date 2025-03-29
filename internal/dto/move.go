package dto

type Move struct {
	PublicGame
	Cell  int  `json:"cell" binding:"required" validate:"required"`
	IsHit bool `json:"is_hit" binding:"required" validate:"required"`
}

type MoveRequest struct {
	ID     uint64 `json:"id" binding:"required" validate:"required"`
	Player string `json:"player" binding:"required" validate:"required"`
	Cell   int    `json:"cell,omitempty"`
}

type CellType string

const (
	Hit  CellType = "hit"
	Miss CellType = "miss"
)
