package domain

import (
	dbdomain "go-Chatting/domain/infra/database"
)

type MessageType string

const (
	PrivateChatType MessageType = "PrivateChatType"
	GroupChatType   MessageType = "GroupChatType"
	StatusType      MessageType = "StatusType"
)

type ChatPayLoad struct {
	Type           MessageType
	SenderID       uint32
	Status         dbdomain.UserStatus
	PrivateMessage dbdomain.PrivateMessage
	GroupMessage   dbdomain.GroupMessage
}
