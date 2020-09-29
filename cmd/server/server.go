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

	info.Version.Name = "PerformCraft 1.15.2"
	info.Version.Protocol = 578
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
		log.Println("[Server]: Failed reading login packet:", err)
		return
	}

	var Username packet.String
	err = p.Scan(&Username)
	log.Printf("[Server]: Starting login for %v\n", Username)

	// Accept login.
	err = conn.WritePacket(packet.Marshal(0x02, packet.String(uuid.MustParse("74242c15-feb0-43b7-8045-2d4a602b2d74").String()), Username))
	if err != nil {
		log.Println("[Server]: Failed sending login accepted packet:", err)
		return
	}

	// Join Game.
	err = conn.WritePacket(packet.Marshal(0x26,
		packet.Int(0),            // EntityID
		packet.UnsignedByte(1),   // Gamemode
		packet.Int(0),            // Dimension
		packet.Long(0),           // HashedSeed
		packet.UnsignedByte(200), // MaxPlayer
		packet.String("default"), // LevelType
		packet.VarInt(15),        // View Distance
		packet.Boolean(false),    // Reduced Debug Info
		packet.Boolean(true),     // Enable respawn screen
	))

	err = conn.WritePacket(packet.Marshal(0x36,
		packet.Double(0), packet.Double(0), packet.Double(0), // XYZ
		packet.Float(0), packet.Float(0), // Yaw Pitch
		packet.Byte(0),   // flag
		packet.VarInt(0), // TP ID
	))

	// Keep the connection alive.
	for {
		p, err := conn.ReadPacket()

		log.Println("[Server]: Inbound packet", p)

		if err != nil {
			log.Println("[Server]: ReadPacket failed:", err)
			return
		}
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

	log.Println(p.ID, p.Data)

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
		log.Println("[Server]: Get play server request")
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
