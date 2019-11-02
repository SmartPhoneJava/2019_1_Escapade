package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"database/sql"
	"time"
)

// createPlayer create player
func (db *DataBase) createPlayer(tx *sql.Tx, user *models.UserPrivateInfo) (int, error) {
	sqlInsert := `
	INSERT INTO Player(name, password, firstSeen, lastSeen) VALUES
		($1, $2, $3, $4)
		RETURNING id;
		`
	t := time.Now()
	row := tx.QueryRow(sqlInsert, user.Name,
		user.Password, t, t)

	var (
		id  int
		err error
	)
	err = row.Scan(&id)
	return id, err
}

func (db *DataBase) updatePlayerPersonalInfo(tx *sql.Tx, user *models.UserPrivateInfo) error {
	sqlStatement := `
			UPDATE Player 
			SET name = $1, password = $2, lastSeen = $3
			WHERE id = $4
			RETURNING id
		`

	row := tx.QueryRow(sqlStatement, user.Name,
		user.Password, time.Now(), user.ID)
	err := row.Scan(&user.ID)
	if err != nil {
		err = re.ErrorUserIsExist()
	}

	return err
}

// updatePlayerLastSeen update users last date seen
func (db *DataBase) updatePlayerLastSeen(tx *sql.Tx, id int) error {
	var (
		sqlStatement = `
			UPDATE Player 
			SET lastSeen = $1
			WHERE id = $2
		`
		err error
	)
	_, err = tx.Exec(sqlStatement, time.Now(), id)
	return err
}

func (db DataBase) checkBunch(tx *sql.Tx, name string, password string) (int32, *models.UserPublicInfo, error) {
	var (
		sqlStatement = `
			SELECT pl.id, pl.name, r.score, r.time, r.difficult
			FROM Player as pl
			join Record as r 
			on r.player_id = pl.id
			where r.difficult = 0 and password like $1 and name like $2`
		id int32
	)

	row := tx.QueryRow(sqlStatement, password, name)
	user := &models.UserPublicInfo{}
	err := row.Scan(&id, &user.Name, &user.BestScore, &user.BestTime, &user.Difficult)
	return id, user, err
}

// GetPrivateInfo get player's personal info
func (db DataBase) getPrivateInfo(tx *sql.Tx, userID int32) (*models.UserPrivateInfo, error) {

	sqlStatement := "SELECT name, password " +
		"FROM Player where id = $1"

	row := tx.QueryRow(sqlStatement, userID)
	user := &models.UserPrivateInfo{}
	user.ID = int(userID)
	err := row.Scan(&user.Name, &user.Password)

	return user, err
}

// GetUsers returns information about users
// for leaderboard
func (db *DataBase) getUsers(tx *sql.Tx, difficult int, offset int, limit int,
	sort string) ([]*models.UserPublicInfo, error) {

	sqlStatement := `
	SELECT P.id, P.photo_title, P.name,
				 R.score, R.time, R.Difficult
	FROM Player as P
	join Record as R 
	on R.player_id = P.id
	where r.difficult = $1  
	`
	if sort == "score" {
		sqlStatement += ` ORDER BY (score) desc `
	} else {
		sqlStatement += ` ORDER BY (time) `
	}
	sqlStatement += ` OFFSET $2 Limit $3 `

	players := make([]*models.UserPublicInfo, 0, limit)
	rows, err := tx.Query(sqlStatement, difficult, offset, limit)
	if err != nil {
		return players, err
	}
	defer rows.Close()

	for rows.Next() {
		player := &models.UserPublicInfo{}
		err = rows.Scan(&player.ID, &player.FileKey, &player.Name,
			&player.BestScore, &player.BestTime, &player.Difficult)
		if err != nil {
			break
		}
		players = append(players, player)
	}

	return players, err
}

// GetUser returns information about user
func (db *DataBase) getUser(tx *sql.Tx, userID int32, difficult int) (*models.UserPublicInfo, error) {

	sqlStatement := `
	SELECT P.id, P.photo_title, P.name,
				 R.score, R.time, R.Difficult
	FROM Player as P
	join Record as R 
	on R.player_id = P.id
	where R.player_id = $1 and
		R.difficult = $2
	`

	player := &models.UserPublicInfo{}
	row := tx.QueryRow(sqlStatement, userID, difficult)
	err := row.Scan(&player.ID, &player.FileKey, &player.Name,
		&player.BestScore, &player.BestTime, &player.Difficult)

	return player, err
}

// deletePlayer delete all information about user
func (db *DataBase) deletePlayer(tx *sql.Tx, user *models.UserPrivateInfo) error {
	sqlStatement := `
	DELETE FROM Player where name=$1 and password=$2
	RETURNING ID
		`
	row := tx.QueryRow(sqlStatement, user.Name, user.Password)

	err := row.Scan(&user.ID)
	if err != nil {
		utils.Debug(true, "cant delete player")
	}

	return err
}
