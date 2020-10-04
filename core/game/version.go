package game

import "github.com/paulhobbel/performcraft/core/proto"

type Version interface {
	GetProtocol() proto.ProtocolVersion
	GetName() string
}

func NewVersion(name string, protocol int) Version {
	return &versionImpl{
		protocol: proto.ProtocolVersion(protocol),
		name:     name,
	}
}

type versionImpl struct {
	protocol proto.ProtocolVersion
	name     string
}

func (v versionImpl) GetProtocol() proto.ProtocolVersion {
	return v.protocol
}

func (v versionImpl) GetName() string {
	return v.name
}
