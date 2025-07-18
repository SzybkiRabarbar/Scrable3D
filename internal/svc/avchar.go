package svc

import (
	"scrable3/internal/model"
	"scrable3/internal/repo"
	"time"

	"github.com/google/uuid"
)

type AvCharService interface {
	CreateMany(player *model.Player, n int) (*[]model.AvChar, error)
	GetWithPlayerUUID(playerUUID uuid.UUID) (*[]model.AvChar, error)
	DeleteMany(charsIDs *[]int64) error
}

type avCharService struct {
	repository repo.Repository
}

func NewAvCharService(r repo.Repository) AvCharService {
	return &avCharService{
		repository: r,
	}
}

var n = -1

func (service *avCharService) drawValue() string {
	// randomIndex := rand.Intn(len(cfg.ALLOWED_CHARACTERS))
	// return string(cfg.ALLOWED_CHARACTERS[randomIndex])
	s := "LAMPLAMPLAMPX"
	n += 1
	return string(s[n%len(s)])
}

func (service *avCharService) CreateMany(
	player *model.Player,
	n int,
) (*[]model.AvChar, error) {
	avChars := []model.AvChar{}
	for range n {
		randomValue := service.drawValue()
		avChar := model.AvChar{
			CreateDate: time.Now(),
			UpdateDate: time.Now(),
			PlayerUUID: player.UUID,
			Value:      randomValue,
		}
		err := service.repository.InsertAvChar(&avChar)
		if err != nil {
			return &avChars, err
		}
		avChars = append(avChars, avChar)
	}

	return &avChars, nil
}

func (service *avCharService) GetWithPlayerUUID(
	playerUUID uuid.UUID,
) (*[]model.AvChar, error) {
	avChars, err := service.repository.SelectAvCharsByPlayerID(playerUUID)
	return avChars, err
}

func (service *avCharService) DeleteMany(charsIDs *[]int64) error {
	for _, id := range *charsIDs {
		err := service.repository.DeleteAvCharByID(id)
		if err != nil {
			return err
		}
	}
	return nil
}
