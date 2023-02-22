package main

import (
	"log"
	"net"

	msgqueuedomain "github.com/QingYunTasha/go-Chatting/domain/delivery/message_queue"

	msgqueuepkg "github.com/QingYunTasha/go-Chatting/pkg/message_queue"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	msgqueuedomain.RegisterMessageQueueServer(s, &msgqueuedomain.Server{
		Mq: msgqueuepkg.NewMessageQueue(32),
	})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
