package cmd

import (
	"GithubRepository/golang_websocket_scribbler/api"
	db "GithubRepository/golang_websocket_scribbler/database"
	env "GithubRepository/golang_websocket_scribbler/environment"
)

func Run() {
	env.GenerateEnvVar()
	db.SetupDB()
	defer db.DB.CloseConnection()
	api.SetupAPI()
}
