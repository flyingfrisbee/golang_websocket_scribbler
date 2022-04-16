package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
)

type Hub struct {
	sync.RWMutex
	Clients          map[*Client]int64
	Broadcast        chan []byte
	Register         chan *Client
	Unregister       chan *Client
	Words            []string
	TurnNumber       int   //-> tiap ganti giliran increment++
	CurrentlyDrawing int64 //(uid) -> dapet dari clients[turnNumber % len(clients)]
}

type GameStat struct {
	CurrentlyDrawing int64
	Answer           string
	Players          []PlayerStat
}

type PlayerStat struct {
	UID   int64
	Name  string
	Score int
}

func newHub() *Hub {
	words := []string{}
	for k := range ListOfItems {
		words = append(words, k)
	}

	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]int64),
		Words:      words,
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

			//tiap orang masuk balikin semua status: list[nama, skor, uid] & currently drawing
			h.RLock()
			ShowGameStatToPlayers(h)
			h.RUnlock()

		case client := <-h.Unregister:

			h.Lock()
			orderPlaceholder := client.Order

			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}

			if len(h.Clients) == 0 {
				return
			}

			for cl := range h.Clients {
				if cl.Order > orderPlaceholder {
					cl.Order--
				}
			}
			h.Unlock()

			//tiap orang keluar balikin game stat
			h.RLock()
			ShowGameStatToPlayers(h)
			h.RUnlock()

		case message := <-h.Broadcast:

			switch message[0] {

			case '3':

				h.Lock()
				msg := strings.Split(string(message), ",")
				receivedUID, _ := strconv.ParseInt(msg[1], 10, 64)
				answer := msg[2]

				for cl, uid := range h.Clients {
					if receivedUID == uid {
						cl.HasAnswered = true
						if strings.EqualFold(strings.Replace(h.Words[0], " ", "", 1), strings.Replace(strings.TrimSpace(answer), " ", "", 1)) {
							cl.Score += 2
						}
						break
					}
				}
				for cl, _ := range h.Clients {
					fmt.Println(*cl)
				}
				h.Unlock()

				h.RLock()
				ShowGameStatToPlayers(h)
				h.RUnlock()

			default:
				h.RLock()
				for client := range h.Clients {
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(h.Clients, client)
					}
				}
				h.RUnlock()
			}

			// next turn
			// h.Words = h.Words[1:]
			// if len(h.Words) == 0 {
			// 	//notify client game has ended
			// }
		}
	}
}

func ShowGameStatToPlayers(h *Hub) {
	playerList := make([]PlayerStat, len(h.Clients))
	for cl, uid := range h.Clients {
		playerList[cl.Order].UID = uid
		playerList[cl.Order].Name = cl.Name
		playerList[cl.Order].Score = cl.Score
	}
	gameStat := GameStat{
		CurrentlyDrawing: h.CurrentlyDrawing,
		Answer:           h.Words[0],
		Players:          playerList,
	}

	jsonBytes, err := json.Marshal(gameStat)
	if err != nil {
		log.Println(err)
	}

	for cl := range h.Clients {
		select {
		case cl.Send <- jsonBytes:
		default:
			close(cl.Send)
			delete(h.Clients, cl)
		}
	}
}
