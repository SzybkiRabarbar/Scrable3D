package svc

import (
	"scrable3/internal/model"
	"scrable3/internal/repo"
	"time"

	"github.com/google/uuid"
)

type GameService interface {
	Create() (*model.Game, error)
	GetWithUUID(gameUUID uuid.UUID) (*model.Game, error)
	Update(game *model.Game) error
	Delete(game *model.Game) error
	Refresh(game *model.Game) error
}

type gameService struct {
	repository repo.Repository
}

func NewGameService(r repo.Repository) GameService {
	return &gameService{
		repository: r,
	}
}

func (service *gameService) Create() (*model.Game, error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	game := &model.Game{
		UUID:        newUUID,
		CreateDate:  time.Now(),
		UpdateDate:  time.Now(),
		Turn:        int(0),
		PointsToWin: 20,
	}
	err = service.repository.InsertGame(game)
	return game, err
}

func (service *gameService) GetWithUUID(
	gameUUID uuid.UUID,
) (*model.Game, error) {
	game, err := service.repository.SelectGameByUUID(gameUUID)
	return game, err
}

func (service *gameService) Update(game *model.Game) error {
	game.UpdateDate = time.Now()
	err := service.repository.UpdateGame(game)
	return err
}

func (service *gameService) Delete(game *model.Game) error {
	err := service.repository.DeleteGame(game)
	return err
}

func (service *gameService) Refresh(game *model.Game) error {
	newGame, err := service.repository.SelectGameByUUID(game.UUID)
	if err != nil {
		return err
	}
	*game = *newGame
	return nil
}
