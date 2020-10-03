package r578

import (
	"github.com/paulhobbel/performcraft/pkg/common"
	"github.com/paulhobbel/performcraft/pkg/proto"
)

func init() {
	r578 := proto.Registry.RegisterProtocol(proto.R578)

	// Register Handshaking packets
	r578.
		RegisterClientPacket(common.Handshaking, HandshakingHandshake, &ClientPacketHandshake{})

	// Register Status packets
	r578.
		RegisterClientPacket(common.Status, StatusRequest, &ClientPacketStatusRequest{}).
		RegisterClientPacket(common.Status, StatusPing, &PacketStatusPing{}).
		RegisterServerPacket(common.Status, StatusResponse, &ServerPacketStatusResponse{}).
		RegisterServerPacket(common.Status, StatusPong, &PacketStatusPing{})

	// Register Login packets
	r578.
		RegisterClientPacket(common.Login, LoginStart, &ClientPacketLoginStart{}).
		RegisterServerPacket(common.Login, LoginDisconnect, &ServerPacketLoginDisconnect{}).
		RegisterServerPacket(common.Login, LoginSuccess, &ServerPacketLoginSuccess{}).
		RegisterServerPacket(common.Login, LoginSetCompression, &ServerPacketSetCompression{})

	// Register Play packets
	r578.
		RegisterServerPacket(common.Play, PlayJoinGame, &ServerPacketPlayJoinGame{}).
		RegisterServerPacket(common.Play, PlayPlayerPositionLookServer, &ServerPacketPlayPlayerPositionLook{})
}
