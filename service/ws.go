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

var clients = make(map[*websocket.Conn]string)
var bdcast = make(chan string)

func HandleWsConnection(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer ws.Close()

	clients[ws] = "rgb(0,0,0)"

	// init board game
	bdcast <- initBoard(20)
	bdcast <- initColorPicker()

	for {
		var buff map[string]any
		err := ws.ReadJSON(&buff)
		if cellID, ok := buff["cellID"]; ok {
			setCellColor(cellID, clients[ws])
		}
		if colorPicker, ok := buff["colorPicker"]; ok {
			setClientColor(ws, colorPicker)
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

func setClientColor(client *websocket.Conn, colorPicker any) {

}

func setCellColor(cellID any, rgbColor string) {
	id, isString := cellID.(string)
	if isString {
		htmlID := `id="boardCell_` + id + `"`
		vals := ` hx-vals='{"cellID": "` + id + `"}'`
		color := ` style="background-color: ` + rgbColor + `;"`
		bdcast <- `<div class="cell" ` + htmlID + vals + color + ` hx-swap-oob='outerHTML:#boardCell_` + cellID.(string) + `' ws-send></div>`
	}
}

func initColorPicker() string {
	htmlString := "<div id='color_picker' class='class-picker' hx-swap-oob='innerHTML:#color_picker'>"
	htmlString += "<div id='color_pick_0' class='color-pick selected'>"
	htmlString += "</div>"
	for i := 0; i < 3; i++ {
		htmlString += "<div id='color_pick_" + strconv.Itoa(i+1) + "' class='color-pick'>"
		htmlString += "</div>"
	}
	htmlString += "</div>"
	return htmlString
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
