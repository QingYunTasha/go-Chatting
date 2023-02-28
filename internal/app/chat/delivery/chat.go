package delivery

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"strconv"

	chatdeliverydomain "github.com/QingYunTasha/go-Chatting/domain/delivery/chat"
	ormdomain "github.com/QingYunTasha/go-Chatting/domain/infra/orm"
	chatusecasedomain "github.com/QingYunTasha/go-Chatting/domain/usecase/chat"
	ormfactory "github.com/QingYunTasha/go-Chatting/internal/infra/orm/factory"
	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type PostPayload struct {
	UserID            uint32
	PrivateChatLastID map[uint32]uint32
	GroupChatLastID   map[uint32]uint32
}

type ChatHandler struct {
	OrmRepo     ormfactory.OrmRepository
	ChatUsecase chatusecasedomain.ChatUsecaseRepository
	ConnMgr     chatdeliverydomain.ConnectionManager
}

func NewChatHandler(chatUsecase chatusecasedomain.ChatUsecaseRepository) *ChatHandler {
	return &ChatHandler{
		ChatUsecase: chatUsecase,
	}
}

func (h *ChatHandler) UserConnect(c *gin.Context) {
	pl := PostPayload{}
	if err := c.BindJSON(&pl); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	conn, _, _, err := ws.UpgradeHTTP(c.Request, c.Writer)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	h.ConnMgr.AddConnection(pl.UserID, &conn)
	defer h.ConnMgr.RemoveConnection(pl.UserID)

	stopCh := make(chan bool)
	errCh := make(chan error)
	go h.ProcessHistoryCatchUp(conn, pl.UserID, pl.PrivateChatLastID, pl.GroupChatLastID, stopCh, errCh)
	go h.processReadWebSocket(conn, stopCh, errCh)

	err = <-errCh
	log.Println(err)
	stopCh <- true
}

func (h *ChatHandler) processReadWebSocket(conn net.Conn, stopCh chan bool, errCh chan error) {
	for {
		select {
		case <-stopCh:
			return
		default:
			msg, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				errCh <- err
				return
			}

			if op == ws.OpClose {
				errCh <- errors.New("WebSocket closed")
				return
			}

			var chatPayload chatdeliverydomain.ChatPayLoad
			err = json.Unmarshal(msg, &chatPayload)
			if err != nil {
				errCh <- err
				return
			}

			switch chatPayload.Type {
			case chatdeliverydomain.PrivateChatType:
				err = h.ProcessPrivateMessage(chatPayload.PrivateMessage)
			case chatdeliverydomain.GroupChatType:
				err = h.ProcessGroupMessage(chatPayload.GroupMessage)
			case chatdeliverydomain.StatusType:
				err = h.ProcessStatus(chatPayload.SenderID, chatPayload.Status)
			default:
				log.Print("unknown chat type:", chatPayload.Type)
			}
			if err != nil {
				errCh <- err
			}
		}
	}
}

func (h *ChatHandler) ProcessHistoryCatchUp(conn net.Conn, userID uint32, PrivateChatLastID map[uint32]uint32, GroupChatLastID map[uint32]uint32, stopCh chan bool, errCh chan error) {
	for user2ID, lastID := range PrivateChatLastID {
		messages, err := h.OrmRepo.PrivateMessage.GetBetweenUsersAfterMsgID(userID, user2ID, lastID)
		if err != nil {
			errCh <- err
			return
		}
		err = h.SendUnsendMessage(conn, messages)
		if err != nil {
			errCh <- err
			return
		}
	}

}

func (h *ChatHandler) ProcessPrivateMessage(msg ormdomain.PrivateMessage) error {
	if err := h.OrmRepo.PrivateMessage.Create(&msg); err != nil {
		return err
	}

	receiverID := msg.ReceiverID
	conn, ok := h.ConnMgr.GetConnection(receiverID)
	if !ok {
		return errors.New("connection not found")
	}

	payLoadBytes, err := json.Marshal(&chatdeliverydomain.ChatPayLoad{
		Type:           chatdeliverydomain.PrivateChatType,
		PrivateMessage: msg,
	})
	if err != nil {
		return err
	}

	err = wsutil.WriteServerMessage(*conn, ws.OpText, payLoadBytes)
	return err
}

func (h *ChatHandler) ProcessGroupMessage(msg ormdomain.GroupMessage) error {
	if err := h.OrmRepo.GroupMessage.Create(&msg); err != nil {
		return err
	}

	if err := h.PublishGroupMessage(msg); err != nil {
		return err
	}
	return nil
}

func (h *ChatHandler) ProcessStatus(userID uint32, status ormdomain.UserStatus) error {
	err := h.OrmRepo.User.Update(userID, &ormdomain.User{Status: status})
	if err != nil {
		return err
	}

	err = h.SendUserStatus(userID, status)

	return err
}

func (h *ChatHandler) SendUnsendMessage(conn net.Conn, messages []ormdomain.PrivateMessage) error {
	return nil
}

func (h *ChatHandler) PublishGroupMessage(groupMessage ormdomain.GroupMessage) error {
	return nil
}

func (h *ChatHandler) SendUserStatus(userID uint32, status ormdomain.UserStatus) error {
	return nil
}

func stringToUint32(str string) (uint32, error) {
	ui64, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, err
	}

	ui := uint32(ui64)

	return ui, nil
}
