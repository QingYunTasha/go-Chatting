package domain

type MessageQueue struct {
	Queue chan GroupMessage
}

func NewMessageQueue(size int) *MessageQueue {
	return &MessageQueue{
		Queue: make(chan GroupMessage, size),
	}
}
