package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
)

type Resp struct {
	Success bool
}

var (
	mutex = sync.RWMutex{}
	hubs  = map[string]*Hub{}
)

func DeleteHub(hub *Hub) {
	mutex.Lock()
	defer mutex.Unlock()
	for k, v := range hubs {
		if v == hub {
			delete(hubs, k)
			return
		}
	}
}

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	mutex.RLock()
	_, ok := hubs[params["roomName"]]
	mutex.RUnlock()

	if ok {
		SendResponse(w, r, false, http.StatusNotAcceptable)
		return
	}

	SendResponse(w, r, true, http.StatusOK)
}

func JoinRoom(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	mutex.RLock()
	defer mutex.RUnlock()

	_, ok := hubs[params["roomName"]]
	if !ok {
		SendResponse(w, r, false, http.StatusNotAcceptable)
		return
	}

	length := len(hubs[params["roomName"]].Clients)
	if length >= 4 {
		SendResponse(w, r, false, http.StatusNotAcceptable)
		return
	}

	SendResponse(w, r, true, http.StatusOK)
}

func CreateWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	mutex.RLock()
	_, ok := hubs[params["roomName"]]
	mutex.RUnlock()

	if !ok {
		hub := newHub()
		go hub.run()

		mutex.Lock()
		hubs[params["roomName"]] = hub
		mutex.Unlock()

		serveWs(hub, w, r, params["playerName"])
		return
	}

	mutex.RLock()
	length := len(hubs[params["roomName"]].Clients)
	mutex.RUnlock()

	if length < 4 {
		serveWs(hubs[params["roomName"]], w, r, params["playerName"])
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/createroom/{roomName}", CreateRoom).Methods("GET")
	router.HandleFunc("/joinroom/{roomName}", JoinRoom).Methods("GET")
	router.HandleFunc("/ws/{roomName}/{playerName}", CreateWebsocketConnection).Methods("GET")

	port := os.Getenv("PORT")

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		log.Println(err)
	}
}
