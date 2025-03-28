package service

type MoveService interface {
	GenerateHouseMove() (error, int)
}

type moveServiceImpl struct {
}

func NewMoveService() *moveServiceImpl {
	return &moveServiceImpl{}
}

func (g *moveServiceImpl) GenerateHouseMove() (error, int) {
	return nil, 0
}
