package domain

import (
	ormdomain "github.com/QingYunTasha/go-Chatting/domain/infra/orm"
)

type MessageType string

const (
	PrivateChatType MessageType = "PrivateChatType"
	StatusType      MessageType = "StatusType"
)

type ChatPayLoad struct {
	Type           MessageType
	SenderID       uint32
	Status         ormdomain.UserStatus
	PrivateMessage ormdomain.PrivateMessage
	GroupMessage   ormdomain.GroupMessage
}
