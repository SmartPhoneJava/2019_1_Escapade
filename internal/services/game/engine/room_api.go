package engine

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
	action_ "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine/action"
)

// RoomRequestsI handle requests to toom
// Strategy Pattern
type RoomRequestsI interface {
	Handle(conn *Connection, rr *RoomRequest)
}

// RoomAPI implement RoomRequestsI
type RoomAPI struct {
	s  synced.SyncI
	m  MessagesI
	c  RClientI
	e  EventsI
	se RSendI
}

func (room *RoomAPI) build(builder RBuilderI) {
	builder.BuildSync(&room.s)
	builder.BuildMessages(&room.m)
	builder.BuildConnectionEvents(&room.c)
	builder.BuildEvents(&room.e)
	builder.BuildSender(&room.se)
}

// Init configure dependencies with other components of the room
func (room *RoomAPI) Init(builder RBuilderI) {
	room.build(builder)
}

// Handle processes the request came from the user
func (room *RoomAPI) Handle(conn *Connection, rr *RoomRequest) {
	utils.Debug(false, "start")
	go room.s.Do(func() {
		if rr.IsGet() {
			room.GetRoom(conn)
		} else if rr.IsSend() {
			room.handleSent(conn, rr.Send)
		} else if rr.Message != nil {
			room.handleMessage(conn, rr.Message)
		}
	})
}

func (room *RoomAPI) handleMessage(conn *Connection, message *models.Message) {
	room.s.DoWithOther(conn, func() {
		HandleMessage(conn, message, room.m)
	})
}

func (room *RoomAPI) handleSent(conn *Connection, request *RoomSend) {
	switch {
	case request.Messages != nil:
		room.GetMessages(conn, request.Messages)
	case request.Cell != nil:
		utils.Debug(false, "PostCell")
		room.PostCell(conn, request.Cell)
	case request.Action != nil:
		room.PostAction(conn, *request.Action)
	}
}

// GetRoom handle "GET /room", return all room information
func (room *RoomAPI) GetRoom(conn *Connection) {
	room.s.DoWithOther(conn, func() {
		room.se.Room(conn)
	})
}

// GetMessages handle "GET /messages", return all room messages
func (room *RoomAPI) GetMessages(conn *Connection, settings *models.Messages) {
	room.s.DoWithOther(conn, func() {
		Messages(conn, settings, room.m.Messages())
	})
}

// PostCell  handle "POST /cell" processes the Cell came from the user
func (room *RoomAPI) PostCell(conn *Connection, cell *Cell) {
	utils.Debug(false, "PostCell try")
	room.s.DoWithOther(conn, func() {
		utils.Debug(false, "PostCell try do")
		room.e.OpenCell(conn, cell)
	})
}

// PostAction handle "POST /action" processes the Cell came from the user
func (room *RoomAPI) PostAction(conn *Connection, action int) {
	room.s.DoWithOther(conn, func() {
		switch action {
		case action_.BackToLobby:
			room.c.BackToLobby(conn)
		case action_.Disconnect:
			room.c.Disconnect(conn)
		case action_.Reconnect:
			room.c.Reconnect(conn)
		case action_.GiveUp:
			room.c.GiveUp(conn)
		case action_.Restart:
			room.c.Restart(conn)
		}
	})
}
