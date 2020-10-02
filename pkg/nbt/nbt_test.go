package nbt

import (
	"io/ioutil"
	"testing"
)

func TestUnmarshal_string(t *testing.T) {
	var data = []byte{
		0x08, 0x00, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x00, 0x09,
		0x42, 0x61, 0x6e, 0x61, 0x6e, 0x72, 0x61, 0x6d, 0x61,
	}

	var name string
	if err := Unmarshal(data, &name); err != nil {
		t.Fatal(err)
	}

	t.Logf("Unmarshaled string: %v", name)

}

func TestUnmarshal_compound(t *testing.T) {
	var data = []byte{
		0x0a, 0x00, 0x0b, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20,
		0x77, 0x6f, 0x72, 0x6c, 0x64, 0x08, 0x00, 0x04, 0x6e,
		0x61, 0x6d, 0x65, 0x00, 0x09, 0x42, 0x61, 0x6e, 0x61,
		0x6e, 0x72, 0x61, 0x6d, 0x61, 0x00,
	}

	var compound struct {
		Name string `nbt:"name"`
	}
	if err := Unmarshal(data, &compound); err != nil {
		t.Fatal(err)
	}

	t.Logf("Unmarshaled compound: %v", compound)
}

func TestUnmarshal_servers(t *testing.T) {
	data, err := ioutil.ReadFile("./servers.dat")
	if err != nil {
		t.Fatal(err)
	}

	type ServersData struct {
		Servers []struct {
			IP   string `nbt:"ip"`
			Name string `nbt:"name"`
			Icon string `nbt:"icon"`
		} `nbt:"servers"`
	}

	var servers ServersData
	if err := Unmarshal(data, &servers); err != nil {
		t.Fatal(err)
	}

	t.Logf("Unmarshaled servers: %v", servers)
}

func TestUnmarshal_level(t *testing.T) {
	data, err := ioutil.ReadFile("./level.dat")
	if err != nil {
		t.Fatal(err)
	}

	var levelData interface{}
	if err := Unmarshal(data, &levelData); err != nil {
		t.Fatal(err)
	}

	t.Logf("Unmarshaled level data: %v", levelData)
}
