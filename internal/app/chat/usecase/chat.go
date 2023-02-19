package usecase

import (
	"context"
	"log"

	ormdomain "github.com/QingYunTasha/go-Chatting/domain/infra/orm"
	msg_queue "github.com/QingYunTasha/go-Chatting/internal/app/msg_queue/delivery"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NotificationUsecase struct {
	Stream msg_queue.MessageQueue_PublishMessageClient
}

func NewNotificationUsecase() *NotificationUsecase {
	var opts []grpc.DialOption
	conn, err := grpc.Dial(":50051", opts)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	client := msg_queue.NewMessageQueueClient(conn)
	stream, err := client.PublishMessage(context.TODO())
	if err != nil {
		log.Fatalf("client.PublishMessage failed: %v", err)
	}

	return &NotificationUsecase{
		Stream: stream,
	}
}

func (u *NotificationUsecase) PublishMessage(ctx context.Context, msg ormdomain.GroupMessage) error {
	if err := u.Stream.Send(&msg_queue.Message{
		ID:        msg.ID,
		SenderID:  msg.SenderID,
		Group:     msg.Group,
		Timestamp: timestamppb.New(msg.Timestamp),
		Content:   msg.Content,
	}); err != nil {
		return err
	}

	return nil
}

type ChatUsecase struct {
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
