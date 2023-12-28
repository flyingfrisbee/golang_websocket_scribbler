package game

type ActionCode int

const (
	UpdateGameInfo ActionCode = iota
	Drawing
	ClearDrawing
	Answer
	GameFinished
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
