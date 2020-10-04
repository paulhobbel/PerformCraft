package entity

type Sender interface {
	SendMessage(message ...interface{})
}
