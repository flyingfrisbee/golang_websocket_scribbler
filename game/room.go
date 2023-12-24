package game

type Room struct {
	// All players
	Players map[int]*Player
	// Inbound messages from the players.
	MsgFromPlayer chan []byte
	// Register requests from players.
	Register chan *Player
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
}

func (r *Room) Run() {
	defer func() {
		close(r.MsgFromPlayer)
		close(r.Unregister)
		close(r.Register)
		HubObj.removeRoomByID(r.ID)
	}()

	for {
		select {
		case player := <-r.Register:
			if len(r.Players) >= 4 {
				player.AckChan <- false
				break
			}
			r.Players[player.ID] = player
			r.TurnOrder = append(r.TurnOrder, player.ID)
			player.AckChan <- true
		case player := <-r.Unregister:
			if _, ok := r.Players[player.ID]; ok {
				shouldCloseRoom := r.unregisterPlayer(player)
				if shouldCloseRoom {
					return
				}
			}
		case message := <-r.MsgFromPlayer:
			for _, player := range r.Players {
				select {
				case player.MsgToPlayer <- message:
				default:
					// Might happen when for some reason the client cannot
					// process the message, since MsgToPlayer is buffered
					// channel that means something wrong occured on player
					shouldCloseRoom := r.unregisterPlayer(player)
					if shouldCloseRoom {
						return
					}
				}
			}
		}
	}
}

func (r *Room) unregisterPlayer(p *Player) bool {
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

	return len(r.Players) == 0
}

func CreateRoom(roomID string) *Room {
	return &Room{
		Players:        make(map[int]*Player),
		MsgFromPlayer:  make(chan []byte),
		Register:       make(chan *Player),
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
