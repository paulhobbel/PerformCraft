package proto

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/paulhobbel/performcraft/core/base"
	"github.com/paulhobbel/performcraft/core/bufio"
)

type ServerPacketStatusResponse struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
		Sample []struct {
			Name string    `json:"name"`
			ID   uuid.UUID `json:"id"`
		}
	} `json:"players"`
	Description string `json:"description"`
	FavIcon     string `json:"favicon,omitempty"`
}

func (p ServerPacketStatusResponse) ID() base.PacketID {
	return StatusResponse
}

func (p *ServerPacketStatusResponse) Read(b bufio.ByteBuffer) error {
	status, err := b.ReadString()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(status), &p)
}

func (p ServerPacketStatusResponse) Write(b bufio.ByteBuffer) error {
	status, err := json.Marshal(p)
	if err != nil {
		return err
	}

	return b.WriteString(string(status))
}

func (p ServerPacketStatusResponse) String() string {
	status, _ := json.Marshal(p)
	return fmt.Sprintf("ServerPacketStatusResponse{%s}", status)
}
