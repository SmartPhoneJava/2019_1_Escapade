package game

import (
	"fmt"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
)

// Winner determine who won the game
func (room *Room) Winner() (idWin int) {
	if room.done() {
		return 0
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	max := 0.

	players := room.Players.RPlayers()
	for id, player := range players {
		if player.Points > max {
			max = player.Points
			idWin = id
		}
	}
	return
}

// flagFound is called, when somebody find cell flag
func (room *Room) FlagFound(founder Connection, found *Cell) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	which := 0
	for _, flag := range room.Players.Flags() {
		if flag.Cell.X == found.X && flag.Cell.Y == found.Y {
			which = flag.Cell.PlayerID
		}
	}

	if which == founder.ID() {
		return
	}

	fmt.Println(which, "was found by", founder.ID())

	room.Players.IncreasePlayerPoints(founder.Index(), 30)

	killConn, index := room.Players.Connections.SearchByID(which)
	if index >= 0 {
		room.Kill(killConn, ActionFlagLost)
	}
}

// isAlive check if connection is player and he is not died
func (room *Room) isAlive(conn *Connection) bool {
	index := conn.Index()
	return index >= 0 && !room.Players.Player(index).Finished
}

// kill make user die and check for finish battle
func (room *Room) Kill(conn *Connection, action int) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	// cause all in pointers
	if room.Status > StatusRunning {
		return
	}

	if room.isAlive(conn) {
		room.SetFinished(conn)

		cell := room.Players.Flag(conn.Index())
		room.Field.SetCellFlagTaken(&cell.Cell)

		if room.Players.Capacity() <= room.killed()+1 {
			room.FinishGame(false)
		}
		pa := *room.addAction(conn.ID(), action)
		go room.sendAction(pa, room.All)
	}
	return
}

// GiveUp kill connection, that call it
func (room *Room) GiveUp(conn *Connection) {
	room.Kill(conn, ActionGiveUp)
}

// flagExists find players with such flag. This - flag owner
func (room *Room) flagExists(cell Cell, this *Connection) (found bool, conn *Connection) {
	var player int
	flags := room.Players.Flags()
	for index, flag := range flags {
		if (flag.Cell.X == cell.X) && (flag.Cell.Y == cell.Y) {
			if this == nil || index != this.Index() {
				found = true
				player = index
			}
			break
		}
	}
	if !found {
		return
	}
	conn = room.Players.Connections.SearchByIndex(player)
	return
}

func (room *Room) SetAndSendNewCell(conn Connection) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	found := true
	// create until it become unique
	var cell Cell
	for found {
		cell = room.Field.CreateRandomFlag(conn.ID())
		found, _ = room.flagExists(cell, nil)
	}
	if room.Players.SetFlag(conn, cell) {
		room.prepare.Stop()
		go room.StartGame()
		//
	}
	response := models.RandomFlagSet(cell)
	conn.SendInformation(response)
}

// setFlag handle user wanna set flag
func (room *Room) SetFlag(conn *Connection, cell *Cell) bool {
	if room.done() {
		return false
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	// if user try set flag after game launch
	if room.Status != StatusFlagPlacing {
		response := models.FailFlagSet(cell, re.ErrorBattleAlreadyBegan())
		conn.SendInformation(response)
		return false
	}

	if !room.Field.IsInside(cell) {
		response := models.FailFlagSet(cell, re.ErrorCellOutside())
		conn.SendInformation(response)
		return false
	}

	if !room.isAlive(conn) {
		response := models.FailFlagSet(cell, re.ErrorPlayerFinished())
		conn.SendInformation(response)
		return false
	}

	if found, prevConn := room.flagExists(*cell, conn); found {
		go room.SetAndSendNewCell(*conn)
		go room.SetAndSendNewCell(*prevConn)
		return true
	}

	if room.Players.SetFlag(*conn, *cell) {
		room.prepare.Stop()
		go room.StartGame()
		//
	}
	return true
}

// setFlags set players flags to field
// call it if game has already begun
func (room *Room) setFlags() {
	flags := room.Players.Flags()
	for _, cell := range flags {
		room.Field.SetFlag(&cell.Cell)
	}
}

// fillField set flags and mines
func (room *Room) FillField() {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	fmt.Println("fillField", room.Field.Height, room.Field.Width, len(room.Field.Matrix))

	room.setFlags()
	room.Field.SetMines()

}

func (room *Room) addAction(id int, action int) (pa *PlayerAction) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	pa = NewPlayerAction(id, action)
	room.setToHistory(pa)
	return
}
