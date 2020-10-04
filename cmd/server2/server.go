package main

import (
	"github.com/paulhobbel/performcraft/impl"
)

func main() {
	srv := impl.NewServer()
	srv.GetNetwork().Listen(25569)

	//<- signal.Notify(ch)
}
