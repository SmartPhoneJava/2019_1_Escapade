package database

import (
	"database/sql"
	"escapade/internal/models"
	re "escapade/internal/return_errors"

	//
	_ "github.com/lib/pq"
)

// В будущем добавить, чтобы отдельно была проверка на
// на корректность, отдельно на sql  инъекции
func ValidatePrivateUI(user *models.UserPrivateInfo) (err error) {

	if !models.ValidatePassword(user.Password) {
		err = re.ErrorInvalidPassword()
		return
	}

	if !models.ValidatePlayerName(user.Name) && !models.ValidateEmail(user.Email) {
		err = re.ErrorInvalidNameOrEmail()
		return
	}

	return
}

// GetPlayerIDbyName get player's id by his hame
func (db *DataBase) GetPlayerIDbyName(username string) (id int, err error) {
	sqlStatement := `SELECT id FROM Player WHERE name = $1`
	row := db.Db.QueryRow(sqlStatement, username)

	err = row.Scan(&id)
	return
}

// GetPlayerNamebyID get player's name by his id
func (db *DataBase) GetPlayerNamebyID(id int) (username string, err error) {
	sqlStatement := `SELECT name FROM Player WHERE id = $1`
	row := db.Db.QueryRow(sqlStatement, id)

	err = row.Scan(&username)
	return
}

// GetNameByEmail get player's name by his email
func (db DataBase) GetNameByEmail(email string) (name string, err error) {
	sqlStatement := "SELECT name " +
		"FROM Player where email=$1"

	row := db.Db.QueryRow(sqlStatement, email)

	if err = row.Scan(&name); err != nil {
		return
	}
	return
}

// GetNameByEmail get player's name by his email
func (db DataBase) GetPasswordEmailByName(name string) (email string, password string, err error) {
	sqlStatement := "SELECT email, password " +
		"FROM Player where name like $1"

	row := db.Db.QueryRow(sqlStatement, name)

	if err = row.Scan(&email, &password); err != nil {
		return
	}
	return
}

// isNameUnique checks if there are Players with
// this('taken') name and returns corresponding error if yes
func (db DataBase) isNameUnique(taken string) error {
	sqlStatement := "SELECT name " +
		"FROM Player where name=$1"

	row := db.Db.QueryRow(sqlStatement, taken)

	var tmp string
	if err := row.Scan(&tmp); err != sql.ErrNoRows {
		if err == nil {
			return re.ErrorNameIstaken()
		}
		return err
	}
	return nil
}

// isEmailUnique checks if there are Players with
// this('taken') email and returns corresponding error if yes
func (db DataBase) isEmailUnique(taken string) error {
	sqlStatement := "SELECT name " +
		"FROM Player where email=$1"

	row := db.Db.QueryRow(sqlStatement, taken)

	var tmp string
	if err := row.Scan(&tmp); err != sql.ErrNoRows {
		if err == nil {
			return re.ErrorEmailIstaken()
		}
		return err
	}
	return nil
}

// confirmUnique confirm that user.Email and user.Name
// dont use by another Player
func (db DataBase) confirmUnique(user *models.UserPrivateInfo) (err error) {
	if err = db.isEmailUnique(user.Email); err != nil {
		return
	}
	err = db.isNameUnique(user.Name)
	return
}

func (db DataBase) checkBunch(field string, password string) (id int, err error) {
	// If checkBunchNamePass cant find brunch name-password
	if id, err = db.checkBunchNamePass(field, password); err != nil {
		// and checkBunchEmailPass cant find brunch email-password
		if id, err = db.checkBunchEmailPass(field, password); err != nil {
			err = re.ErrorWrongPassword()
			return // then password wrong
		}
	}
	err = nil
	return
}

// confirmRightPass checks that Player with such
// password and name exists
func (db DataBase) checkBunchNamePass(username string, password string) (id int, err error) {

	sqlStatement := "SELECT id FROM Player where name like $1 and password like $2"
	row := db.Db.QueryRow(sqlStatement, username, password)
	err = row.Scan(&id)
	return
}

// confirmRightPass checks that Player with such
// password and name exists
func (db DataBase) checkBunchEmailPass(email string, password string) (id int, err error) {
	sqlStatement := "SELECT id FROM Player where email like $1 and password like $2"
	row := db.Db.QueryRow(sqlStatement, email, password)
	err = row.Scan(&id)
	return
}

// confirmRightEmail checks that Player with such
// email and name exists
func (db DataBase) confirmEmailNamePassword(user *models.UserPrivateInfo) error {
	sqlStatement := "SELECT 1 FROM Player where name like $1 and password like $2 and email like $3"

	row := db.Db.QueryRow(sqlStatement, user.Name, user.Password, user.Email)
	var res int
	err := row.Scan(&res)
	return err
}

func (db *DataBase) deletePlayer(user *models.UserPrivateInfo) error {
	sqlStatement := `
	DELETE FROM Player where name=$1 and password=$2 and email=$3
		`
	_, err := db.Db.Exec(sqlStatement, user.Name,
		user.Password, user.Email)

	return err
}

func (db *DataBase) createPlayer(user *models.UserPrivateInfo) (id int, err error) {
	sqlInsert := `
	INSERT INTO Player(name, password, email) VALUES
    ($1, $2, $3);
		`
	_, err = db.Db.Exec(sqlInsert, user.Name, user.Password, user.Email)

	if err != nil {
		return
	}

	sqlGetID := `SELECT id FROM Player WHERE name = $1`
	row := db.Db.QueryRow(sqlGetID, user.Name)

	err = row.Scan(&id)
	return
}

// UpdatePlayerByName gets name of Player from
// relation Session, cause we know that user has session
func (db *DataBase) UpdatePlayerByName(curName string, user *models.UserPrivateInfo) (err error) {
	var (
		curEmail     string
		curPass      string
		sqlStatement string
		oldName      string
	)

	oldName = curName
	if curEmail, curPass, err = db.GetPasswordEmailByName(curName); err != nil {
		return
	}

	if user.Email != curEmail && user.Email != "" {
		if !models.ValidateEmail(user.Email) {
			return re.ErrorInvalidEmail()
		}
		if err = db.isEmailUnique(user.Email); err != nil {
			return re.ErrorInvalidEmail()
		}
		curEmail = user.Email
	}

	if user.Password != curPass && user.Password != "" {
		curPass = user.Password
	}

	if user.Name != curName && user.Name != "" {
		if !models.ValidateString(user.Name) {
			return re.ErrorInvalidName()
		}
		if err = db.isNameUnique(user.Name); err != nil {
			return re.ErrorInvalidName()
		}
		curName = user.Name
	}

	sqlStatement = `
			UPDATE Player 
			SET name = $1, email = $2, password = $3  
			WHERE name like $4
		`
	_, err = db.Db.Exec(sqlStatement, curName, curEmail, curPass, oldName)

	if err != nil {
		return
	}

	return
}
