package main

import (
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/paulhobbel/performcraft/pkg/common"
	"github.com/paulhobbel/performcraft/pkg/net"
	v2 "github.com/paulhobbel/performcraft/pkg/net/packet/v2"
	"github.com/paulhobbel/performcraft/pkg/proto"
	"github.com/paulhobbel/performcraft/pkg/proto/r578"
	"github.com/paulhobbel/performcraft/pkg/server"
	"github.com/paulhobbel/performcraft/pkg/server/network/handler"
	"log"
)

var listenAddr = flag.String("l", "localhost:25569", "listen address")

func createStatusResponse() *r578.ServerPacketStatusResponse {
	info := &r578.ServerPacketStatusResponse{}

	info.Version.Name = "PerformCraft 1.15.2"
	info.Version.Protocol = 578
	info.Players.Max = 1337
	info.Players.Online = 420
	//info.Players.Sample =
	info.Description = "PerformCraft Research Server"

	return info
}

func main() {
	flag.Parse()

	fmt.Println("Starting PerformCraft v0.0.0")

	srv := server.NewServer(proto.R578)

	srv.SubscribeAs(handler.HandshakeHandler)

	//srv.SubscribeAs(func(handshake *r578.ClientPacketHandshake, conn *net.Conn) {
	//	log.Printf("[Server]: Got handshake: %v", handshake)
	//	conn.SetProtocolState(handshake.State)
	//})

	srv.SubscribeAs(func(request *r578.ClientPacketStatusRequest, conn *net.Conn) {
		log.Println("[Server]: Got status request, sending info")
		conn.WritePacket(createStatusResponse())
	})

	srv.SubscribeAs(func(ping *r578.PacketStatusPing, conn *net.Conn) {
		log.Println("[Server]: Sending ping")
		conn.WritePacket(ping)
	})

	srv.SubscribeAs(func(login *r578.ClientPacketLoginStart, conn *net.Conn) {
		log.Printf("[Server]: Got login start: %+v", login)

		//if ver != proto.R578 {
		//	log.Printf("[Server]: Incorrect protocol version, got %d, expected: %d", ver, proto.R578)
		//
		//	translation := "multiplayer.disconnect.outdated_client"
		//	if ver > proto.R578 {
		//		translation = "multiplayer.disconnect.outdated_server"
		//	}
		//
		//	conn.WritePacket(&r578.ServerPacketLoginDisconnect{
		//		Reason: chat.Message{
		//			Translate: translation,
		//			With:      []string{"1.15.2"},
		//		},
		//	})
		//	return
		//}

		err := conn.WritePacket(&r578.ServerPacketLoginSuccess{
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
		conn.SetProtocolState(common.Play)

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
	})

	srv.SubscribeAs(func(packet *v2.BasePacket) {
		log.Printf("[Server]: Received unknown packet: %+v", packet)
	})

	if err := srv.Listen(25569); err != nil {
		panic(err)
	}
}
