package main

import (
	"sync"
)

type Hub struct {
	sync.RWMutex
	Clients          map[*Client]int64
	Broadcast        chan []byte
	Register         chan *Client
	Unregister       chan *Client
	TurnNumber       int   //-> tiap ganti giliran increment++
	CurrentlyDrawing int64 //(uid) -> dapet dari clients[turnNumber % len(clients)]
}

func newHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]int64),
		TurnNumber: 1,
	}
}

func (h *Hub) run() {
	defer func() {
		close(h.Broadcast)
		close(h.Register)
		close(h.Unregister)
		DeleteHub(h)
	}()

	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = client.UID
			if len(h.Clients) == 1 {
				h.CurrentlyDrawing = client.UID
			}
			h.Unlock()
		case client := <-h.Unregister:
			h.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}

			if len(h.Clients) == 0 {
				return
			}

			for cl := range h.Clients {
				if cl.Order > 1 {
					cl.Order--
				}
			}
			h.Unlock()
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
