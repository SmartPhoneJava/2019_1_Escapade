package handlers

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine"

	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/photo"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

func gameOnline(db *database.DataBase, c *config.Configuration) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}
		result := connect(rw, r, db, c)
		api.SendResult(rw, result)
	}
}

func connect(rw http.ResponseWriter, r *http.Request, db *database.DataBase,
	c *config.Configuration) api.Result {
	const place = "GameOnline"
	var (
		err    error
		userID int32
		ws     *websocket.Conn
		user   *models.UserPublicInfo
		roomID string
	)

	utils.Debug(false, "GameOnline")
	lobby := engine.GetLobby()
	if lobby == nil {
		return api.NewResult(http.StatusInternalServerError, place, nil, re.ServerWrapper(err))
	}

	roomID = api.StringFromPath(r, "id", "")

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  c.WebSocket.ReadBufferSize,
		WriteBufferSize: c.WebSocket.WriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	if userID, err = api.GetUserIDFromAuthRequest(r); err != nil {
		userID = lobby.Anonymous()
	}

	if userID < 0 {
		user = &models.UserPublicInfo{
			Name:    "Anonymous" + strconv.Itoa(rand.Intn(10000)),
			ID:      int32(userID),
			FileKey: photo.GetDefaultAvatar(),
		}
	} else {
		if user, err = db.GetUser(userID, 0); err != nil {
			return api.NewResult(http.StatusNotFound, place, nil, re.NoUserWrapper(err))
		}
	}
	photo.GetImages(user)

	if ws, err = upgrader.Upgrade(rw, r, rw.Header()); err != nil {
		if _, ok := err.(websocket.HandshakeError); ok {
			err = re.ErrorHandshake()
		} else {
			err = re.ErrorNotWebsocket()
		}
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	conn := engine.NewConnection(ws, user, lobby)
	conn.Launch(c.WebSocket, roomID)

	return api.NewResult(0, place, nil, nil)
}
