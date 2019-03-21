package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brkr/goreal"
	"log"
	"net/http"
	"reflect"
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
			log.Println(keys[0])
			client.Put("name", keys[0])
			client.Put("profilePic", "https://ptetutorials.com/images/user-profile.png")
			gameServer.JoinRoom("lobby", client)
		}

	}

	gameServer.RegisterRoom("lobby", lobby)
	gameServer.Start()
}

type Lobby struct {
	goreal.Room
	Users map[*goreal.Client]User
	id    int32
}

type User struct {
	Status  string         `json:"status"`
	Name    string         `json:"name"`
	Id      string         `json:"id"`
	Profile string         `json:"profilePicture"`
	Client  *goreal.Client `json:"-"`
}

type OnlineResponse struct {
	Users []User `json:"users"`
}

type JsonData struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}
type Message struct {
	FromId       string `json:"from_id"`
	ToId         string `json:"to_id"`
	Name         string `json:"name"`
	Message      string `json:"message"`
	ProfilePhoto string `json:"profilePicture"`
	DateStr      string `json:"date_str"`
}

func (l *Lobby) OnJoinRequest(client *goreal.Client) bool {
	log.Println("lobby onJoinRequest")
	if name, ok := client.Get("name").(string); ok {
		if name == "berker" {
			client.SendMessageStr("atildiniz..")
			client.CloseConnection()
			return false
		}
	}
	return true
}

func (l *Lobby) OnInit() {
	// saniyede bir kez
	l.Config.SimulationTick = 1
}

func (l *Lobby) OnMessage(client *goreal.Client, message []byte) {
	log.Printf("Lobby: Message : %s", string(message))

	dict := &JsonData{}

	if err := json.Unmarshal(message, dict); err != nil {
		log.Println("json data bulunamadi")
	}

	log.Println(dict)

	if dict.Event == "private_message" {

		log.Println("private message")

		fmt.Println(reflect.TypeOf(dict.Data).String())

		if m, ok := dict.Data.(map[string]interface{}); ok {
			l.sendPrivateMessage(client, m)
		}
	}
}

func (l *Lobby) OnClientJoin(client *goreal.Client) {
	log.Println("lobby onClientJoin")

	username := client.Get("name")
	log.Printf("Join %s", username)
	l.id++

	sid := fmt.Sprintf("id_%d", l.id)

	user := &User{Name: "", Status: "online", Client: client, Id: sid}
	if name, ok := username.(string); ok {
		log.Println("casting")
		user.Name = name
	}

	if profilePic, ok := client.Get("profilePic").(string); ok {
		user.Profile = profilePic
	}

	l.Users[client] = *user

	log.Println(l.Users[client])
	// first
	l.sendClientOwnInfo(client, *user)

	l.sendOnlineToAllUser()
}

func (l *Lobby) sendClientOwnInfo(client *goreal.Client, user User) {
	jsonData := &JsonData{Data: user, Event: "userinfo"}
	if jsonString, err := json.Marshal(jsonData); err == nil {
		client.SendMessage(jsonString)
	} else {
		log.Println(err)
	}

}

func (l *Lobby) sendOnlineToAllUser() {

	response := &OnlineResponse{Users: make([]User, len(l.Users))}

	i := 0
	for _, v := range l.Users {
		response.Users[i] = v
		i++
	}

	jsonString, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
	}

	l.BroadcastMessageByte(jsonString)

}

func (l *Lobby) OnLeave(client *goreal.Client) {
	user := l.Users[client]
	log.Printf("%s kullanici odadan ayrildi", user.Name)
	delete(l.Users, client)

	l.sendOnlineToAllUser()
}

func (l *Lobby) OnUpdate(delta float64) {
	//log.Println("update game simulation delta time:  ", delta)
	// log.Printf("sending status. User count %d", len(l.Users))
}

// private message
func (l *Lobby) sendPrivateMessage(client *goreal.Client, data map[string]interface{}) {

	from, err := l.findUserByClient(client)

	if err != nil {
		log.Println(err)
		return
	}

	for _, to := range l.Users {
		if to.Id == data["to"] {
			log.Println("send to oooo")
			msg := &Message{FromId: from.Id, ToId: to.Id, Message: data["message"].(string), DateStr: time.Now().String(), ProfilePhoto: from.Profile}

			jsonData := &JsonData{Event: "private_message", Data: msg}
			jsonM, _ := json.Marshal(jsonData)

			// send message to users
			to.Client.SendMessage(jsonM)
			from.Client.SendMessage(jsonM)
		}
	}

}

func (l *Lobby) findUserByClient(client *goreal.Client) (User, error) {
	var u User
	for _, v := range l.Users {
		if v.Client == client {
			return v, nil
		}
	}
	return u, errors.New("user not found")
}

func NewLobby() *Lobby {
	return &Lobby{Users: make(map[*goreal.Client]User)}
}
