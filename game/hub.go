package game

import (
	"sync"

	"github.com/google/uuid"
)

type hub struct {
	mtx      sync.RWMutex
	roomColl map[string]*Room
}

func (h *hub) AddRoomToHub() *Room {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	// Loop until the produced id is unique
	var id string
	for {
		id = uuid.New().String()
		if _, ok := h.roomColl[id]; !ok {
			break
		}
	}

	room := CreateRoom(id)
	h.roomColl[id] = room
	return room
}

func (h *hub) FindRoomByID(id string) *Room {
	h.mtx.RLock()
	defer h.mtx.RUnlock()

	room, ok := h.roomColl[id]
	if !ok {
		return nil
	}

	return room
}

func newHub() *hub {
	return &hub{
		roomColl: make(map[string]*Room),
	}
}

var (
	HubObj *hub = newHub()
)
