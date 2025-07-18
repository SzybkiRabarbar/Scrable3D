package model

import (
	"time"

	"github.com/google/uuid"
)

type Game struct {
	UUID        uuid.UUID
	CreateDate  time.Time
	UpdateDate  time.Time
	Turn        int
	PointsToWin int64
}

var GameMigrationSQL = map[string]string{
	"sqlite3": `-- Game
CREATE TABLE IF NOT EXISTS games(
    uuid BLOB PRIMARY KEY,
	create_date INTEGER NOT NULL,
	update_date INTEGER,
    turn INTEGER NOT NULL,
    points_to_win INTEGER NOT NULL
);
`,
}
