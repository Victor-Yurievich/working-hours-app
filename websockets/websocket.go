package websockets

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var WS *websocket.Conn // Ask Lior about scopes

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	WS = conn
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Client successfully connected...")
	reader(WS)
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func LogUserOut() {
	if WS == nil {
		return
	}
	if err := WS.WriteMessage(1, []byte("/logout")); err != nil {
		log.Println(err)
		return
	}
}
