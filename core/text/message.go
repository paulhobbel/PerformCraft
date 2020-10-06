package text

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Text string `json:"text,omitempty"`

	Bold          bool   `json:"bold,omitempty"`
	Italic        bool   `json:"italic,omitempty"`
	Underlined    bool   `json:"underlined,omitempty"`
	Strikethrough bool   `json:"strikethrough,omitempty"`
	Obfuscated    bool   `json:"obfuscated,omitempty"`
	Color         string `json:"color,omitempty"`

	Translate string    `json:"translate,omitempty"`
	With      []string  `json:"with,omitempty"`
	Extra     []Message `json:"extra,omitempty"`
}

func (m Message) String() string {
	bs, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("Message%s", string(bs))
}
