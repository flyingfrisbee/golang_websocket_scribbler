package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

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
	router.HandleFunc("/{roomName}/{playerName}", CreateRoom).Methods("GET")

	err := http.ListenAndServe("localhost:8080", router)
	if err != nil {
		log.Println(err)
	}
}
