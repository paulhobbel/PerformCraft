package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/paulhobbel/performcraft/pkg/net"
	"github.com/paulhobbel/performcraft/pkg/net/packet"
	"io"
	"log"
)

var listenAddr = flag.String("l", "localhost:25569", "listen address")

type ServerInfoPlayer struct {
	Name string    `json:"name"`
	ID   uuid.UUID `json:"id"`
}

type ServerInfo struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
		Sample []ServerInfoPlayer
	} `json:"players"`
	Description string `json:"description"`
	FavIcon     string `json:"favicon,omitempty"`
}

func createListResponse() string {
	var info ServerInfo

	info.Version.Name = "PerformCraft 1.16.3"
	info.Version.Protocol = 753
	info.Players.Max = 1337
	info.Players.Online = 420
	info.Players.Sample = []ServerInfoPlayer{}
	info.Description = "PerformCraft Research Server"

	data, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}

	return string(data)
}

func acceptListPing(conn *net.Conn) {
	for i := 0; i < 2; i++ {
		p, err := conn.ReadPacket()
		if err != nil {
			return
		}

		switch p.ID {
		case 0x00: // list
			log.Println("[Server]: Sending list info")
			err = conn.WritePacket(packet.Marshal(0x00, packet.String(createListResponse())))
		case 0x01: // ping
			log.Println("[Server]: Sending ping")
			err = conn.WritePacket(p)
		}

		if err != nil {
			return
		}
	}
}

func acceptLogin(conn *net.Conn) {
	p, err := conn.ReadPacket()
	if err != nil {
		return
	}

	switch p.ID {
	case 0x00:
		var Username packet.String
		err = p.Scan(&Username)
		log.Printf("[Server]: Starting login for %v", Username)

		err = conn.WritePacket(packet.Marshal(0x02, packet.String("74242c15-feb0-43b7-8045-2d4a602b2d74"), Username))
	}

	if err != nil {
		return
	}
}

func handleConn(conn *net.Conn) {
	defer conn.Close()

	// handshake
	p, err := conn.ReadPacket()
	if err != nil {
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			log.Println("[Server]: Failed handshaking: ", err)
		}
		return
	}

	log.Println(p.Data)

	var (
		Protocol, Intention packet.VarInt
		ServerAddress       packet.String
		ServerPort          packet.UnsignedShort
	)

	err = p.Scan(&Protocol, &ServerAddress, &ServerPort, &Intention)
	if err != nil {
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			log.Println("[Server]: Failed handshaking: ", err)
		}
		log.Println("[Server]: Failed handshaking: ", err)
	}

	log.Printf("[Server]: Got handshake { proto: %v, addr: %v, port: %v, intention: %v }", Protocol, ServerAddress, ServerPort, Intention)

	switch Intention {
	case 1: // status
		log.Println("[Server]: Got server list info request")
		acceptListPing(conn)
	case 2: // login
		log.Println("[Server]: Sawwy, not supporting login procedure yet.")
		acceptLogin(conn)
	}
}

func main() {
	flag.Parse()

	fmt.Println("Starting PerformCraft v0.0.0")

	listener, err := net.ListenMC(*listenAddr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Listening at: %v\n", *listenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(&conn)
	}
}
