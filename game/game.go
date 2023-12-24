package game

type ActionCode int

const (
	UpdateGameInfo ActionCode = iota // sendData: gameInfo
	Drawing                          // receivedData: []drawing, sendData: []drawing
	ClearDrawing                     // receivedData: nil, sendData: nil
	Answer                           // receivedData: string, sendData: player
	SyncDrawing                      // sendData: []drawing from cache
	GameFinished                     // sendData: player that wins
)

type userMessage struct {
	Code     ActionCode  `json:"code"`
	SenderID int         `json:"sender_id"`
	Data     interface{} `json:"data"`
}

type drawingCoordinate struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type gameInfo struct {
	RoomID         string       `json:"room_id"`
	Players        []playerInfo `json:"players"`
	CurrentTurnIdx int          `json:"current_turn_index"`
	Word           string       `json:"word"`
}

type playerInfo struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Score        int    `json:"score"`
	ScreenWidth  int    `json:"screen_width"`
	ScreenHeight int    `json:"screen_height"`
	HasAnswered  bool   `json:"has_answered"`
}

func generateGameInfo() {

}

// USER CAN DRAW
// get coordinates of the lines, write to cache, broadcast msg

// USER CAN CLEAR THE DRAWING
// delete cache, broadcast msg

// USER CAN ANSWER
// validate answer, grant score, update gameinfo, potentially change turn

// USER CAN JOIN ROOM
// send draw cache, update gameinfo

// USER CAN LEAVE ROOM
// update gameinfo, potentially change turn

// NO MORE DRAWING OBJECTS
// game finished, declare winner, close room