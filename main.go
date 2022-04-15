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
	if _, ok := hubs[params["roomName"]]; !ok {
		hub := newHub()
		go hub.run()

		mutex.Lock()
		hubs[params["roomName"]] = hub
		mutex.Unlock()

		serveWs(hub, w, r)
		return
	}

	serveWs(hubs[params["roomName"]], w, r) // TODO: later switch to /join/{roomName}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/create/{roomName}", CreateRoom).Methods("GET")

	// router.HandleFunc("join", func(w http.ResponseWriter, r *http.Request) {
	//	do something
	// }).Methods("GET")

	err := http.ListenAndServe("localhost:8080", router)
	if err != nil {
		log.Println(err)
	}
}

// buat handlefunc buat konek dari android -> server
// ada 2 handlefunc -> bisa create / bisa join
// -create handlefunc- cek di map[string]*Hub kalo !exist -> create hub baru ke map (kasi mutex.Lock), balikin uid ke user (json), nyalain client listener,
// type Hub struct {
//  sync.RWMutex
// 	clients []*Client
// 	broadcast chan []byte
// 	register chan *Client
// 	unregister chan *Client
// 	turnNumber int -> tiap ganti giliran increment++
// 	currentlyDrawing int64(uid) -> dapet dari clients[turnNumber % len(clients)]
// }

// type Client struct {
// 	uid int64
// 	hub *Hub
// 	conn *websocket.Conn
// 	send chan []byte
// }
