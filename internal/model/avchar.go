package model

import (
	"time"

	"github.com/google/uuid"
)

type AvChar struct {
	ID         int64
	CreateDate time.Time
	UpdateDate time.Time
	PlayerUUID uuid.UUID
	Value      string
}

var AvCharMigrationSQL = map[string]string{
	"sqlite3": `-- AvChar
CREATE TABLE IF NOT EXISTS available_characters(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	create_date INTEGER NOT NULL,
	update_date INTEGER,
    player_uuid BLOB NOT NULL,
    val CHAR(1) NOT NULL,
    FOREIGN KEY (player_uuid) REFERENCES players(uuid) ON DELETE CASCADE
);
`,
}
