package main

import (
	"github.com/brkr/goreal"
	"log"
	"net/http"
	"time"
)

func main() {
	gameServer := goreal.NewGameServer(2020)
	log.Println(gameServer.Port)
	lobby := NewLobby()
	gameServer.InitClient = func(w http.ResponseWriter, r *http.Request, client *goreal.Client) {
		log.Println("new user handler")
		keys, ok := r.URL.Query()["name"]
		if ok {
			log.Println( keys)
			gameServer.JoinRoom("lobby", client)
		}

	}

	gameServer.RegisterRoom("lobby", lobby)
	gameServer.Start()

}



type Lobby struct {
	goreal.Room

}

func (l *Lobby) OnJoinRequest(client *goreal.Client) bool {
	log.Println("lobby onJoinRequest")

	if len(l.Clients) >= 2 {
		log.Println("room is full")
		return false
	}

	return true
}

func (l *Lobby) OnInit() {


}

func (l *Lobby) OnMessage(client *goreal.Client, message []byte) {
	log.Printf("Lobby: Message : %s", string(message))
	time.Sleep(4 * time.Second)
	client.SendMessageStr("{\"users\": [{\"name\": \"cl1\",\"id\": 1},{\"name\": \"cl1\",\"id\": 2}]}")
}

func (l *Lobby) OnClientJoin(client *goreal.Client) {
	log.Println("lobby onClientJoin")


}

func (l *Lobby) OnLeave(client *goreal.Client) {
	client.SendMessageStr("You're kicked from room.")
}

func (l *Lobby) OnUpdate(delta float64) {
	//log.Println("update game simulation delta time:  ", delta)
}

func NewLobby() *Lobby {
	return &Lobby{}
}

