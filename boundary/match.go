package boundary

import (
	"encoding/json"
	"log"
	"net/http"
	"pong/controller"
	"pong/entity"
	"pong/repository"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MatchHandler struct {
	playerstore *repository.UserStore
	simulator   *controller.Lobby
	upgrader    websocket.Upgrader
}

func NewMatchHanlder(p *repository.UserStore, s *controller.Lobby, u websocket.Upgrader) MatchHandler {
	return MatchHandler{playerstore: p, simulator: s, upgrader: u}
}

func (m MatchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if len(idStr) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Println("User did not submit id")
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Println("User submitted invalid id: ", idStr, " error: ", err)
		return
	}
	player, ok := m.playerstore.Get(id)
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Println("User submitted unkown id: ", idStr, " error: ", err)
		return
	}
	conn, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Could not establish websocket connection err: ", err)
		return
	}
	server := make(chan controller.ClientMessage)
	client := make(chan controller.ServerMessage)
	go read(conn, server, *player)
	go write(conn, client, *player)
	observer := controller.Observer{Player: player, Server: server, Client: client}
	m.simulator.Add(observer)
}

func read(conn *websocket.Conn, server chan<- controller.ClientMessage, player entity.Player) {
	log.Println(player, " WS read goroutine opened")
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Websocket closed")
			break
		}
		clientMessage := controller.ClientMessage{}
		err = json.Unmarshal(message, &clientMessage)
		if err != nil {
			log.Println("Player submited invalid input: ", message, " err: ", err)
			continue
		}
		// log.Println(clientMessage)
		server <- clientMessage
	}
	log.Println(player, " WS read goroutine closed")
	conn.Close()
}

func write(conn *websocket.Conn, client <-chan controller.ServerMessage, player entity.Player) {
	log.Println(player, " WS write goroutine opened")
	for message := range client {
		json, err := json.Marshal(message.Game)
		if err != nil {
			log.Println("Could not create JSON from: ", message, " err: ", err)
			continue
		}
		conn.WriteJSON(json)
		// log.Println("Message sent to ", player)
	}
	log.Println(player, " WS write goroutine closed")
	conn.Close()
}
