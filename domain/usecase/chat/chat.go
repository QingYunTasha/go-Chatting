package domain

import (
	"context"
	"net"
	"net/http"

	ormdomain "github.com/QingYunTasha/go-Chatting/domain/infra/orm"
	"github.com/gin-gonic/gin"
)

type ChatUsecaseRepository interface {
	UserConnect(ctx context.Context, r *http.Request, w gin.ResponseWriter, userID uint32, PrivateChatLastID map[uint32]uint32, GroupChatLastID map[uint32]uint32) error
	ProcessReadWebSocket(conn net.Conn) error
	ProcessHistoryCatchUp(userID uint32, PrivateChatLastID map[uint32]uint32, GroupChatLastID map[uint32]uint32) error
	ProcessPrivateMessage(privateMessage ormdomain.PrivateMessage) error
	ProcessGroupMessage(groupMessage ormdomain.GroupMessage) error
	ProcessStatus(userID uint32, status ormdomain.UserStatus) error
	getUserConnection(userID uint32) error
	publishGroupMessage(groupMessage ormdomain.GroupMessage) error
}
