package service

import (
	"fmt"
	"io"
	"log"
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

	// init board game
	bdcast <- initBoard(3)

	for {
		var buff map[string]any
		err := ws.ReadJSON(&buff)
		if cellID, ok := buff["cellID"]; ok {
			htmlID := `id="boardCell_` + cellID.(string) + `"`
			vals := ` hx-vals='{"cellID": "` + cellID.(string) + `"}'`
			bdcast <- `<div class="cell colored" ` + htmlID + vals + ` hx-swap-oob='outerHTML:#boardCell_` + cellID.(string) + `' ws-send>x</div>`
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

func initBoard(size int) string {
	htmlString := "<div hx-swap-oob='innerHTML:#board'>"
	for i := 0; i < size; i++ {
		htmlString += "<div class='line'>"
		for y := 0; y < size; y++ {
			id := strconv.Itoa(i) + "_" + strconv.Itoa(y)
			htmlID := `id="boardCell_` + id + `"`
			vals := ` hx-vals='{"cellID": "` + id + `"}'`
			htmlString += `<div class="cell" ` + htmlID + vals + ` ws-send></div>`
		}
		htmlString += "</div>"
	}
	htmlString += "</div>"
	return htmlString
}

func init() {
	go handleMessages()
}
