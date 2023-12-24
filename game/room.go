package game

import (
	"fmt"
	"sync"
	"time"
)

type Room struct {
	// All players
	Players map[int]*Player
	// Inbound messages from the players.
	MsgFromPlayer chan []byte
	// Unregister requests from players.
	Unregister chan *Player
	// Unique id for the room
	ID string
	// Index of currently drawing player from turnOrder
	CurrentTurnIdx int
	// Player ids based on the order of who entered the room
	TurnOrder []int
	// Name of the object to be drawn
	Words []string
	// Lock
	mtx sync.Mutex
	// Extra check so player can't join room that's due to removal
	isClosed bool
}

func (r *Room) Run() {
	ticker := time.NewTicker(10 * time.Second)
	defer func() {
		ticker.Stop()
		close(r.MsgFromPlayer)
		close(r.Unregister)
		HubObj.removeRoomByID(r.ID)
	}()

	for {
		select {
		case player := <-r.Unregister:
			r.mtx.Lock()
			if _, ok := r.Players[player.ID]; ok {
				r.unregisterPlayer(player)
			}
			r.mtx.Unlock()
		case message := <-r.MsgFromPlayer:
			r.mtx.Lock()
			for _, player := range r.Players {
				select {
				case player.MsgToPlayer <- message:
				default:
					// Might happen when for some reason the client cannot
					// process the message, since MsgToPlayer is buffered
					// channel that means something wrong occured on player
					r.unregisterPlayer(player)
				}
			}
			r.mtx.Unlock()
		case <-ticker.C:
			r.mtx.Lock()
			if len(r.Players) == 0 {
				defer r.mtx.Unlock()
				return
			}
			r.mtx.Unlock()
		}
	}
}

func (r *Room) registerPlayer(p *Player) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if len(r.Players) >= 4 || r.isClosed {
		return fmt.Errorf("a room can only hosts up to four people")
	}

	r.Players[p.ID] = p
	r.TurnOrder = append(r.TurnOrder, p.ID)
	return nil
}

func (r *Room) unregisterPlayer(p *Player) {
	close(p.MsgToPlayer)
	delete(r.Players, p.ID)

	i := 0
	for i < len(r.TurnOrder) {
		if r.TurnOrder[i] == p.ID {
			break
		}
		i++
	}

	removedPlayerIsDrawing := p.ID == r.TurnOrder[r.CurrentTurnIdx]
	if removedPlayerIsDrawing {
		r.Words = r.Words[1:]
	}

	r.TurnOrder = append(r.TurnOrder[:i], r.TurnOrder[i+1:]...)
	if r.CurrentTurnIdx > len(r.TurnOrder)-1 {
		r.CurrentTurnIdx = 0
	}

	r.isClosed = len(r.Players) == 0
	// should update game info
}

func CreateRoom(roomID string) *Room {
	return &Room{
		Players:        make(map[int]*Player),
		MsgFromPlayer:  make(chan []byte),
		Unregister:     make(chan *Player),
		ID:             roomID,
		CurrentTurnIdx: 0,
		TurnOrder:      []int{},
		Words:          generateWordsInRandomOrder(),
	}
}

func generateWordsInRandomOrder() []string {
	l := len(ListOfItems)
	res := make([]string, l)

	i := 0
	for key := range ListOfItems {
		res[i] = key
		i++
	}

	return res
}
