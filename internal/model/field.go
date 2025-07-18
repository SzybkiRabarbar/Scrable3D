package model

import (
	"time"

	"github.com/google/uuid"
)

type Field struct {
	ID         int64
	CreateDate time.Time
	UpdateDate time.Time
	GameUUID   uuid.UUID
	PlayerUUID uuid.UUID
	AppendNum  int
	Value      string
	PosX       int
	PosY       int
	PosZ       int
}

var FieldMigrationSQL = map[string]string{
	"sqlite3": `-- Field
CREATE TABLE IF NOT EXISTS fields(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	create_date INTEGER NOT NULL,
	update_date INTEGER,
    game_uuid BLOB NOT NULL,
    player_uuid BLOB NOT NULL,
    append_num INTEGER NOT NULL,
    val CHAR(1) NOT NULL,
    pos_x INTEGER NOT NULL,
    pos_y INTEGER NOT NULL,
    pos_z INTEGER NOT NULL,
    FOREIGN KEY (game_uuid) REFERENCES games(uuid) ON DELETE CASCADE,
    FOREIGN KEY (player_uuid) REFERENCES players(uuid) ON DELETE CASCADE
);
`,
}
