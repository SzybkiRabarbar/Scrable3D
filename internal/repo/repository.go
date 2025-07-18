package repo

import (
	"scrable3/internal/model"

	"github.com/google/uuid"
)

type Repository interface {
	Migrate() error
	CloseConn() error

	InsertGame(game *model.Game) error
	UpdateGame(game *model.Game) error
	SelectGameByUUID(gameUUID uuid.UUID) (*model.Game, error)
	DeleteGame(game *model.Game) error

	InsertPlayer(player *model.Player) error
	SelectPlayerByUUID(playerUUID uuid.UUID) (*model.Player, error)
	UpdatePlayer(updatedPlayer *model.Player) error

	InsertField(field *model.Field) error
	SelectFieldsByGameID(gameUUID uuid.UUID) (*[]model.Field, error)
	DeleteField(field *model.Field) error

	InsertAvChar(avChar *model.AvChar) error
	SelectAvCharsByPlayerID(playerUUID uuid.UUID) (*[]model.AvChar, error)
	DeleteAvCharByID(avCharID int64) error
}
