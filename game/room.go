package game

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type NextTurnTrigger int

const (
	CurrentTurnPlayerLeft NextTurnTrigger = iota
	AllPlayersHaveSubmittedAnswer
)

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
	// Drawing cache
	Cache [][]interface{}
	// Number of people that has submitted their answer
	AnswersCount int
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
			roomIsFullyPopulated := len(r.Players) >= 4
			if roomIsFullyPopulated {
				player.AckChan <- false
				break
			}

			_, duplicatePlayer := r.Players[player.ID]
			if duplicatePlayer {
				player.AckChan <- false
				break
			}

			gameConcludes := len(r.Words) == 0
			if gameConcludes {
				// do nothing, wait until all players left the room
				player.AckChan <- false
				break
			}

			r.Players[player.ID] = player
			r.TurnOrder = append(r.TurnOrder, player.ID)
			player.AckChan <- true
			// GameInfo
			msg := r.generateGameInfo()
			_, err := r.sendMessageToPlayers(&msg)
			if err != nil {
				log.Println(err)
				break
			}
			// Sync drawing
			if len(r.Cache) == 0 {
				break
			}
			msg.Code = Drawing
			msg.Data = r.Cache
			jsonBytes, err := json.Marshal(&msg)
			if err != nil {
				log.Println(err)
				break
			}
			player.MsgToPlayer <- jsonBytes
		case player := <-r.Unregister:
			_, ok := r.Players[player.ID]
			if !ok {
				break
			}

			proceedNextTurn := r.TurnOrder[r.CurrentTurnIdx] == player.ID
			shouldCloseRoom := r.unregisterPlayer(player)
			if shouldCloseRoom {
				return
			}

			gameConcludes := len(r.Words) == 0
			if gameConcludes {
				// do nothing, wait until all players left the room
				break
			}

			msg := r.generateGameInfo()
			_, err := r.sendMessageToPlayers(&msg)
			if err != nil {
				log.Println(err)
				break
			}

			if !proceedNextTurn {
				break
			}

			gameConcludes = r.nextTurn(CurrentTurnPlayerLeft)
			if gameConcludes {
				msg := userMessage{
					Code: GameFinished,
					Data: nil,
				}
				_, err := r.sendMessageToPlayers(&msg)
				if err != nil {
					log.Println(err)
					break
				}
				break
			}

			msg = r.generateGameInfo()
			_, err = r.sendMessageToPlayers(&msg)
			if err != nil {
				log.Println(err)
				break
			}
		case message := <-r.MsgFromPlayer:
			gameConcludes := len(r.Words) == 0
			if gameConcludes {
				// do nothing, wait until all players left the room
				break
			}

			var msg userMessage
			err := json.Unmarshal(message, &msg)
			if err != nil {
				log.Println(err)
				break
			}

			err = r.handleIncomingMessage(&msg)
			if err != nil {
				log.Println(err)
				break
			}

			_, err = r.sendMessageToPlayers(&msg)
			if err != nil {
				log.Println(err)
				break
			}

			playerNotSolo := len(r.Players) > 1
			// Player that is drawing cannot give answer, hence - 1
			everyoneAlreadyAnswered := r.AnswersCount == (len(r.Players) - 1)
			proceedNextTurn := playerNotSolo && everyoneAlreadyAnswered
			if !proceedNextTurn {
				break
			}

			gameConcludes = r.nextTurn(AllPlayersHaveSubmittedAnswer)
			if gameConcludes {
				msg := userMessage{
					Code: GameFinished,
					Data: nil,
				}
				_, err := r.sendMessageToPlayers(&msg)
				if err != nil {
					log.Println(err)
					break
				}
				break
			}

			msg = r.generateGameInfo()
			_, err = r.sendMessageToPlayers(&msg)
			if err != nil {
				log.Println(err)
				break
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

	r.TurnOrder = append(r.TurnOrder[:i], r.TurnOrder[i+1:]...)
	if r.CurrentTurnIdx > len(r.TurnOrder)-1 {
		r.CurrentTurnIdx = 0
	}

	return len(r.Players) == 0
}

func (r *Room) nextTurn(cause NextTurnTrigger) bool {
	r.Cache = nil
	r.Words = r.Words[1:]
	r.AnswersCount = 0

	for _, id := range r.TurnOrder {
		r.Players[id].HasAnswered = false
	}

	switch cause {
	case AllPlayersHaveSubmittedAnswer:
		r.CurrentTurnIdx++
		if r.CurrentTurnIdx > len(r.TurnOrder)-1 {
			r.CurrentTurnIdx = 0
		}
	}

	return len(r.Words) == 0
}

func (r *Room) sendMessageToPlayers(msg *userMessage) (bool, error) {
	if msg.Code == Answer && msg.Data == nil {
		return false, nil
	}

	jsonBytes, err := json.Marshal(&msg)
	if err != nil {
		return false, err
	}

	for _, id := range r.TurnOrder {
		player := r.Players[id]
		select {
		case player.MsgToPlayer <- jsonBytes:
		default:
			// Might happen when for some reason the client cannot
			// process the message, since MsgToPlayer is buffered
			// channel that means something wrong occured on player
			shouldCloseRoom := r.unregisterPlayer(player)
			if shouldCloseRoom {
				return true, nil
			}
		}
	}

	return false, nil
}

// Update game data based on the action code and modify message accordingly
func (r *Room) handleIncomingMessage(msg *userMessage) error {
	switch msg.Code {
	case Drawing:
		coords, ok := msg.Data.([]interface{})
		if !ok {
			return fmt.Errorf("failed when converting message from player: drawing")
		}

		r.Cache = append(r.Cache, coords)
		msg.Data = [][]interface{}{coords}
	case ClearDrawing:
		r.Cache = nil
	case Answer:
		ans, ok := msg.Data.(string)
		if !ok {
			return fmt.Errorf("failed when converting message from player: answer")
		}

		player := r.Players[msg.SenderID]

		if len(r.Cache) != 0 || r.AnswersCount != 0 {
			// i'm thinking about a case where a player left and then new turn
			// has begun, but this person already answered the previous turn's drawing
			playerGuessedCorrectly := strings.EqualFold(ans, r.Words[0])
			if playerGuessedCorrectly {
				player.Score++
			}
			player.HasAnswered = true
			r.AnswersCount++
			msg.Data = player.mapToPlayerInfo()
			break
		}
		msg.Data = nil
	}

	return nil
}

func (r *Room) generateGameInfo() userMessage {
	playersInfo := make([]playerInfo, len(r.TurnOrder))
	for idx, id := range r.TurnOrder {
		playersInfo[idx] = r.Players[id].mapToPlayerInfo()
	}
	gameInfo := gameInfo{
		RoomID:         r.ID,
		Players:        playersInfo,
		CurrentTurnIdx: r.CurrentTurnIdx,
		Word:           r.Words[0],
	}

	msg := userMessage{
		Code: UpdateGameInfo,
		Data: gameInfo,
	}

	return msg
}

func (r *Room) findWinner() userMessage {
	var winner *Player
	max := 0
	for _, id := range r.TurnOrder {
		p := r.Players[id]
		if max < p.Score {
			max = p.Score
			winner = p
		}
	}

	return userMessage{
		Code: GameFinished,
		Data: winner.mapToPlayerInfo(),
	}
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
