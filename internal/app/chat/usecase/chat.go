package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"

	ormdomain "github.com/QingYunTasha/go-Chatting/domain/infra/orm"
	chatdomain "github.com/QingYunTasha/go-Chatting/domain/usecase/chat"
	ormfactory "github.com/QingYunTasha/go-Chatting/internal/infra/orm/factory"
	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type ChatUsecase struct {
	OrmRepo ormfactory.OrmRepository
	ConnMgr ConnectionManager
}

func NewChatUsecase(ormRepo ormfactory.OrmRepository) *ChatUsecase {
	return &ChatUsecase{}
}

func (u *ChatUsecase) HandleReadWebSocket(conn net.Conn, messageChan chan []byte) {
	for {
		msg, op, err := wsutil.ReadClientData(conn)
		if err != nil {
			log.Println("WebSocket read error:", err)
		}

		switch op {
		case ws.OpText:
			log.Println("Text message received:", string(msg))
			messageChan <- msg
		case ws.OpBinary:
			log.Println("Binary message received:", msg)
			messageChan <- msg
		case ws.OpClose:
			log.Println("WebSocket closed")
			return
		}
	}
}

func (u *ChatUsecase) UserConnect(ctx context.Context, r *http.Request, w gin.ResponseWriter, user1ID, user2ID, lastID uint32) error {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return err
	}
	defer conn.Close()

	messageChan := make(chan []byte)
	defer close(messageChan)
	go u.HandleReadWebSocket(conn, messageChan)

	err = u.ProcessHistoryCatchUp(conn, user1ID, user2ID, lastID)
	if err != nil {
		log.Println("History catch-up error:", err)
		return err
	}

	var chatPayload chatdomain.ChatPayLoad
	for msg := range messageChan {
		err := json.Unmarshal(msg, &chatPayload)
		if err != nil {
			log.Println("Unmarshal webSocket message error:", err)
			return err
		}

		switch chatPayload.Type {
		case chatdomain.PrivateChatType:
			err = u.ProcessPrivateChat(chatPayload)
		case chatdomain.GroupChatType:
			err = u.ProcessGroupChat(chatPayload)
		case chatdomain.StatusType:
			err = u.ProcessStatus(chatPayload)
		default:
			log.Print("unknown chat type:", chatPayload.Type)
		}
		if err != nil {
			log.Println("Process webSocket message error:", err)
			return err
		}

	}
	return nil
}

func (u *ChatUsecase) ProcessHistoryCatchUp(conn net.Conn, user1ID, user2ID, lastID uint32) error {
	messages, err := u.OrmRepo.PrivateMessage.GetBetweenUsersAfterMsgID(user1ID, user2ID, lastID)
	if err != nil {
		log.Println("get unsend message error:", err)
		return err
	}
	err = u.SendUnsendMessage(conn, messages)
	if err != nil {
		log.Println("send message error:", err)
		return err
	}

	return nil
}

func (u *ChatUsecase) SendUnsendMessage(conn net.Conn, messages []ormdomain.PrivateMessage) error {
	for _, msg := range messages {
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			log.Fatal(err)
		}
		wsutil.WriteServerMessage(conn, ws.OpText, msgBytes)
	}
	return nil
}

func (u *ChatUsecase) ProcessPrivateChat(chatPayLoad chatdomain.ChatPayLoad) error {
	if err := u.OrmRepo.PrivateMessage.Create(&chatPayLoad.PrivateMessage); err != nil {
		return err
	}

	receiverID := chatPayLoad.PrivateMessage.ReceiverID
	conn, ok := u.ConnMgr.GetConnection(receiverID)
	if !ok {
		return errors.New("connection not found")
	}

	payLoadBytes, err := json.Marshal(chatPayLoad)
	if err != nil {
		return err
	}

	err = wsutil.WriteServerMessage(*conn, ws.OpText, payLoadBytes)
	return err
}

func (u *ChatUsecase) ProcessGroupChat(chatPayLoad chatdomain.ChatPayLoad) error {
	if err := u.OrmRepo.GroupMessage.Create(&chatPayLoad.GroupMessage); err != nil {
		return err
	}

	if err := SendGroupMessage(chatPayLoad.GroupMessage); err != nil {
		return err
	}
	return nil
}

func (u *ChatUsecase) ProcessStatus(chatPayLoad chatdomain.ChatPayLoad) error {
	userID := chatPayLoad.SenderID
	err := u.OrmRepo.User.Update(userID, &ormdomain.User{Status: chatPayLoad.Status})
	if err != nil {
		return err
	}

	SendUserStatus(userID, chatPayLoad.Status)

	return nil
}

func SendUserStatus(userID uint32, status ormdomain.UserStatus) error {
	// TODO: SendStatusToPrivate()
	SendStatusToFriends()
	SendStatusToGroup()
	return nil
}

func SendGroupMessage(ormdomain.GroupMessage) error {
	GetSenderGroups()

	GetGroupOnlineUsers()
	SendMessageToUsers()

	return nil
}

func SendStatusToFriends() {}
func SendStatusToGroup()   {}
func GetSenderGroups()     {}
func GetGroupOnlineUsers() {}
func SendMessageToUsers()  {}
