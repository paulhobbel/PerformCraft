package main

import (
	"flag"
	"fmt"
	"github.com/paulhobbel/performcraft/pkg/common"
	"github.com/paulhobbel/performcraft/pkg/net"
	"github.com/paulhobbel/performcraft/pkg/proto"
	"github.com/paulhobbel/performcraft/pkg/proto/r578"
	"log"

	_ "github.com/paulhobbel/performcraft/pkg/proto/r578"
	_ "github.com/paulhobbel/performcraft/pkg/proto/r753"
)

var sourceAddr = flag.String("s", "localhost:25569", "source address")
var targetAddr = flag.String("t", "localhost:25565", "target address")

type GameState int

const (
	Handshake GameState = iota
	Login
	Play
)

func handleConn(clientConn *net.Conn) {
	serverConn, err := net.DialMC(*targetAddr)
	if err != nil {
		panic(err)
	}

	//log.Println("Got new client")

	//var gameState GameState = Login

	// TODO: Place elsewhere
	def := proto.Registry.GetDefinition(proto.R578)

	clientConn.Decoder.SetPacketFactory(def.GetClientPacket)
	serverConn.Decoder.SetPacketFactory(def.GetServerPacket)

	// tempfix for concurrency
	serverConn.SetProtocolState(common.Login)

	// Handle client
	go func(clientConn, serverConn *net.Conn) {
		defer func() {
			defer clientConn.Close()
			defer serverConn.Close()
			log.Println("Stopped listening for client packets")
		}()

		log.Println("Listening for client packets")

		for {
			p, err := clientConn.ReadPacket()
			if err != nil {
				log.Println("[Client] Failed reading client packet: ", err)
				return
			}

			log.Printf("[Client]: Incoming packet %+v\n", p)

			switch clientConn.State {
			case common.Handshake:
				switch packet := p.(type) {
				case *r578.ClientPacketHandshake:
					clientConn.SetProtocolState(common.ProtocolState(packet.State))
					serverConn.SetProtocolState(common.ProtocolState(packet.State))
				}
			}

			//switch gameState {
			//case Play:
			//	if p.ID() == 0x12 { // Player Position
			//		var pX packet.Double
			//		var pY packet.Double
			//		var pZ packet.Double
			//		var pOnGround packet.Boolean
			//
			//		p.Scan(&pX, &pY, &pZ, &pOnGround)
			//
			//		if !pOnGround {
			//			log.Println("[Client]: Client is falling...overriding :P")
			//			pOnGround = true
			//			p = packet.Marshal(p.ID, &pX, &pY, &pZ, &pOnGround)
			//		}
			//
			//		//log.Printf("[Client] X: %v, Y: %v, Z: %v, IsOnGround: %v", pX, pY, pZ, pOnGround)
			//	} else if p.ID() == 0x13 { // Player Position And Rotation
			//		var pX packet.Double
			//		var pY packet.Double
			//		var pZ packet.Double
			//		var pYaw packet.Float
			//		var pPitch packet.Float
			//		var pOnGround packet.Boolean
			//
			//		p.Scan(&pX, &pY, &pZ, &pYaw, &pPitch, &pOnGround)
			//
			//		if !pOnGround {
			//			log.Println("[Client]: Client is falling while moving & rotating...overriding :P")
			//			pOnGround = true
			//			p = packet.Marshal(p.ID, &pX, &pY, &pZ, &pYaw, &pPitch, &pOnGround)
			//		}
			//
			//		//log.Printf("[Client] Yaw: %v, Pitch: %v, IsOnGround: %v", pYaw, pPitch, pOnGround)
			//	} else if p.ID() == 0x14 { // Player Rotation
			//		var pYaw packet.Float
			//		var pPitch packet.Float
			//		var pOnGround packet.Boolean
			//
			//		p.Scan(&pYaw, &pPitch, &pOnGround)
			//
			//		if !pOnGround {
			//			log.Println("[Client]: Client is falling while rotating...overriding :P")
			//			pOnGround = true
			//			p = packet.Marshal(p.ID, &pYaw, &pPitch, &pOnGround)
			//		}
			//	} else {
			//		log.Printf("[Client]: Incoming packet %v\n", p)
			//	}
			//default:
			//	log.Printf("[Client]: Incoming packet %v\n", p)
			//}

			err = serverConn.WritePacket(p)
			if err != nil {
				log.Println("[Client] Failed sending client packet: ", err)
				return
			}
		}
	}(clientConn, serverConn)

	// Handle server
	go func(clientConn, serverConn *net.Conn) {
		defer func() {
			defer clientConn.Close()
			defer serverConn.Close()
			log.Println("Stopped listening for server packets")
		}()

		log.Println("Listening for server packets")

		for {
			p, err := serverConn.ReadPacket()
			if err != nil {
				log.Println("[Server] Failed reading server packet: ", err)
				return
			}

			//log.Printf("[Server]: Incoming packet %+v\n", p)

			switch serverConn.State {
			case common.Login:
				switch packet := p.(type) {
				case *r578.ServerPacketSetCompression:
					serverConn.SetThreshold(int(packet.Threshold))
					log.Printf("[Server]: Got SetCompression packet, set to threshold: %v\n", 256)

					// Let's not send the SetCompression packet to client as this messes a lot of shit up.
					continue
				}

				if p.ID() == 0x24 {
					serverConn.SetProtocolState(common.Play)
					clientConn.SetProtocolState(common.Play)
				}
			}

			//switch gameState {
			//case Login:
			//	//// Auth success
			//	//if p.ID() == 0x02 {
			//	//	var (
			//	//		clientUUID     packet.UUID
			//	//		clientUsername packet.String
			//	//	)
			//	//
			//	//	p.Scan(&clientUUID, &clientUsername)
			//	//
			//	//	log.Printf("[Server]: Got successful login { uuid: %v, username: %v }", clientUUID, clientUsername)
			//	//}
			//
			//	// Switch to play state
			//	if p.ID() == 0x24 {
			//	//	var (
			//	//		EntityID         packet.Int
			//	//		IsHardcore       packet.Boolean
			//	//		Gamemode         packet.UnsignedByte
			//	//		PreviousGamemode packet.UnsignedByte
			//	//		WorldNames       packet.StringArray
			//	//		DimensionCodec   packet.NBT
			//	//		Dimension        packet.NBT
			//	//		WorldName        packet.String
			//	//	)
			//	//
			//	//	err = p.Scan(&EntityID, &IsHardcore, &Gamemode, &PreviousGamemode, &WorldNames, &DimensionCodec, &WorldName)
			//	//	if err != nil {
			//	//		log.Printf("Failed scanning LoginStartPacket: %v", err)
			//	//	}
			//	//	log.Printf("LoginStart: {entityId: %v, isHardcore: %v, gamemode: %v, prevGamemode: %v, worldNames: %v, dimensionCodec: %v, dimension: %v, worldName: %v}",
			//	//		EntityID, IsHardcore, Gamemode, PreviousGamemode, WorldNames, DimensionCodec.V, Dimension, WorldName)
			//
			//		gameState = Play
			//	}
			//	//case Play:
			//	//	log.Printf("[Server]: Incoming packet %v\n", p)
			//}

			err = clientConn.WritePacket(p)
			if err != nil {
				log.Println("[Server] Failed sending server packet: ", err)
				return
			}

		}
	}(clientConn, serverConn)
}

func main() {
	flag.Parse()

	fmt.Printf("Starting proxy to: %v\n", *targetAddr)

	listener, err := net.ListenMC(*sourceAddr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Listening at: %v\n", *sourceAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(&conn)
	}
}
