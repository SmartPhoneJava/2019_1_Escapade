package database

import (
	idb "github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

/*
	[ObjectName]UseCase opens a connection to the database, but cannot perform
	operations on it. All operations are performed using [ObjectName]Repository
	classes. It([ObjectName]Repository) cant open a database connection itself
	and all it's methods require a database connection as input data.
	[ObjectName]UseCase is responsible for closing the connection,
	opening transactions and other actions related directly to the connection.
	[ObjectName]UseCase call methods of [ObjectName]Repository to perform
	operations in database
*/

type UserUseCaseI interface {
	idb.UserCaseI
	Init(user UserRepositoryI, record RecordRepositoryI)

	EnterAccount(name, password string) (int32, error)
	CreateAccount(user *models.UserPrivateInfo) (int, error)
	UpdateAccount(userID int32, user *models.UserPrivateInfo) (err error)
	DeleteAccount(user *models.UserPrivateInfo) error

	FetchAll(difficult int, page int, perPage int,
		sort string) (players []*models.UserPublicInfo, err error)
	FetchOne(userID int32, difficult int) (*models.UserPublicInfo, error)

	PagesCount(perPage int) (int, error)
}

type UserRepositoryI interface {
	create(tx idb.TransactionI, user *models.UserPrivateInfo) (int, error)
	delete(tx idb.TransactionI, user *models.UserPrivateInfo) error

	updateNamePassword(tx idb.TransactionI, user *models.UserPrivateInfo) error
	checkNamePassword(tx idb.TransactionI, name string,
		password string) (int32, *models.UserPublicInfo, error)
	fetchNamePassword(tx idb.TransactionI,
		userID int32) (*models.UserPrivateInfo, error)

	updateLastSeen(tx idb.TransactionI, id int) error

	fetchAll(tx idb.TransactionI, params UsersSelectParams) ([]*models.UserPublicInfo, error)
	fetchOne(tx idb.TransactionI, userID int32,
		difficult int) (*models.UserPublicInfo, error)

	pagesCount(dbI idb.DatabaseI, perPage int) (int, error)
}

type GameUseCaseI interface {
	idb.UserCaseI
	Init(game GameRepositoryI)

	Create(game *models.Game) (int32, int32, error)
	Save(info models.GameInformation) error

	FetchOneGame(roomID string) (models.GameInformation, error)
	FetchAllGames(userID int32) ([]models.GameInformation, error)

	FetchAllRoomsID(userID int32) ([]string, error)
}

type GameRepositoryI interface {
	createGame(tx idb.TransactionI, game *models.Game) (int32, error)
	updateGame(tx idb.TransactionI, game *models.Game) error

	createGamers(tx idb.TransactionI, GameID int32, gamers []models.Gamer) error
	createField(tx idb.TransactionI, gameID int32,
		field models.Field) (int32, error)
	createActions(tx idb.TransactionI, GameID int32, actions []models.Action) error
	createCells(tx idb.TransactionI, FieldID int32, cells []models.Cell) error

	fetchOneGame(tx idb.TransactionI, roomID string) (models.Game, error)
	fetchAllCells(tx idb.TransactionI, fieldID int) ([]models.Cell, error)
	fetchAllGamers(tx idb.TransactionI, gameID int32) ([]models.Gamer, error)
	fetchAllActions(tx idb.TransactionI, gameID int32) ([]models.Action, error)
	fetchOneField(tx idb.TransactionI, gameID int32) (int, models.Field, error)

	fetchAllRoomsID(tx idb.TransactionI, userID int32) ([]string, error)
}

type RecordUseCaseI interface {
	idb.UserCaseI
	Init(record RecordRepositoryI)

	Update(id int32, record *models.Record) error
}

type RecordRepositoryI interface {
	create(tx idb.TransactionI, id int) error
	update(tx idb.TransactionI, id int32, record *models.Record) error
}

type ImageUseCaseI interface {
	idb.UserCaseI
	Init(image ImageRepositoryI)

	Update(filename string, userID int32) error
	FetchByName(name string) (string, error)
	FetchByID(id int32) (string, error)
}

type ImageRepositoryI interface {
	update(dbI idb.DatabaseI, filename string, userID int32) error
	fetchByName(dbI idb.DatabaseI, name string) (string, error)
	fetchByID(dbI idb.DatabaseI, id int32) (string, error)
}
