package main

import (
	"log"
	"net"

	ormdomain "github.com/QingYunTasha/go-Chatting/domain/infra/orm"
	msg_queue "github.com/QingYunTasha/go-Chatting/internal/app/msg_queue/delivery"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	msg_queue.RegisterMessageQueueServer(s, &msg_queue.Server{
		Mq: ormdomain.NewMessageQueue(32),
	})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
