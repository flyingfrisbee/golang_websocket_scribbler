package api

import (
	"GithubRepository/golang_websocket_scribbler/game"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	roomID string
)

func createRoom(c *gin.Context) {
	_, err := strconv.Atoi(c.Request.URL.Query().Get("userId"))
	if err != nil {
		writeResponse(c, nil, "cannot convert userId to integer", http.StatusBadRequest)
		return
	}

	room := game.HubObj.AddRoomToHub()
	go room.Run()
	a, _ := c.Writer.(http.ResponseWriter)
	game.ServeWs(room, a, c.Request)

	// TODO: Delete this later, for testing purpose only
	roomID = room.ID
}

func joinRoom(c *gin.Context) {
	// roomID := c.Param("roomId")
	room := game.HubObj.FindRoomByID(roomID)
	a, _ := c.Writer.(http.ResponseWriter)
	game.ServeWs(room, a, c.Request)
}

func defineRoomRoute(r *gin.Engine) {
	roomRoute := r.Group("/rooms")
	roomRoute.GET("/create", createRoom)
	roomRoute.GET("/join/:roomId", joinRoom)
}
