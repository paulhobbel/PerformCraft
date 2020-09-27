package main

import (
	"flag"
	"fmt"
	"github.com/paulhobbel/performcraft/pkg/net"
	"github.com/paulhobbel/performcraft/pkg/net/packet"
	"log"
)

var sourceAddr = flag.String("s", "localhost:25569", "source address")
var targetAddr = flag.String("t", "localhost:25565", "target address")

func handleConn(clientConn *net.Conn) {
	serverConn, err := net.DialMC(*targetAddr)
	if err != nil {
		panic(err)
	}

	//log.Println("Got new client")

	// Handle client
	go func() {
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

			// Player position
			if p.ID == 0x12 {
				var pX packet.Double
				var pY packet.Double
				var pZ packet.Double
				var pOnGround packet.Boolean

				p.Scan(&pX, &pY, &pZ, &pOnGround)

				log.Printf("[Client] X: %v, Y: %v, Z: %v, IsOnGround: %v", pX, pY, pZ, pOnGround)
			} else {
				log.Printf("[Client]: Incoming packet %v\n", p)
			}

			err = serverConn.WritePacket(p)
			if err != nil {
				log.Println("[Client] Failed sending client packet: ", err)
				return
			}
		}
	}()

	// Handle server
	go func() {
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

			//log.Printf("[Server]: Incoming packet %v\n", p)

			// Got SetCompression packet, setup client and server.
			if p.ID == 0x03 {
				var threshold packet.VarInt
				p.Scan(&threshold)
				serverConn.SetThreshold(int(threshold))

				log.Printf("[Server]: Got SetCompression packet, set to threshold: %v\n", threshold)

				// Let's not send the SetCompression packet to client as this messes a lot of shit up.
				continue
			}

			err = clientConn.WritePacket(p)
			if err != nil {
				log.Println("[Server] Failed sending server packet: ", err)
				return
			}

		}
	}()
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
