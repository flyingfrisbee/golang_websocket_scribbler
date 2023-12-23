package api

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

type genericResponse struct {
	Data       interface{} `json:"data"`
	Message    string      `json:"message"`
	StatusCode int         `json:"status_code"`
}

func SetupAPI() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	defineUserRoute(r)
	defineRoomRoute(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	r.Run(addr)
}

func writeResponse(
	c *gin.Context,
	data interface{},
	msg string,
	code int,
) {
	resp := genericResponse{
		Data:       data,
		Message:    msg,
		StatusCode: code,
	}
	c.JSON(code, resp)
}
