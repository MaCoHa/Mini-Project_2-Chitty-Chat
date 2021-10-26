package main

import (
	"context"
	"example/Mini_Project_2_Chitty-Chat/chat"
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	port = ":8008"
)

type Server struct {
	chat.UnimplementedChatServiceServer
}

func (s *Server) Broadcast(ctx context.Context, in *chat.Msg) (*chat.Response, error) {
	log.Printf("Received Broadcast request")
	return &chat.Response{Message: "Yo"}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	chat.RegisterChatServiceServer(s, &Server{})
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
