package service

import (
	"battledak-server/internal/dto"
	"fmt"
)

type CryptoService interface {
	IsZero(cellValue dto.CellType) (error, bool)
}

type cryptoServiceImpl struct {
}

func NewCryptoService() *cryptoServiceImpl {
	return &cryptoServiceImpl{}
}

func (g *cryptoServiceImpl) IsZero(cellValue dto.CellType) (error, bool) {
	if cellValue == "" || cellValue == "0" {
		return fmt.Errorf("Invalid value"), true
	}

	if cellValue == "5feceb66ffc86f38d952786c6d696c79c2dbc239dd4e91b46729d73a27fb57e9" {
		return nil, true
	}

	return nil, false
}
