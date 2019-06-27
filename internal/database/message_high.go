package database

import (
	"database/sql"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// CreateMessage create message
func (db *DataBase) CreateMessage(message *models.Message,
	inRoom bool, gameID string) (id int, err error) {
	sqlInsert := `
	INSERT INTO GameChat(player_id, name, in_room, roomID, message, time) VALUES
		($1, $2, $3, $4, $5, $6)
		RETURNING ID;
		`
	row := db.Db.QueryRow(sqlInsert, message.User.ID, message.User.Name, inRoom,
		gameID, message.Text, message.Time)

	if err = row.Scan(&id); err != nil {
		utils.Debug(true, "cant create message", err.Error())
		return
	}

	return
}

// UpdateMessage update message
func (db *DataBase) UpdateMessage(message *models.Message) (id int, err error) {
	sqlInsert := `
	Update GameChat set message = $1, edited = true where id = $2
		RETURNING ID;
		`
	row := db.Db.QueryRow(sqlInsert, message.Text, message.ID)

	if err = row.Scan(&id); err != nil {
		utils.Debug(true, "cant update message", err.Error())
		return
	}

	return
}

// DeleteMessage delete message
func (db *DataBase) DeleteMessage(message *models.Message) (id int, err error) {
	sqlInsert := `
	Delete From GameChat where id = $1
		RETURNING ID;
		`
	row := db.Db.QueryRow(sqlInsert, message.ID)

	if err = row.Scan(&id); err != nil {
		utils.Debug(true, "cant delete message", err.Error())
		return
	}

	return
}

// LoadMessages load messages from database
func (db *DataBase) LoadMessages(inRoom bool, gameID string) (messages []*models.Message, err error) {

	var (
		tx *sql.Tx
	)

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if messages, err = db.getMessages(tx, inRoom, gameID); err != nil {
		utils.Debug(true, "cant load messages", err.Error())
		return
	}

	err = tx.Commit()
	return
}
