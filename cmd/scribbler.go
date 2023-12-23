package cmd

import (
	"GithubRepository/golang_websocket_scribbler/api"
	db "GithubRepository/golang_websocket_scribbler/database"
	env "GithubRepository/golang_websocket_scribbler/environment"
)

func Run() {
	// required endpoints:
	// create room (generate random uuid as the key of the map)
	// join room (find the map based on the key provided)
	env.GenerateEnvVar()
	db.SetupDB()
	defer db.DB.CloseConnection()
	api.SetupAPI()
}
