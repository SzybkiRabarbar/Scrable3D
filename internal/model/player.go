package model

import (
	"time"

	"github.com/google/uuid"
)

type Player struct {
	UUID       uuid.UUID
	CreateDate time.Time
	UpdateDate time.Time
	GameUUID   uuid.UUID
	Points     int64
	Appends    int
}

var PlayerMigrationSQL = map[string]string{
	"sqlite3": `-- Player
CREATE TABLE IF NOT EXISTS players(
    uuid BLOB PRIMARY KEY,
	create_date INTEGER NOT NULL,
	update_date INTEGER,
    game_uuid BLOB NOT NULL,
    points INTEGER NOT NULL,
    appends INTEGER NOT NULL,
    FOREIGN KEY (game_uuid) REFERENCES games(uuid) ON DELETE CASCADE
);
`,
}
