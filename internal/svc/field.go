package svc

import (
	"scrable3/internal/dto"
	"scrable3/internal/model"
	"scrable3/internal/repo"
	"time"

	"github.com/google/uuid"
)

type FieldService interface {
	Create(
		gameUUID uuid.UUID,
		playerUUID uuid.UUID,
		playerAppendNum int,
		value string,
		pos [3]int,
	) (*model.Field, error)
	CreateMany(
		gameUUID uuid.UUID,
		playerUUID uuid.UUID,
		playerAppendNum int,
		data *[]dto.FieldData,
	) (*[]model.Field, error)
	GetWithGameUUID(gameUUID uuid.UUID) (*[]model.Field, error)
	Delete(field *model.Field) error
	// CreateInitialData(
	// 	game *model.Game,
	// 	player *model.Player,
	// )( *model.Field, error)
}

type fieldService struct {
	repository repo.Repository
}

func NewFieldService(r repo.Repository) FieldService {
	return &fieldService{
		repository: r,
	}
}

func (service *fieldService) Create(
	gameUUID uuid.UUID,
	playerUUID uuid.UUID,
	playerAppendNum int,
	value string,
	pos [3]int,
) (*model.Field, error) {
	field := &model.Field{
		CreateDate: time.Now(),
		UpdateDate: time.Now(),
		GameUUID:   gameUUID,
		PlayerUUID: playerUUID,
		AppendNum:  playerAppendNum,
		Value:      value,
		PosX:       pos[0],
		PosY:       pos[1],
		PosZ:       pos[2],
	}
	err := service.repository.InsertField(field)
	return field, err
}

func (service *fieldService) CreateMany(
	gameUUID uuid.UUID,
	playerUUID uuid.UUID,
	playerAppendNum int,
	data *[]dto.FieldData,
) (*[]model.Field, error) {
	fields := []model.Field{}
	for _, fieldData := range *data {
		field := model.Field{
			GameUUID:   gameUUID,
			PlayerUUID: playerUUID,
			AppendNum:  playerAppendNum,
			Value:      fieldData.Value,
			PosX:       fieldData.Pos[0],
			PosY:       fieldData.Pos[1],
			PosZ:       fieldData.Pos[2],
		}
		err := service.repository.InsertField(&field)
		if err != nil {
			return &fields, err
		}
		fields = append(fields, field)
	}
	return &fields, nil
}

func (service *fieldService) GetWithGameUUID(gameUUID uuid.UUID) (*[]model.Field, error) {
	fields, err := service.repository.SelectFieldsByGameID(gameUUID)
	return fields, err
}

func (service *fieldService) Delete(field *model.Field) error {
	return service.repository.DeleteField(field)
}
