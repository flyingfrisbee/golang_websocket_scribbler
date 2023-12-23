package api

import (
	"GithubRepository/golang_websocket_scribbler/game"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createRoomRequest struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}

func createRoom(c *gin.Context) {
	// create room in hub
	// create room
	// create client

	// a, _ := c.Writer.(http.ResponseWriter)
	// game.ServeWs(hub, a, c.Request)
}

func joinRoom(c *gin.Context) {
	// roomID := c.Param("room_id")

}

func defineRoomRoute(r *gin.Engine) {
	hub := game.CreateRoom()
	go hub.Run()

	roomRoute := r.Group("/rooms")
	// roomRoute.POST("/create", createRoom)
	// roomRoute.GET("/join/:room_id", joinRoom)

	roomRoute.GET("/", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	}))
	roomRoute.GET("/ws", func(c *gin.Context) {
		a, _ := c.Writer.(http.ResponseWriter)
		game.ServeWs(hub, a, c.Request)
	})
}
