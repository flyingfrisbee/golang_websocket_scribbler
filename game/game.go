package game

type ActionCode int

const (
	UpdateGameInfo ActionCode = iota // sendData: gameInfo
	Drawing                          // receivedData: []drawing, sendData: []drawing
	ClearDrawing                     // receivedData: nil, sendData: nil
	Answer                           // receivedData: string, sendData: player
	GameFinished                     // sendData: nil, check for winner from client
)

type userMessage struct {
	Code     ActionCode  `json:"code"`
	SenderID int         `json:"sender_id"`
	Data     interface{} `json:"data"`
}

type drawingCoordinate struct {
	StartX float32 `json:"start_x"`
	StartY float32 `json:"start_y"`
	EndX   float32 `json:"end_x"`
	EndY   float32 `json:"end_y"`
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
