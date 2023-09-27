package service

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"

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

	rand := rand.New(rand.NewSource(rand.Int63()))

	for {
		var buff map[string]any
		err := ws.ReadJSON(&buff)
		if _, ok := buff["test"]; ok {
			//bdcast <- `<div hx-swap-oob='innerHTML:#msg'>` + msg.(string) + `</div>`
			bdcast <- `<div hx-swap-oob='innerHTML:#msg'>` + strconv.Itoa(rand.Int()) + `</div>`
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
	err := client.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		fmt.Println("Error writing message:", err)
		return
	}

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
