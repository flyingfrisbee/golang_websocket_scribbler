package api

import (
	db "GithubRepository/golang_websocket_scribbler/database"
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

	id, err := db.DB.CreateUser(req.Username)
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

type deleteUserRequest struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}

func deleteUser(c *gin.Context) {
	jsonBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		writeResponse(c, nil, "failed to read request body", http.StatusBadRequest)
		return
	}
	defer c.Request.Body.Close()

	var req deleteUserRequest
	err = json.Unmarshal(jsonBytes, &req)
	if err != nil {
		writeResponse(c, nil, "failed to parse JSON request", http.StatusBadRequest)
		return
	}

	err = db.DB.DeleteUser(req.UserID, req.Username)
	if err != nil {
		writeResponse(c, nil, err.Error(), http.StatusBadRequest)
		return
	}

	writeResponse(c, nil, "success deleted user", http.StatusOK)
}

func deleteUserPublic(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "index.html")
}

func defineUserRoute(r *gin.Engine) {
	userRoute := r.Group("/users")
	userRoute.POST("/register", createUser)
	userRoute.POST("/delete", deleteUser)
	userRoute.GET("/delete-public", deleteUserPublic)
}
