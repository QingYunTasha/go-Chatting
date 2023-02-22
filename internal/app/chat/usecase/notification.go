package usecase

import (
	"context"
	"log"

	msgqueuedomain "github.com/QingYunTasha/go-Chatting/domain/delivery/message_queue"
	ormdomain "github.com/QingYunTasha/go-Chatting/domain/infra/orm"
	orm "github.com/QingYunTasha/go-Chatting/internal/infra/orm/factory"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NotificationUsecase struct {
	Stream  msgqueuedomain.MessageQueue_PublishMessageClient
	OrmRepo orm.OrmRepository
}

func NewNotificationUsecase() *NotificationUsecase {
	var opts grpc.DialOption
	conn, err := grpc.Dial(":50051", opts)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	client := msgqueuedomain.NewMessageQueueClient(conn)
	ctx := context.Background()
	stream, err := client.PublishMessage(ctx)
	if err != nil {
		log.Fatalf("client.PublishMessage failed: %v", err)
	}

	return &NotificationUsecase{
		Stream: stream,
	}
}

func (u *NotificationUsecase) PublishMessage(ctx context.Context, msg ormdomain.GroupMessage) error {
	if err := u.Stream.Send(&msgqueuedomain.Message{
		ID:        msg.ID,
		SenderID:  msg.SenderID,
		GroupID:   msg.GroupID,
		Timestamp: timestamppb.New(msg.Timestamp),
		Content:   msg.Content,
	}); err != nil {
		return err
	}

	return nil
}
