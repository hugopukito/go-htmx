package service

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]any)
var bdcast = make(chan string)

func HandleWsConnection(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer ws.Close()

	clients[ws] = struct{}{}

	for {
		var buff map[string]any
		err := ws.ReadJSON(&buff)
		if _, ok := buff["test"]; ok {
			bdcast <- "<div id='messageContent' hx-swap-oob='true'>Message: </div>"
		}

		if err != nil {
			delete(clients, ws)
			break
		}
	}
}

func handleMessages() {
	for {
		msg := <-bdcast
		messageClients(msg)
	}
}

func messageClients(msg string) {
	for client := range clients {
		messageClient(client, msg)
	}
}

func messageClient(client *websocket.Conn, msg string) {
	err := client.WriteJSON(msg)

	if err != nil && unsafeError(err) {
		log.Printf("error: %v", err)
		client.Close()
		delete(clients, client)
	}
}

// If a message is sent while a client is closing, ignore the error
func unsafeError(err error) bool {
	return !websocket.IsCloseError(err, websocket.CloseGoingAway) && err != io.EOF
}

func init() {
	go handleMessages()
}
