package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/brkr/goreal"
)

func main() {
	gameServer := goreal.NewGameServer(2020)
	log.Println(gameServer.Port)
	lobby := NewLobby()
	gameServer.InitClient = func(w http.ResponseWriter, r *http.Request, client *goreal.Client) {
		log.Println("new user handler")
		keys, ok := r.URL.Query()["name"]
		if ok {
			log.Println(keys)
			client.Put("name", keys)
			gameServer.JoinRoom("lobby", client)
		}

	}

	gameServer.RegisterRoom("lobby", lobby)
	gameServer.Start()

}

type Lobby struct {
	goreal.Room
	Users map[*goreal.Client]User
}

type User struct {
	Status string
	Name   string
	Client *goreal.Client
}

type OnlineResponse struct {
	Users []User
}

func (l *Lobby) OnJoinRequest(client *goreal.Client) bool {
	log.Println("lobby onJoinRequest")
	return true
}

func (l *Lobby) OnInit() {
	l.Config.SimulationTick = 1

}

func (l *Lobby) OnMessage(client *goreal.Client, message []byte) {
	log.Printf("Lobby: Message : %s", string(message))
	// time.Sleep(4 * time.Second)
	// client.SendMessageStr("{\"users\": [{\"name\": \"cl1\",\"id\": 1},{\"name\": \"cl1\",\"id\": 2}]}")

}

func (l *Lobby) OnClientJoin(client *goreal.Client) {
	log.Println("lobby onClientJoin")

	username := client.Get("name")
	log.Printf("asdasd %s", username)

	name, err := username.(string)
	if err {

	}

	log.Printf(name)

	user := &User{Name: name, Status: "online", Client: client}
	log.Printf("client %s connected.", name)
	l.Users[client] = *user

	log.Println(l.Users[client])

}

func (l *Lobby) OnLeave(client *goreal.Client) {
	client.SendMessageStr("You're kicked from room.")
}

func (l *Lobby) OnUpdate(delta float64) {
	//log.Println("update game simulation delta time:  ", delta)
	log.Printf("sending status. User count %d", len(l.Users))

	response := &OnlineResponse{Users: make([]User, len(l.Users))}

	i := 0
	for _, v := range l.Users {
		response.Users[i] = v
	}

	jsonString, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
	}

	l.BroadcastMessageByte(jsonString)
}

func NewLobby() *Lobby {
	return &Lobby{Users: make(map[*goreal.Client]User)}
}
