package login

import (
	"github.com/google/uuid"
	"github.com/paulhobbel/performcraft/core/game"
	"github.com/paulhobbel/performcraft/core/network"
	"github.com/paulhobbel/performcraft/core/proto"
)

func LoginStartHandler(request *proto.ClientPacketLoginStart, session network.Session) {
	//id, err := uuid.FromBytes([]byte("OfflinePlayer:" + request.Name))
	//if err != nil {
	//	log.Printf("[Login]: Failed generating id: %v", err)
	//}

	// Enable compression
	session.WritePacket(&proto.ServerPacketSetCompression{Threshold: 256})
	session.SetThreshold(256)

	profile := &game.Profile{
		UUID: uuid.New(),
		Name: request.Name,
	}

	// Cache players
	session.WritePacket(&proto.ServerPacketLoginSuccess{
		UUID: profile.UUID.String(),
		Name: profile.Name,
	})

	session.SetPlayerProfile(profile)

	// We go play
	session.InitPlayer()
}
