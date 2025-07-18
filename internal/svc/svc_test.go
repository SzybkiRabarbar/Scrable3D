package svc

import (
	"errors"
	"scrable3/internal/dto"
	"scrable3/internal/mock"
	"scrable3/internal/model"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func raiseErr(t *testing.T, service string, method string, err error) {
	t.Errorf("%v.%v error; %v", service, method, err)
}

func TestGameService(t *testing.T) {
	sn := "GameService"
	mc := gomock.NewController(t)
	defer mc.Finish()

	mockRepo := mock.NewMockRepository(mc)

	gameService := NewGameService(mockRepo)

	// *
	mn := "Create()"
	mockRepo.EXPECT().InsertGame(gomock.Any()).Return(nil)

	createdGame, err := gameService.Create()
	if err != nil {
		raiseErr(t, sn, mn, err)
	}
	createdGameUUID := createdGame.UUID

	// *
	mn = "Update()"
	mockRepo.EXPECT().UpdateGame(gomock.Any()).Return(nil)

	createdGame.Turn = 1
	err = gameService.Update(createdGame)
	if err != nil {
		raiseErr(t, sn, mn, err)
	}
	if createdGame.Turn != 1 {
		err = errors.New("Unexpected data manipulation")
		raiseErr(t, sn, mn, err)
	}

	// *
	mn = "GetWithUUID()"
	mockRepo.EXPECT().
		SelectGameByUUID(createdGameUUID).
		Return(&model.Game{UUID: createdGameUUID}, nil)

	fetchedGame, err := gameService.GetWithUUID(createdGameUUID)
	if err != nil {
		raiseErr(t, sn, mn, err)
	}
	if fetchedGame.UUID != createdGameUUID {
		err = errors.New("Unexpected data manipulation")
		raiseErr(t, sn, mn, err)
	}

	// *
	mn = "Refresh()"
	mockRepo.EXPECT().SelectGameByUUID(createdGame.UUID).Return(createdGame, nil)

	err = gameService.Refresh(createdGame)
	if err != nil {
		raiseErr(t, sn, mn, err)
	}

	// *
	mn = "Delete()"
	mockRepo.EXPECT().DeleteGame(gomock.Any()).Return(nil)

	err = gameService.Delete(createdGame)
	if err != nil {
		raiseErr(t, sn, mn, err)
	}
}

func TestPlayerService(t *testing.T) {
	sn := "PlayerService"
	mc := gomock.NewController(t)
	defer mc.Finish()

	mockRepo := mock.NewMockRepository(mc)
	game := &model.Game{UUID: uuid.UUID{}}

	playerService := NewPlayerService(mockRepo)

	// *
	mn := "Create()"
	mockRepo.EXPECT().InsertPlayer(gomock.Any()).Return(nil)

	createdPlayer, err := playerService.Create(game)
	if err != nil {
		raiseErr(t, sn, mn, err)
	}
	createdPlayerUUID := createdPlayer.UUID

	// *
	mn = "Update()"
	mockRepo.EXPECT().UpdatePlayer(gomock.Any()).Return(nil)
	createdPlayer.Points = 150
	createdPlayer.Appends = 1

	err = playerService.Update(createdPlayer)
	if err != nil {
		raiseErr(t, sn, mn, err)
	}
	if createdPlayer.Points != 150 || createdPlayer.Appends != 1 {
		err = errors.New("Unexpected data manipulation")
		raiseErr(t, sn, mn, err)
	}

	// *
	mn = "GetWithUUID()"
	mockRepo.EXPECT().
		SelectPlayerByUUID(createdPlayerUUID).
		Return(&model.Player{UUID: createdPlayerUUID, GameUUID: game.UUID}, nil)

	fetchedPlayer, err := playerService.GetWithUUID(createdPlayerUUID)
	if err != nil {
		raiseErr(t, sn, mn, err)
	}
	if fetchedPlayer.UUID != createdPlayerUUID ||
		fetchedPlayer.GameUUID != game.UUID {
		err = errors.New("Unexpected data manipulation")
		raiseErr(t, sn, mn, err)
	}

	// *
	mn = "Refresh()"
	mockRepo.EXPECT().UpdatePlayer(gomock.Any()).Return(nil)

	err = playerService.Update(createdPlayer)
	if err != nil {
		raiseErr(t, sn, mn, err)
	}
}

func TestFieldService(t *testing.T) {
	sn := "FieldService"
	mc := gomock.NewController(t)
	defer mc.Finish()

	mockRepo := mock.NewMockRepository(mc)
	game := &model.Game{UUID: uuid.UUID{}}
	player := &model.Player{UUID: uuid.UUID{}}

	fieldService := NewFieldService(mockRepo)

	// *
	mn := "Create()"
	mockRepo.EXPECT().InsertField(gomock.Any()).Return(nil)

	createdField, err := fieldService.Create(game.UUID, player.UUID, player.Appends, "A", [3]int{1, 2, 3})
	if err != nil {
		raiseErr(t, sn, mn, err)
	}

	// *
	mn = "CreateMany()"
	mockRepo.EXPECT().InsertField(gomock.Any()).Return(nil)
	mockRepo.EXPECT().InsertField(gomock.Any()).Return(nil)
	mockRepo.EXPECT().InsertField(gomock.Any()).Return(nil)
	data := &[]dto.FieldData{
		{
			Value: "A",
			Pos:   [3]int{1, 2, 3},
		},
		{
			Value: "B",
			Pos:   [3]int{4, 5, 6},
		},
		{
			Value: "C",
			Pos:   [3]int{7, 8, 9},
		},
	}

	createdFields, err := fieldService.CreateMany(game.UUID, player.UUID, player.Appends, data)
	if err != nil {
		raiseErr(t, sn, mn, err)
	}
	for i, f := range *createdFields {
		if f.Value != (*data)[i].Value {
			err := errors.New("Value not matching")
			raiseErr(t, sn, mn, err)
		}
	}

	// *
	mn = "GetWithGameUUID()"
	mockRepo.EXPECT().SelectFieldsByGameID(game.UUID).Return(createdFields, nil)

	fetchedFields, err := fieldService.GetWithGameUUID(game.UUID)
	if err != nil {
		raiseErr(t, sn, mn, err)
	}
	if fetchedFields != createdFields {
		err = errors.New("Wrong data returned")
		raiseErr(t, sn, mn, err)
	}
	// *
	mn = "Delete()"
	mockRepo.EXPECT().DeleteField(gomock.Any()).Return(nil)

	err = fieldService.Delete(createdField)
	if err != nil {
		raiseErr(t, sn, mn, err)
	}
}

func TestAvCharService(t *testing.T) {
	sn := "AvCharService"
	mc := gomock.NewController(t)
	defer mc.Finish()

	mockRepo := mock.NewMockRepository(mc)
	game := &model.Game{UUID: uuid.UUID{}}
	player := &model.Player{UUID: uuid.UUID{}, GameUUID: game.UUID}

	avCharService := NewAvCharService(mockRepo)

	//*
	mn := "CreateMany()"
	mockRepo.EXPECT().InsertAvChar(gomock.Any()).Return(nil)
	mockRepo.EXPECT().InsertAvChar(gomock.Any()).Return(nil)
	mockRepo.EXPECT().InsertAvChar(gomock.Any()).Return(nil)

	availableChars, err := avCharService.CreateMany(player, 3)
	if err != nil {
		raiseErr(t, sn, mn, err)
	}
	if len(*availableChars) != 3 {
		err := errors.New("Wrong number of AvChars")
		raiseErr(t, sn, mn, err)
	}
}
