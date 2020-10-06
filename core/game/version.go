package game

import "github.com/paulhobbel/performcraft/core/base"

type Version interface {
	GetProtocol() base.ProtocolVersion
	GetName() string
}

func NewVersion(name string, protocol int) Version {
	return &versionImpl{
		protocol: base.ProtocolVersion(protocol),
		name:     name,
	}
}

type versionImpl struct {
	protocol base.ProtocolVersion
	name     string
}

func (v versionImpl) GetProtocol() base.ProtocolVersion {
	return v.protocol
}

func (v versionImpl) GetName() string {
	return v.name
}
