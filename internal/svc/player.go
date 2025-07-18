package svc

import (
	"scrable3/internal/model"
	"scrable3/internal/repo"
	"time"

	"github.com/google/uuid"
)

type PlayerService interface {
	Create(game *model.Game) (*model.Player, error)
	GetWithUUID(playerUUID uuid.UUID) (*model.Player, error)
	Update(player *model.Player) error
	Refresh(player *model.Player) error
}

type playerService struct {
	repository repo.Repository
}

func NewPlayerService(r repo.Repository) PlayerService {
	return &playerService{
		repository: r,
	}
}

func (service *playerService) Create(game *model.Game) (*model.Player, error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	player := &model.Player{
		UUID:       newUUID,
		CreateDate: time.Now(),
		UpdateDate: time.Now(),
		GameUUID:   game.UUID,
		Points:     0,
		Appends:    0,
	}
	err = service.repository.InsertPlayer(player)
	return player, err
}

func (service *playerService) GetWithUUID(
	playerUUID uuid.UUID,
) (*model.Player, error) {
	player, err := service.repository.SelectPlayerByUUID(playerUUID)
	return player, err
}

func (service *playerService) Update(player *model.Player) error {
	player.UpdateDate = time.Now()
	err := service.repository.UpdatePlayer(player)
	return err
}

func (service *playerService) Refresh(player *model.Player) error {
	newPlayer, err := service.repository.SelectPlayerByUUID(player.UUID)
	if err != nil {
		return err
	}
	*player = *newPlayer
	return nil
}
