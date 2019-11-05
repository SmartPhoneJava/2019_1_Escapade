package handlers

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	idb "github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine"

	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	ih "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/photo"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"

	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type GameHandler struct {
	c          *config.Configuration
	upgraderWS websocket.Upgrader
	user       database.UserUseCaseI
	game       database.GameUseCaseI
}

func (h *GameHandler) InitWithPostgresql(c *config.Configuration) error {

	var (
		db     = &idb.PostgresSQL{}
		user   = &database.UserRepositoryPQ{}
		record = &database.RecordRepositoryPQ{}
		game   = &database.GameRepositoryPQ{}
	)
	return h.Init(c, db, user, record, game)
}

func (h *GameHandler) GameDB(c *config.Configuration) database.GameUseCaseI {
	return h.game
}

func (h *GameHandler) Init(c *config.Configuration,
	DB idb.DatabaseI,
	userDB database.UserRepositoryI,
	recordDB database.RecordRepositoryI,
	gameDB database.GameRepositoryI,
) error {

	err := DB.Open(c.DataBase)
	if err != nil {
		return err
	}

	h.c = c
	h.upgraderWS = websocket.Upgrader{
		ReadBufferSize:  c.WebSocket.ReadBufferSize,
		WriteBufferSize: c.WebSocket.WriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	h.user = &database.UserUseCase{}
	h.user.Init(userDB, recordDB)

	err = h.user.Use(DB)
	if err != nil {
		return err
	}

	h.game = &database.GameUseCase{}
	h.game.Init(gameDB)

	err = h.game.Use(DB)
	if err != nil {
		return err
	}
	return nil
}

func (h *GameHandler) Close() {
	h.user.Close()
}

func (h *GameHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	ih.Route(rw, r, ih.MethodHandlers{
		http.MethodGet:     h.connect,
		http.MethodOptions: nil})
}

func (h *GameHandler) connect(rw http.ResponseWriter, r *http.Request) api.Result {
	const place = "GameOnline"
	var (
		ws   *websocket.Conn
		user *models.UserPublicInfo
	)

	roomID := api.StringFromPath(r, "id", "")

	lobby := engine.GetLobby()
	if lobby == nil {
		return api.NewResult(http.StatusInternalServerError, place, nil, re.ErrorServer())
	}

	err := h.prepareUser(r, h.user, lobby)
	if err != nil {
		return api.NewResult(http.StatusInternalServerError, place, nil, err)
	}

	if ws, err = h.upgraderWS.Upgrade(rw, r, rw.Header()); err != nil {
		if _, ok := err.(websocket.HandshakeError); ok {
			err = re.ErrorHandshake()
		} else {
			err = re.ErrorNotWebsocket()
		}
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	conn := engine.NewConnection(ws, user, lobby)
	conn.Launch(h.c.WebSocket, roomID)
	// code 0 mean nothing to send to client
	return api.NewResult(0, place, nil, nil)
}

func (h *GameHandler) prepareUser(r *http.Request, userDB database.UserUseCaseI, lobby *engine.Lobby) error {
	var (
		userID int32
		err    error
		user   *models.UserPublicInfo
	)
	if userID, err = api.GetUserIDFromAuthRequest(r); err != nil {
		userID = lobby.Anonymous()
	}

	if userID < 0 {
		user = h.anonymousUser(userID)
	} else {
		if user, err = userDB.FetchOne(userID, 0); err != nil {
			return re.NoUserWrapper(err)
		}
	}
	photo.GetImages(user)
	return nil
}

func (h *GameHandler) anonymousUser(userID int32) *models.UserPublicInfo {
	return &models.UserPublicInfo{
		Name:    "Anonymous" + strconv.Itoa(rand.Intn(10000)), // в конфиг
		ID:      int32(userID),
		FileKey: photo.GetDefaultAvatar(),
	}
}
