package main

import (
	"github.com/paulhobbel/performcraft/core"
	"github.com/paulhobbel/performcraft/handlers/handshake"
	"github.com/paulhobbel/performcraft/handlers/login"
	"github.com/paulhobbel/performcraft/handlers/status"
)

func main() {
	srv := core.Instance()

	// TODO: Find better way to register these
	srv.GetNetwork().Subscribe(handshake.HandshakeHandler)
	srv.GetNetwork().Subscribe(status.StatusPingHandler)
	srv.GetNetwork().Subscribe(status.StatusRequestHandler)
	srv.GetNetwork().Subscribe(login.LoginStartHandler)

	srv.Start()

	//<- signal.Notify(ch)
}
