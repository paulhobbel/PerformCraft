package chat

import (
	"encoding/json"
	"fmt"
	"io"
)

type Message message

type message struct {
	Text string `json:"text,omitempty"`

	Bold          bool   `json:"bold,omitempty"`
	Italic        bool   `json:"italic,omitempty"`
	Underlined    bool   `json:"underlined,omitempty"`
	Strikethrough bool   `json:"strikethrough,omitempty"`
	Obfuscated    bool   `json:"obfuscated,omitempty"`
	Color         string `json:"color,omitempty"`

	Translate string    `json:"translate,omitempty"`
	With      []string  `json:"with,omitempty"`
	Extra     []message `json:"extra,omitempty"`
}

func (m *Message) UnmarshalJSON(msg []byte) (err error) {
	if len(msg) == 0 {
		return io.EOF
	}

	// Just raw text
	if msg[0] == '"' {
		err = json.Unmarshal(msg, &m.Text)
	} else {
		err = json.Unmarshal(msg, (*message)(m))
	}

	return
}

func (m Message) String() string {
	bs, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("Message%s", string(bs))
}
