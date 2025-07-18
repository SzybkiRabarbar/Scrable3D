package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"scrable3/internal/model"
	"time"

	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

type sqlite3Repository struct {
	db DB
}

// TODO dodaj daty do wszyskiego, create, update w sqlite trzymaj date w int

func NewSqlite3Connection(db DB) (Repository, error) {
	repo := &sqlite3Repository{db: db}
	err := repo.Migrate()
	if err != nil {
		return nil, err
	}

	return repo, nil
}

// |PRIVATE| //

func (repo *sqlite3Repository) checkSqlErr(err error) error {
	if err == nil {
		return nil
	}

	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return ErrDuplicate
		}
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotExists
	}

	return err
}

// |PUBLIC| //

func (repo *sqlite3Repository) Migrate() error {
	// Enable foreign key support
	_, err := repo.db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return err
	}
	d := "sqlite3"
	migrationQuery := model.GameMigrationSQL[d] +
		model.PlayerMigrationSQL[d] +
		model.FieldMigrationSQL[d] +
		model.AvCharMigrationSQL[d]
	_, err = repo.db.Exec(migrationQuery)
	return err
}

func (repo *sqlite3Repository) CloseConn() error {
	return repo.db.Close()
}

// * Game * //

func (repo *sqlite3Repository) InsertGame(game *model.Game) error {
	_, err := repo.db.Exec(
		`INSERT INTO games(
			uuid, 
			create_date,
			update_date,
			turn, 
			points_to_win
		) values(?,?,?,?,?)`,
		game.UUID,
		game.CreateDate.Unix(),
		game.UpdateDate.Unix(),
		game.Turn,
		game.PointsToWin,
	)
	return repo.checkSqlErr(err)
}

func (repo *sqlite3Repository) SelectGameByUUID(
	gameUUID uuid.UUID,
) (*model.Game, error) {
	row := repo.db.QueryRow("SELECT * FROM games WHERE uuid = ?", gameUUID)

	var game model.Game
	var createDate int64
	var updateDate int64
	err := row.Scan(&game.UUID, &createDate, &updateDate, &game.Turn, &game.PointsToWin)
	game.CreateDate = time.Unix(createDate, 0)
	game.UpdateDate = time.Unix(updateDate, 0)
	return &game, repo.checkSqlErr(err)
}

func (repo *sqlite3Repository) UpdateGame(game *model.Game) error {
	res, err := repo.db.Exec(
		`UPDATE games SET 
			turn = ?, 
			points_to_win = ?,
			update_date = ?,
		WHERE uuid = ?`,
		game.Turn,
		game.PointsToWin,
		game.UpdateDate.Unix(),
		game.UUID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUpdateFailed
	}

	return nil
}

func (repo *sqlite3Repository) DeleteGame(game *model.Game) error {
	res, err := repo.db.Exec("DELETE FROM games WHERE uuid = ?", game.UUID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return err
}

// * Player * //

func (repo *sqlite3Repository) InsertPlayer(player *model.Player) error {
	fmt.Println(player)
	_, err := repo.db.Exec(
		`INSERT INTO players(
			uuid, 
			create_date,
			update_date,
			game_uuid, 
			points, 
			appends
		) values(?,?,?,?,?,?)`,
		player.UUID,
		player.CreateDate.Unix(),
		player.UpdateDate.Unix(),
		player.GameUUID,
		player.Points,
		player.Appends,
	)
	return repo.checkSqlErr(err)
}

func (repo *sqlite3Repository) SelectPlayerByUUID(
	playerUUID uuid.UUID,
) (*model.Player, error) {
	row := repo.db.QueryRow("SELECT * FROM players WHERE uuid = ?", playerUUID)

	var createDate int64
	var updateDate int64
	var player model.Player
	err := row.Scan(
		&player.UUID, &createDate, &updateDate,
		&player.GameUUID, &player.Points, &player.Appends,
	)
	player.CreateDate = time.Unix(createDate, 0)
	player.UpdateDate = time.Unix(updateDate, 0)
	if err = repo.checkSqlErr(err); err != nil {
		return nil, err
	}
	return &player, nil
}

func (repo *sqlite3Repository) UpdatePlayer(
	updatedPlayer *model.Player,
) error {
	res, err := repo.db.Exec(
		`UPDATE players SET 
			game_uuid = ?,
			update_date = ?,
			points = ?, 
			appends = ? 
		WHERE uuid = ?`,
		updatedPlayer.GameUUID,
		updatedPlayer.UpdateDate.Unix(),
		updatedPlayer.Points,
		updatedPlayer.Appends,
		updatedPlayer.UUID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUpdateFailed
	}

	return nil
}

// * Field * //

func (repo *sqlite3Repository) InsertField(field *model.Field) error {
	res, err := repo.db.Exec(`
		INSERT INTO fields(
			game_uuid, create_date, update_date, player_uuid, append_num,
			val, pos_x, pos_y, pos_z
		) values(
			?,?,?,?,?,?,?,?,?
		)`,
		field.GameUUID,
		field.CreateDate.Unix(),
		field.UpdateDate.Unix(),
		field.PlayerUUID,
		field.AppendNum,
		field.Value,
		field.PosX,
		field.PosY,
		field.PosZ,
	)
	if err = repo.checkSqlErr(err); err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	field.ID = id

	return nil
}

func (repo *sqlite3Repository) SelectFieldsByGameID(
	gameUUID uuid.UUID,
) (*[]model.Field, error) {
	rows, err := repo.db.Query(
		"SELECT * FROM fields WHERE game_uuid = ?",
		gameUUID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fields []model.Field

	for rows.Next() {
		var createDate int64
		var updateDate int64
		var lt model.Field
		err := rows.Scan(
			&lt.ID, &createDate, &updateDate,
			&lt.GameUUID, &lt.PlayerUUID, &lt.AppendNum, &lt.Value,
			&lt.PosX, &lt.PosY, &lt.PosZ,
		)
		lt.CreateDate = time.Unix(createDate, 0)
		lt.UpdateDate = time.Unix(updateDate, 0)
		if err = repo.checkSqlErr(err); err != nil {
			return &fields, err
		}
		fields = append(fields, lt)
	}
	if err = rows.Err(); err != nil {
		return &fields, err
	}
	return &fields, nil
}

func (repo *sqlite3Repository) DeleteField(field *model.Field) error {
	res, err := repo.db.Exec("DELETE FROM fields WHERE id = ?", field.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return err
}

// * AvChar * //

func (repo *sqlite3Repository) InsertAvChar(avChar *model.AvChar) error {
	res, err := repo.db.Exec(`
		INSERT INTO available_characters(
			create_date,
			update_date,
			player_uuid, 
			val
		) values(
			?,?,?,?
		)`,
		avChar.CreateDate.Unix(),
		avChar.UpdateDate.Unix(),
		avChar.PlayerUUID,
		avChar.Value,
	)
	if err = repo.checkSqlErr(err); err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	avChar.ID = id

	return nil
}

func (repo *sqlite3Repository) SelectAvCharsByPlayerID(
	playerUUID uuid.UUID,
) (*[]model.AvChar, error) {
	rows, err := repo.db.Query(
		"SELECT * FROM available_characters WHERE player_uuid = ?",
		playerUUID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var available_characters []model.AvChar

	for rows.Next() {
		var createDate int64
		var updateDate int64
		var lt model.AvChar
		err := rows.Scan(
			&lt.ID, &createDate, &updateDate, &lt.PlayerUUID, &lt.Value,
		)
		lt.CreateDate = time.Unix(createDate, 0)
		lt.UpdateDate = time.Unix(updateDate, 0)
		if err = repo.checkSqlErr(err); err != nil {
			return &available_characters, err
		}
		available_characters = append(available_characters, lt)
	}
	if err = rows.Err(); err != nil {
		return &available_characters, err
	}
	return &available_characters, nil
}

func (repo *sqlite3Repository) DeleteAvCharByID(avCharID int64) error {
	res, err := repo.db.Exec(
		"DELETE FROM available_characters WHERE id = ?", avCharID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return err
}
