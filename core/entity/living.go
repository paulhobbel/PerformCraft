package entity

type LivingEntity interface {
	Entity

	GetHealth() float64
	SetHealth(health float64)
}
