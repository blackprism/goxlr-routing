package main

import (
	"github.com/blackprism/goxlr-routing/Config"
	"github.com/blackprism/goxlr-routing/GOXLR"
	"github.com/blackprism/goxlr-routing/Key"
)

func main() {
	payload := GOXLR.Payload{}
	keyListener := Key.NewListener(Config.Json("config.json"), payload)
	keyListener.Register()

	websocket := NewWebsocket(keyListener)
	websocket.Listen("127.0.0.1:6805")
}
