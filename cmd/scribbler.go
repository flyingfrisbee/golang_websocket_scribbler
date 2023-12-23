package cmd

import (
	"GithubRepository/golang_websocket_scribbler/api"
	db "GithubRepository/golang_websocket_scribbler/database"
	env "GithubRepository/golang_websocket_scribbler/environment"

	"log"
	"net/http"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func Run() {
	// required endpoints:
	// create room (generate random uuid as the key of the map)
	// join room (find the map based on the key provided)
	env.GenerateEnvVar()
	db.SetupDB()
	defer db.DB.CloseConnection()
	api.SetupAPI()
}
