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
var bdcastAll = make(chan string)
var colors = []string{
	"rgb(5, 59, 80)",
	"rgb(23, 107, 135)",
	"rgb(100, 204, 197)",
}

func HandleWsConnection(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer ws.Close()

	clients[ws] = colors[0]

	// init board game
	bdcastAll <- initBoard(100)
	bdcastAll <- initColorPicker(0)

	for {
		var buff map[string]any
		err := ws.ReadJSON(&buff)
		if cellID, ok := buff["cellID"]; ok {
			setCellColor(ws, cellID)
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
		msg := <-bdcastAll
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

func setClientColor(ws *websocket.Conn, colorPicker any) {
	idString, ok := colorPicker.(string)
	id, ok2 := strconv.Atoi(idString)
	if ok && ok2 == nil {
		clients[ws] = colors[id]
		messageClient(ws, initColorPicker(id))
	}
}

func setCellColor(ws *websocket.Conn, cellID any) {
	id, ok := cellID.(string)
	if ok {
		htmlID := `id="boardCell_` + id + `"`
		vals := ` hx-vals='{"cellID": "` + id + `"}'`
		color := ` style="background-color: ` + clients[ws] + `;"`
		bdcastAll <- `<div class="cell" ` + htmlID + vals + color + ` hx-swap-oob='outerHTML:#boardCell_` + cellID.(string) + `' ws-send></div>`
	}
}

func initColorPicker(targetID int) string {
	htmlString := "<div id='color_picker' hx-swap-oob='innerHTML:#color_picker'>"
	for i := 0; i < len(colors); i++ {
		id := strconv.Itoa(i)
		vals := ` hx-vals='{"colorPicker": "` + id + `"}'`
		color := ` style="background-color: ` + colors[i]
		if i == targetID {
			color += `; border: 3px solid rgb(255, 75, 145);"`
		} else {
			color += `;"`
		}
		htmlString += "<div id='color_pick_" + id + "'" + vals + color + "' class='color-pick' ws-send>"
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
