package game

import (
	"encoding/json"
	"escapade/internal/models"
	"escapade/internal/utils"
	"fmt"
	"sync"
)

// Request connect Connection and his message
type Request struct {
	Connection *Connection
	Message    []byte
}

// Lobby there are all rooms and users placed
type Lobby struct {
	AllRooms  *Rooms `json:"allRooms,omitempty"`
	FreeRooms *Rooms `json:"freeRooms,omitempty"`

	Waiting map[int]*Connection `json:"waiting,omitempty"`
	Playing map[int]*Connection `json:"playing,omitempty"`

	// connection joined lobby
	ChanJoin chan *Connection `json:"-"`
	// connection left lobby
	chanLeave chan *Connection
	//chanRequest   chan *LobbyRequest
	chanBroadcast chan *Request

	semJoin    chan bool
	semRequest chan bool
}

// NewLobby create new instance of Lobby
func NewLobby(roomsCapacity, maxJoin, maxRequest int) *Lobby {

	lobby := &Lobby{
		AllRooms:  NewRooms(roomsCapacity),
		FreeRooms: NewRooms(roomsCapacity),

		Waiting: make(map[int]*Connection),
		Playing: make(map[int]*Connection),

		ChanJoin:  make(chan *Connection),
		chanLeave: make(chan *Connection),
		//chanRequest:   make(chan *LobbyRequest),
		chanBroadcast: make(chan *Request),

		semJoin:    make(chan bool, maxJoin),
		semRequest: make(chan bool, maxRequest),
	}
	return lobby
}

// CloseRoom free room resources
func (lobby *Lobby) CloseRoom(room *Room) {
	// if not in freeRooms nothing bad will happen
	// there is check inside, it will just return without errors
	lobby.FreeRooms.Remove(room)
	lobby.AllRooms.Remove(room)
	lobby.sendTAILRooms()
}

// createRoom create room, add to all and free rooms
// and run it
func (lobby *Lobby) createRoom(rs *models.RoomSettings) *Room {

	name := utils.RandomString(16) // вынести в кофиг
	room := NewRoom(rs, name, lobby)
	if !lobby.AllRooms.Add(room) {
		return nil
	}

	lobby.FreeRooms.Add(room)
	go lobby.sendTAILRooms() // inform all about new room
	return room
}

// Join handle user join to lobby
func (lobby *Lobby) Join(new *Connection) {
	lobby.semJoin <- true
	defer func() {
		<-lobby.semJoin
	}()

	// find such player
	old := lobby.AllRooms.SearchPlayer(new)

	if old != nil {
		old.room.RecoverPlayer(old, new)
		return // found
	}

	// find such observer
	old = lobby.AllRooms.SearchObserver(new)
	if old != nil {
		old.room.RecoverObserver(old, new)
		return // found
	}

	// player is new
	lobby.sendRooms(new)
	lobby.addWaiter(new)
	new.debug("new waiter")
	lobby.sendTAILPeople()
}

// Leave handle user leave lobby
func (lobby *Lobby) Leave(conn *Connection) {

	fmt.Println("disconnected -  #", conn.GetPlayerID())
	if !conn.InRoom() {
		lobby.removeWaiter(conn)
		lobby.sendTAILPeople()
	} else {
		lobby.removePlayer(conn)
		conn.room.addAction(conn, ActionDisconnect)
		conn.room.sendHistory(conn.room.all())
	}
	conn.Kill([]byte("Good day!"))
	return
}

// ----- handle room status
// roomStart - room remove from free
func (lobby *Lobby) roomStart(room *Room) {
	lobby.FreeRooms.Remove(room)
	lobby.sendTAILRooms()
}

// roomFinish - room remove from all
func (lobby *Lobby) roomFinish(room *Room) {
	room.Status = StatusFinished
	for _, conn := range room.Players.Get {
		lobby.playerToWaiter(conn)
	}
	lobby.AllRooms.Remove(room)
	lobby.sendTAILRooms()
}

// -----

// ----- handle connection status
func (lobby *Lobby) addWaiter(conn *Connection) {
	lobby.Waiting[conn.GetPlayerID()] = conn
}

func (lobby *Lobby) setWaiterRoom(conn *Connection, room *Room) {
	lobby.Waiting[conn.GetPlayerID()] = conn
}

func (lobby *Lobby) addPlayer(conn *Connection, room *Room) {
	lobby.Playing[conn.GetPlayerID()] = conn
}

func (lobby *Lobby) removeWaiter(conn *Connection) {
	delete(lobby.Waiting, conn.GetPlayerID())
}

func (lobby *Lobby) removePlayer(conn *Connection) {
	delete(lobby.Playing, conn.GetPlayerID())
}

func (lobby *Lobby) waiterToPlayer(conn *Connection) {
	lobby.removeWaiter(conn)
	lobby.addPlayer(conn, conn.room)
}

func (lobby *Lobby) playerToWaiter(conn *Connection) {
	lobby.removePlayer(conn)
	lobby.addWaiter(conn) // it send finished player to lobby
}

// -----

// pickUpRoom find room for player
func (lobby *Lobby) pickUpRoom(conn *Connection, rs *models.RoomSettings) (done bool) {
	// if there is no room
	if lobby.FreeRooms.Empty() {
		// if room capacity ended return nil
		room := lobby.createRoom(rs)
		if room != nil {
			conn.debug("We create your own room, cool!")
			room.addPlayer(conn)
		} else {
			conn.debug("cant create. Why?")
		}
		return room != nil
	}
	conn.debug("We have some rooms!")

	// lets find room for him
	for _, room := range lobby.FreeRooms.Get {
		//if room.SameAs()
		if room.addPlayer(conn) {
			done = true
			break
		}
	}
	if !done {
		Answer(conn, []byte("Error. Cant find room"))
	}
	return done
}

// handleRequest
func (lobby *Lobby) handleRequest(conn *Connection, lr *LobbyRequest) {
	conn.debug("lobby handle conn")
	lobby.semRequest <- true
	defer func() {
		<-lobby.semRequest
	}()

	if lr.IsGet() {
		lobby.requestGet(conn, lr)
	} else if lr.IsSend() {
		if lr.Send.RoomSettings == nil {
			conn.debug("lobby cant execute request")
			return
		}
		lobby.EnterRoom(conn, lr.Send.RoomSettings)
	}
}

// EnterRoom handle user join to room
func (lobby *Lobby) EnterRoom(conn *Connection, rs *models.RoomSettings) {

	done := false
	if _, room := lobby.AllRooms.SearchRoom(rs.Name); room != nil {
		conn.debug("lobby found required room")
		done = room.Enter(conn)
	} else {
		conn.debug("lobby search room for you")
		done = lobby.pickUpRoom(conn, rs)
	}

	if done {
		conn.debug("lobby done")
		lobby.waiterToPlayer(conn)
		lobby.sendTAILPeople()
	} else {
		conn.debug("lobby cant execute request")
	}
}

// sendRooms send rooms info for user
func (lobby *Lobby) sendRooms(conn *Connection) {
	bytes, _ := json.Marshal(lobby.AllRooms)
	conn.SendInformation(bytes)
}

// Run the room in goroutine
func (lobby *Lobby) Run() {

	for {
		select {
		case connection := <-lobby.ChanJoin:
			go lobby.Join(connection)

		case message := <-lobby.chanBroadcast:
			go lobby.analize(message)

		case connection := <-lobby.chanLeave:
			go lobby.Leave(connection)
		}
	}
}

func (lobby *Lobby) analize(req *Request) {
	if !req.Connection.InRoom() {
		var send *LobbyRequest
		if err := json.Unmarshal(req.Message, &send); err != nil {
			bytes, _ := json.Marshal(err)
			req.Connection.SendInformation(bytes)
		} else {
			lobby.handleRequest(req.Connection, send)
		}
	} else {
		if req.Connection.room == nil {
			return
		}
		var send *RoomRequest
		if err := json.Unmarshal(req.Message, &send); err != nil {
			bytes, _ := json.Marshal(err)
			req.Connection.SendInformation(bytes)
		} else {
			req.Connection.room.handleRequest(req.Connection, send)
		}
	}
}

func (lobby *Lobby) sendToAllInLobby(info interface{}) {
	waitJobs := &sync.WaitGroup{}
	bytes, _ := json.Marshal(info)
	fmt.Println("sendToAllInLobby began")
	for _, conn := range lobby.Waiting {
		if !conn.InRoom() {
			waitJobs.Add(1)
			conn.debug("send response")
			conn.sendGroupInformation(bytes, waitJobs)
		} else {
			conn.debug("lobby know about you :)")
		}
	}
	waitJobs.Wait()
}

// send to all in lobby
func (lobby *Lobby) sendTAILRooms() {
	get := &LobbyGet{
		AllRooms:  true,
		FreeRooms: true,
	}
	fmt.Println("sendTAILRooms began")
	send := lobby.makeGetModel(get)
	lobby.sendToAllInLobby(send)
}

func (lobby *Lobby) sendTAILPeople() {
	get := &LobbyGet{
		Waiting: true,
		Playing: true,
	}
	fmt.Println("sendTAILPeople began")
	send := lobby.makeGetModel(get)
	lobby.sendToAllInLobby(send)
}

func (lobby *Lobby) makeGetModel(get *LobbyGet) *Lobby {
	sendLobby := &Lobby{}
	if get.AllRooms {
		sendLobby.AllRooms = lobby.AllRooms
	}
	if get.FreeRooms {
		sendLobby.FreeRooms = lobby.FreeRooms
	}
	if get.Waiting {
		sendLobby.Waiting = lobby.Waiting
	}
	if get.Playing {
		sendLobby.Playing = lobby.Playing
	}
	return sendLobby
}

func (lobby *Lobby) requestGet(conn *Connection, lr *LobbyRequest) {
	sendLobby := lobby.makeGetModel(lr.Get)
	bytes, _ := json.Marshal(sendLobby)
	conn.debug("lobby execute get request")
	conn.SendInformation(bytes)
}
