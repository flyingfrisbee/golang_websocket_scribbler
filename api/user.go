package api

import (
	"GithubRepository/golang_websocket_scribbler/database"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	Username string `json:"username"`
}

type createUserResponse struct {
	UserID int `json:"user_id"`
}

func createUser(c *gin.Context) {
	jsonBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		writeResponse(c, nil, "failed to read request body", http.StatusBadRequest)
		return
	}
	defer c.Request.Body.Close()

	var req createUserRequest
	err = json.Unmarshal(jsonBytes, &req)
	if err != nil {
		writeResponse(c, nil, "failed to parse JSON request", http.StatusBadRequest)
		return
	}

	id, err := database.DB.CreateUser(req.Username)
	if err != nil {
		errMsg := fmt.Sprintf("error occured when creating user: %s", err.Error())
		writeResponse(c, nil, errMsg, http.StatusInternalServerError)
		return
	}

	resp := createUserResponse{
		UserID: id,
	}
	writeResponse(c, &resp, "success creating user", http.StatusOK)
}

func defineUserRoute(r *gin.Engine) {
	userRoute := r.Group("/users")
	userRoute.POST("/register", createUser)
}
