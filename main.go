package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
)

type AppVersion struct {
	MandatoryVersion int `json:"mandatory_version"`
	OptionalVersion  int `json:"optional_version"`
}

type Resp struct {
	Success bool `json:"success"`
}

var (
	AppVersionMandatory = 2
	AppVersionOptional  = 1
	mutex               = sync.RWMutex{}
	hubs                = map[string]*Hub{}
	UID                 = 1
)

func GetUID() int {
	defer func() {
		UID++
	}()

	return UID
}

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

func GetAppVersion(w http.ResponseWriter, r *http.Request) {
	SendAppVersionResponse(w, r, http.StatusOK)
}

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	mutex.RLock()
	_, ok := hubs[params["roomName"]]
	mutex.RUnlock()

	if ok {
		SendResponse(w, r, false, http.StatusOK)
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
		SendResponse(w, r, false, http.StatusOK)
		return
	}

	length := len(hubs[params["roomName"]].Clients)
	if length >= 4 {
		SendResponse(w, r, false, http.StatusOK)
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
		hub := newHub(params["roomName"])
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
	defer func() {
		CloseConnectionMongoDB()
	}()
	CreateConnectionToMongoDB()
	router := mux.NewRouter()
	router.HandleFunc("/version", GetAppVersion).Methods("GET")
	router.HandleFunc("/createroom/{roomName}", CreateRoom).Methods("GET")
	router.HandleFunc("/joinroom/{roomName}", JoinRoom).Methods("GET")
	router.HandleFunc("/ws/{roomName}/{playerName}", CreateWebsocketConnection).Methods("GET")

	port := os.Getenv("PORT")

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		log.Println(err)
	}
}
