package v2

type Encoder struct {
	threshold int
}

func (e Encoder) Marshal(v interface{}) error {
	buf := Buffer{}

	// Write
}
