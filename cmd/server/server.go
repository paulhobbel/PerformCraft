package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/paulhobbel/performcraft/pkg/common"
	"github.com/paulhobbel/performcraft/pkg/net"
	"github.com/paulhobbel/performcraft/pkg/net/chat"
	"github.com/paulhobbel/performcraft/pkg/proto"
	"github.com/paulhobbel/performcraft/pkg/proto/r578"
	"io"
	"log"
)

var listenAddr = flag.String("l", "localhost:25569", "listen address")

func createListResponse() *r578.ServerPacketStatusResponse {
	info := &r578.ServerPacketStatusResponse{}

	info.Version.Name = "PerformCraft 1.15.2"
	info.Version.Protocol = 578
	info.Players.Max = 1337
	info.Players.Online = 420
	//info.Players.Sample =
	info.Description = "PerformCraft Research Server"

	return info
}

func acceptLogin(conn *net.Conn, ver proto.ProtocolVersion) {
	p, err := conn.ReadPacket()
	if err != nil {
		log.Println("[Server]: Failed reading login packet:", err)
		return
	}

	if login, ok := p.(*r578.ClientPacketLoginStart); ok {
		log.Printf("[Server]: Got login start: %+v", login)

		if ver != proto.R578 {
			log.Printf("[Server]: Incorrect protocol version, got %d, expected: %d", ver, proto.R578)

			translation := "multiplayer.disconnect.outdated_client"
			if ver > proto.R578 {
				translation = "multiplayer.disconnect.outdated_server"
			}

			conn.WritePacket(&r578.ServerPacketLoginDisconnect{
				Reason: chat.Message{
					Translate: translation,
					With:      []string{"1.15.2"},
				},
			})
			return
		}

		err = conn.WritePacket(&r578.ServerPacketLoginSuccess{
			Name: login.Name,
			UUID: uuid.MustParse("74242c15-feb0-43b7-8045-2d4a602b2d74").String(),
		})
		if err != nil {
			log.Println("[Server]: Failed sending login accepted packet:", err)
			return
		}

		// Join Game
		err = conn.WritePacket(&r578.ServerPacketPlayJoinGame{
			EntityId:          0,
			GameMode:          1,
			Dimension:         0,
			Seed:              0,
			MaxPlayers:        0,
			LevelType:         "default",
			ViewDistance:      15,
			ReducedDebugInfo:  false,
			EnableSpawnScreen: true,
		})

		// Switch to play state
		conn.SetState(common.Play)

		// Send Position
		err = conn.WritePacket(&r578.ServerPacketPlayPlayerPositionLook{
			X:          0,
			Y:          0,
			Z:          0,
			Yaw:        0,
			Pitch:      0,
			Flags:      0,
			TeleportID: 0,
		})
	}

	//var Username packet.String
	//err = p.Scan(&Username)
	//log.Printf("[Server]: Starting login for %v\n", Username)
	//
	//// Accept login.
	//err = conn.WritePacket(packet.Marshal(0x02, packet.String(uuid.MustParse("74242c15-feb0-43b7-8045-2d4a602b2d74").String()), Username))
	//if err != nil {
	//	log.Println("[Server]: Failed sending login accepted packet:", err)
	//	return
	//}
	//
	//// Join Game.
	//err = conn.WritePacket(packet.Marshal(0x26,
	//	packet.Int(0),            // EntityID
	//	packet.UnsignedByte(1),   // Gamemode
	//	packet.Int(0),            // Dimension
	//	packet.Long(0),           // HashedSeed
	//	packet.UnsignedByte(200), // MaxPlayer
	//	packet.String("default"), // LevelType
	//	packet.VarInt(15),        // View Distance
	//	packet.Boolean(false),    // Reduced Debug Info
	//	packet.Boolean(true),     // Enable respawn screen
	//))
	//
	//err = conn.WritePacket(packet.Marshal(0x36,
	//	packet.Double(0), packet.Double(0), packet.Double(0), // XYZ
	//	packet.Float(0), packet.Float(0), // Yaw Pitch
	//	packet.Byte(0),   // flag
	//	packet.VarInt(0), // TP ID
	//))

	// Keep the connection alive.
	for {
		p, err := conn.ReadPacket()

		log.Printf("[Server]: Inbound packet %+v", p)

		if err != nil {
			log.Println("[Server]: ReadPacket failed:", err)
			return
		}
	}
}

func handleConn(conn *net.Conn) {
	defer conn.Close()

	def := proto.Registry.GetDefinition(proto.R578)

	conn.Decoder.SetPacketFactory(def.GetClientPacket)

	// handshake
	p, err := conn.ReadPacket()
	if err != nil {
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			log.Println("[Server]: Failed handshaking: ", err)
		}
		return
	}

	if handshake, ok := p.(*r578.ClientPacketHandshake); ok {
		conn.SetState(handshake.State)

		if handshake.State == common.Status {
			for i := 0; i < 2; i++ {
				p, err = conn.ReadPacket()
				switch p.(type) {
				case *r578.ClientPacketStatusRequest:
					log.Println("[Server]: Sending list info")
					conn.WritePacket(createListResponse())
				case *r578.PacketStatusPing:
					log.Println("[Server]: Sending ping")
					conn.WritePacket(p)
				}
			}

			return
		}

		if handshake.State == common.Login {
			acceptLogin(conn, proto.ProtocolVersion(handshake.Version))
		}
	}

	//log.Println(p.ID, p.Data)
	//
	//var (
	//	Protocol, Intention packet.VarInt
	//	ServerAddress       packet.String
	//	ServerPort          packet.UnsignedShort
	//)
	//
	//err = p.Scan(&Protocol, &ServerAddress, &ServerPort, &Intention)
	//if err != nil {
	//	if !errors.Is(err, io.ErrUnexpectedEOF) {
	//		log.Println("[Server]: Failed handshaking: ", err)
	//	}
	//	log.Println("[Server]: Failed handshaking: ", err)
	//}
	//
	//log.Printf("[Server]: Got handshake { proto: %v, addr: %v, port: %v, intention: %v }", Protocol, ServerAddress, ServerPort, Intention)
	//
	//switch Intention {
	//case 1: // status
	//	log.Println("[Server]: Got server list info request")
	//	acceptListPing(conn)
	//case 2: // login
	//	log.Println("[Server]: Get play server request")
	//	acceptLogin(conn)
	//}
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
