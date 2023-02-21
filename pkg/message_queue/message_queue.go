package messagequeue

import (
	ormdomain "github.com/QingYunTasha/go-Chatting/domain/infra/orm"
)

type MessageQueue struct {
	Queue chan ormdomain.GroupMessage
}

func NewMessageQueue(size int) *MessageQueue {
	return &MessageQueue{
		Queue: make(chan ormdomain.GroupMessage, size),
	}
}

func (q *MessageQueue) Pop() ormdomain.GroupMessage {
	return <-q.Queue
}

func (q *MessageQueue) Push(msg ormdomain.GroupMessage) {
	q.Queue <- msg
}

func (q *MessageQueue) Len() int {
	return len(q.Queue)
}
