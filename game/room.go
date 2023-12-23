package game

type Room struct {
	// All players
	Players map[int]*Player
	// Inbound messages from the players.
	MsgFromPlayer chan []byte
	// Register requests from the clients.
	Register chan *Player
	// Unregister requests from players.
	Unregister chan *Player
	// Index of currently drawing player from turnOrder
	CurrentTurnIdx int
	// Player ids based on the order of who entered the room
	TurnOrder []int
	// Name of the object to be drawn
	Words []string
}

func (r *Room) Run() {
	for {
		select {
		case player := <-r.Register:
			r.Players[player.ID] = player
		case player := <-r.Unregister:
			if _, ok := r.Players[player.ID]; ok {
				delete(r.Players, player.ID)
				close(player.MsgToPlayer)
			}
		case message := <-r.MsgFromPlayer:
			for _, player := range r.Players {
				select {
				case player.MsgToPlayer <- message:
				default:
					close(player.MsgToPlayer)
					delete(r.Players, player.ID)
				}
			}
		}
	}
}

func CreateRoom() *Room {
	return &Room{
		Players:        make(map[int]*Player),
		MsgFromPlayer:  make(chan []byte),
		Register:       make(chan *Player),
		Unregister:     make(chan *Player),
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
