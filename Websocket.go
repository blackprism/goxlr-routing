package main

import (
	"github.com/blackprism/goxlr-routing/Key"
	gwebsocket "github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Websocket struct {
	listener Key.Listener
	upgrader gwebsocket.Upgrader
}

func NewWebsocket(listener Key.Listener) Websocket {
	var upgrader = gwebsocket.Upgrader{}

	return Websocket{
		listener,
		upgrader,
	}
}

func (w Websocket) Listen(addr string) {
	http.HandleFunc("/", w.handler)

	log.Println("Start websocket server")
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (w Websocket) handler(writer http.ResponseWriter, request *http.Request) {
	connection, err := w.upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Print("Can't convert Http request to websocket:", err)
		return
	}
	defer connection.Close()

	log.Println("Connection from", connection.RemoteAddr())

	w.listener.Listen(connection)
}
