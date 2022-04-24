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
	RoomName         string
	Broadcast        chan []byte
	Register         chan *Client
	Unregister       chan *Client
	Clients          map[*Client]int
	Words            []string
	TurnNumber       int //-> tiap ganti giliran increment++
	CurrentlyDrawing int //(uid) -> dapet dari clients[turnNumber % len(clients)]
}

type GameStat struct {
	CurrentlyDrawing int          `json:"currently_drawing"`
	Answer           string       `json:"answer"`
	Players          []PlayerStat `json:"players"`
}

type PlayerStat struct {
	UID         int    `json:"uid"`
	Name        string `json:"name"`
	Score       int    `json:"score"`
	HasAnswered bool   `json:"has_answered"`
}

func newHub(roomName string) *Hub {
	words := []string{}
	for k := range ListOfItems {
		words = append(words, k)
	}

	return &Hub{
		RoomName:   roomName,
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]int),
		Words:      words,
		TurnNumber: 1,
	}
}

func (h *Hub) run() {
	defer func() {
		h.Lock()
		if len(h.Clients) != 0 {
			for cl := range h.Clients {
				delete(h.Clients, cl)
				close(cl.Send)
			}
		}
		// TODO: clear movement cache mongodb
		DeleteMovementCollection(h.RoomName)
		h.Unlock()

		close(h.Broadcast)
		close(h.Register)
		// close(h.Unregister)
		DeleteHub(h)
	}()

	h.startChannelListener()
}

func (h *Hub) startChannelListener() {
	defer h.Unlock()

	for {
		select {
		case client := <-h.Register:

			h.Lock()
			h.Clients[client] = client.UID
			if len(h.Clients) == 1 {
				h.CurrentlyDrawing = client.UID
			}

			// TODO: check if entry exists in mongodb, if exist: fetch cache movement from mongodb, send through the client channel cl.Send <- data
			if CheckIfMovementCacheExist(h.RoomName) {
				for _, v := range GetMovement(h.RoomName) {
					client.Send <- []byte(v.Text)
				}
			}
			h.Unlock()

			ShowGameStatToPlayers(h)

		case client := <-h.Unregister:

			h.Lock()
			//user forcefully exit the game while in turn, resulting in room closing
			if h.Clients[client] == h.CurrentlyDrawing {
				return
			}

			orderPlaceholder := client.Order

			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}

			// no more player left, close the room
			if len(h.Clients) == 0 {
				return
			}

			for cl := range h.Clients {
				if cl.Order > orderPlaceholder {
					cl.Order--
				}
			}
			h.Unlock()

			ShowGameStatToPlayers(h)

		case message := <-h.Broadcast:

			switch message[0] {

			case '0':
				msg := strings.Split(string(message), ";")
				if len(msg) != 2 {
					return
				}

				receivedUID, err := strconv.Atoi(msg[1])
				//non authorized format, close the room
				if err != nil {
					return
				}

				// TODO: clear mongodb entries for movement cache
				h.Lock()
				DeleteMovementCollection(h.RoomName)
				h.Unlock()

				h.RLock()
				for cl, uid := range h.Clients {
					if receivedUID == uid {
						continue
					}
					select {
					case cl.Send <- message:
					default:
						close(cl.Send)
						delete(h.Clients, cl)
					}
				}
				h.RUnlock()

			case '1':

				msg := strings.Split(string(message), ";")
				if len(msg) != 5 {
					return
				}
				width := msg[2][2:]
				height := msg[3][2:]
				movementList := msg[4]

				receivedUID, err := strconv.Atoi(msg[1])
				//non authorized format, close the room
				if err != nil {
					return
				}

				// TODO: write movement to mongodb
				h.Lock()
				CreateNewMovementEntry(h.RoomName, string(message))
				h.Unlock()

				h.RLock()
				for cl, uid := range h.Clients {
					if receivedUID == uid {
						continue
					}
					select {
					case cl.Send <- []byte(fmt.Sprintf("[%s;%s;%s", width, height, movementList)):
					default:
						close(cl.Send)
						delete(h.Clients, cl)
					}
				}
				h.RUnlock()

			case '2':
				h.Lock()

				for cl, uid := range h.Clients {
					if uid == h.CurrentlyDrawing {
						cl.Score += 2
						break
					}
				}

				h.Words = h.Words[1:]
				if len(h.Words) == 0 {
					for client := range h.Clients {
						select {
						case client.Send <- []byte{'4'}:
						default:
							close(client.Send)
							delete(h.Clients, client)
						}
					}
				} else {
					for k, v := range h.Clients {
						if k.Order == h.TurnNumber%len(h.Clients) {
							h.CurrentlyDrawing = v
						}

						k.HasAnswered = false
					}

					h.TurnNumber++
				}

				h.Unlock()
				if len(h.Words) == 0 {
					ShowEndGameStatToPlayers(h)
				} else {
					ShowGameStatToPlayers(h)
				}

			case '3':

				h.Lock()
				msg := strings.Split(string(message), ";")
				//non authorized msg, close the room
				if len(msg) != 3 {
					return
				}
				receivedUID, err := strconv.Atoi(msg[1])
				//non authorized format, close the room
				if err != nil {
					return
				}
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
				h.Unlock()
				ShowGameStatToPlayers(h)

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
		}
	}
}

func ShowGameStatToPlayers(h *Hub) {
	h.RLock()
	defer h.RUnlock()

	playerList := make([]PlayerStat, len(h.Clients))

	for cl, uid := range h.Clients {
		playerList[cl.Order].UID = uid
		playerList[cl.Order].Name = cl.Name
		playerList[cl.Order].Score = cl.Score
		playerList[cl.Order].HasAnswered = cl.HasAnswered
	}
	gameStat := GameStat{
		CurrentlyDrawing: h.CurrentlyDrawing,
		Answer:           h.Words[0],
		Players:          playerList,
	}

	jsonBytes, err := json.Marshal(gameStat)
	if err != nil {
		log.Println(err)
		return
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

func ShowEndGameStatToPlayers(h *Hub) {
	h.RLock()
	defer h.RUnlock()

	playerList := make([]PlayerStat, len(h.Clients))

	for cl, uid := range h.Clients {
		playerList[cl.Order].UID = uid
		playerList[cl.Order].Name = cl.Name
		playerList[cl.Order].Score = cl.Score
		playerList[cl.Order].HasAnswered = cl.HasAnswered
	}
	gameStat := GameStat{
		CurrentlyDrawing: h.CurrentlyDrawing,
		Answer:           "Finished!",
		Players:          playerList,
	}

	jsonBytes, err := json.Marshal(gameStat)
	if err != nil {
		log.Println(err)
		return
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
