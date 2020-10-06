package proto

import "github.com/paulhobbel/performcraft/core/base"

func init() {
	clientPackets = map[base.ProtocolState]map[base.PacketID]func() base.Packet{
		base.Handshake: {
			HandshakingHandshake: func() base.Packet {
				return &ClientPacketHandshake{}
			},
		},
		base.Login: {
			LoginStart: func() base.Packet {
				return &ClientPacketLoginStart{}
			},
		},
		base.Status: {
			StatusPing: func() base.Packet {
				return &PacketStatusPing{}
			},
			StatusRequest: func() base.Packet {
				return &ClientPacketStatusRequest{}
			},
		},
		base.Play: {
			0x0f: func() base.Packet {
				return &ClientPacketKeepAlive{}
			},
		},
	}
}
