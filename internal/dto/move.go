package dto

type Move struct {
	PublicGame
	Cell     int  `json:"cell,omitempty"`
	Succeded bool `json:"succeded" binding:"required" validate:"required"`
}

type MoveRequest struct {
	ID     uint64 `json:"id" binding:"required" validate:"required"`
	Player string `json:"player" binding:"required" validate:"required"`
	Cell   int    `json:"cell,omitempty"`
}
