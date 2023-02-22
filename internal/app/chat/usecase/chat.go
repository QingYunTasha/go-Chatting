package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"sync"

	ormdomain "github.com/QingYunTasha/go-Chatting/domain/infra/orm"
	chatdomain "github.com/QingYunTasha/go-Chatting/domain/usecase/chat"
	orm "github.com/QingYunTasha/go-Chatting/internal/infra/orm/factory"
	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type ChatUsecase struct {
	OrmRepo orm.OrmRepository
	ConnMgr ConnectionManager
}

type ConnectionManager struct {
	connections map[uint32]*net.Conn
	mutex       sync.RWMutex
}

func (cm *ConnectionManager) AddConnection(userID uint32, conn *net.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.connections[userID] = conn
}

func (cm *ConnectionManager) GetConnection(userID uint32) (*net.Conn, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	conn, ok := cm.connections[userID]
	return conn, ok
}

func (cm *ConnectionManager) RemoveConnection(userID uint32) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	delete(cm.connections, userID)
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
		case chatdomain.StatusType:
			err = u.ProcessStatus(chatPayload)
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

func (u *ChatUsecase) ProcessStatus(chatPayLoad chatdomain.ChatPayLoad) error {
	userID := chatPayLoad.SenderID
	err := u.OrmRepo.User.Update(userID, &ormdomain.User{Status: chatPayLoad.Status})
	if err != nil {
		return err
	}

	return nil
}

func NewChatUsecase() *ChatUsecase {
	return &ChatUsecase{}
}

func (u *ChatUsecase) SendMessage(senderID, receiverID, channelType, message string) error {
	return nil
}

func (u *ChatUsecase) GetMessage(senderID, receiverID, channelType, preMessageID string) error {
	return nil
}

func GetStatus(senderID, receiverID, channelType, status string) error {
	return nil
}

func SendStatus(senderID, receiverID, channelType, status string) error {
	return nil
}
