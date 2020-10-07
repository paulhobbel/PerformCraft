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
			PlayTeleportConfirm: func() base.Packet {
				return &ClientPacketTeleportConfirm{}
			},
			PlayClientSettings: func() base.Packet {
				return &ClientPacketSettings{}
			},
			PlayPluginMessageClient: func() base.Packet {
				return &ClientPacketPluginMessage{}
			},
			PlayKeepAliveClient: func() base.Packet {
				return &ClientPacketKeepAlive{}
			},
		},
	}
}
