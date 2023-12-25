package api

import (
	"GithubRepository/golang_websocket_scribbler/game"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func listRooms(c *gin.Context) {
	writeResponse(c, game.HubObj.ListRooms(), "success fetch all room ids", http.StatusOK)
}

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
}

func joinRoom(c *gin.Context) {
	_, err := strconv.Atoi(c.Request.URL.Query().Get("userId"))
	if err != nil {
		writeResponse(c, nil, "cannot convert userId to integer", http.StatusBadRequest)
		return
	}

	roomID := c.Param("roomId")
	room := game.HubObj.FindRoomByID(roomID)
	if room == nil {
		writeResponse(c, nil, "room not found", http.StatusBadRequest)
		return
	}
	a, _ := c.Writer.(http.ResponseWriter)
	game.ServeWs(room, a, c.Request)
}

func defineRoomRoute(r *gin.Engine) {
	roomRoute := r.Group("/rooms")
	roomRoute.GET("/", listRooms)
	roomRoute.GET("/create", createRoom)
	roomRoute.GET("/join/:roomId", joinRoom)
}
