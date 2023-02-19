package delivery

import (
	"io"

	ormdomain "github.com/QingYunTasha/go-Chatting/domain/infra/orm"
	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	UnimplementedMessageQueueServer
	Mq ormdomain.MessageQueue
}

func (s *Server) PublishMessage(stream MessageQueue_PublishMessageServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&empty.Empty{})
		}
		if err != nil {
			return err
		}

		s.Mq.Queue <- ormdomain.Message{
			ID:        msg.ID,
			SenderID:  msg.SenderID,
			Timestamp: msg.Timestamp.AsTime(),
			Group:     msg.Group,
			Content:   msg.Content,
		}
	}
}
func (s *Server) ConsumeMessage(_ *empty.Empty, stream MessageQueue_ConsumeMessageServer) error {
	for {
		msg := <-s.Mq.Queue

		if err := stream.Send(&Message{
			ID:        msg.ID,
			SenderID:  msg.SenderID,
			Timestamp: timestamppb.New(msg.Timestamp),
			Group:     msg.Group,
			Content:   msg.Content,
		}); err != nil {
			return err
		}
	}
}
