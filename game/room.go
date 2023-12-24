package game

import (
	"fmt"
	"sync"
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
}

func (r *Room) Run() {
	for {
		select {
		case player := <-r.Unregister:
			r.mtx.Lock()
			if _, ok := r.Players[player.ID]; ok {
				delete(r.Players, player.ID)
				close(player.MsgToPlayer)
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
					close(player.MsgToPlayer)
					delete(r.Players, player.ID)
				}
			}
			r.mtx.Unlock()
		}
	}
}

func (r *Room) registerPlayer(p *Player) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if len(r.Players) >= 4 {
		return fmt.Errorf("a room can only hosts up to four people")
	}

	r.Players[p.ID] = p
	r.TurnOrder = append(r.TurnOrder, p.ID)
	return nil
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
