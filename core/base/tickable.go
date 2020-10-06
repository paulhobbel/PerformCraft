package base

type Tickable interface {
	Tick(deltaTicks int)
}
