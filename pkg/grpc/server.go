package grpc

import (
	"context"
	"fmt"
	"net"

	"go-chatting/internal/service/chat"

	"google.golang.org/grpc"
)

// Server defines the gRPC server.
type Server struct {
	grpcServer  *grpc.Server
	chatService chat.Service
}

// NewServer creates a new instance of the gRPC server.
func NewServer(port string, chatService chat.Service) *Server {
	s := &Server{
		grpcServer:  grpc.NewServer(),
		chatService: chatService,
	}

	RegisterChatServiceServer(s.grpcServer, s)

	return s
}

// Start starts the gRPC server.
func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		return err
	}

	fmt.Printf("gRPC server listening on %s\n", lis.Addr())

	if err := s.grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}

// SendMessage is the implementation of the ChatService SendMessage gRPC method.
func (s *Server) SendMessage(ctx context.Context, req *SendMessageRequest) (*SendMessageResponse, error) {
	message, err := s.chatService.SendMessage(ctx, req.UserId, req.Content)
	if err != nil {
		return nil, err
	}

	res := &SendMessageResponse{
		Id:        message.ID,
		UserId:    message.UserID,
		Content:   message.Content,
		Timestamp: message.Timestamp,
	}

	return res, nil
}
