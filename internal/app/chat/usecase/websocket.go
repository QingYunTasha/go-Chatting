package usecase

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"

	ormdomain "github.com/QingYunTasha/go-Chatting/domain/infra/orm"
	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type MessageType string

const (
	PrivateChatType MessageType = "PrivateChatType"
	StatusType      MessageType = "StatusType"
)

type ChatPayLoad struct {
	Type           MessageType
	PrivateMessage ormdomain.PrivateMessage
}

func HandleReadWebSocket(conn net.Conn, messageChan chan []byte) {
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

func UserConnect(ctx context.Context, r *http.Request, w gin.ResponseWriter, user1Id uint32, user2ID uint32) error {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return err
	}
	defer conn.Close()

	messageChan := make(chan []byte)
	defer close(messageChan)
	go HandleReadWebSocket(conn, messageChan)

	err = ProcessHistoryCatchUp(user1Id, user2ID)
	if err != nil {
		log.Println("History catch-up error:", err)
		return err
	}

	var chatPayload ChatPayLoad
	for msg := range messageChan {
		err := json.Unmarshal(msg, &chatPayload)
		if err != nil {
			log.Println("webSocket message parse error:", err)
			return err
		}

		switch chatPayload.Type {
		case PrivateChatType:
			ProcessPrivateChat(chatPayload)
		case StatusType:
			ProcessStatus(chatPayload)
		}
	}
	return nil
}

func ProcessHistoryCatchUp(user1Id, user2ID uint32) error {
	messages, err := GetUnsendMessage(user1Id, user2ID)
	if err != nil {
		log.Println("get unsend message error:", err)
		return err
	}
	err = SendUnsendMessage(messages)
	if err != nil {
		log.Println("send message error:", err)
		return err
	}

	return nil
}

func GetUnsendMessage(user1Id, user2ID uint32) ([]ormdomain.PrivateMessage, error) {
	return nil, nil
}

func SendUnsendMessage([]ormdomain.PrivateMessage) error {
	return nil
}

func ProcessPrivateChat(chatPayLoad ChatPayLoad) {

}

func ProcessStatus(chatPayLoad ChatPayLoad) {

}
