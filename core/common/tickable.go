package common

type Tickable interface {
	Tick(deltaTicks int)
}
