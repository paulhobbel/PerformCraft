package status

import (
	"github.com/paulhobbel/performcraft/core"
	"github.com/paulhobbel/performcraft/core/network"
	"github.com/paulhobbel/performcraft/core/proto"
)

func StatusRequestHandler(request *proto.ClientPacketStatusRequest, session network.Session) {
	server := core.Instance()

	response := &proto.ServerPacketStatusResponse{}

	serverVersion := server.GetVersion()
	response.Version.Name = serverVersion.GetName()
	response.Version.Protocol = int(serverVersion.GetProtocol())
	response.Description = "PerformCraft Research Server"

	response.Players.Online = 0
	response.Players.Max = 20

	session.WritePacket(response)
}
