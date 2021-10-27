package main

import (
	"context"
	pb "example/Mini_Project_2_Chitty-Chat/chat"
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	port = ":8008"
)

var userList *pb.UserList

type ChatServiceServer struct {
	pb.UnimplementedChatServiceServer
}

func (s *ChatServiceServer) Connect(ctx context.Context, user *pb.User) (*pb.Response, error) {
	userList.Users = append(userList.Users, user)
	userList.MessageMap[user.Username] = &pb.Message{User: user, Text: ""}
	return &pb.Response{Status: "Join Successful"}, nil
}

func (s *ChatServiceServer) Publish(ctx context.Context, msg *pb.Message) (*pb.Response, error) {
	resp, err := s.Broadcast(ctx, msg)
	if err != nil {
		log.Fatalf("Could not send message: %v", err)
	}
	return resp, nil
}

func (s *ChatServiceServer) Broadcast(ctx context.Context, msg *pb.Message) (*pb.Response, error) {
	log.Printf("Msg:" + msg.Text)
	userList.MessageMap[msg.User.Username] = msg
	return &pb.Response{Status: "Message Received"}, nil
}

func (s *ChatServiceServer) Listen(ctx context.Context, user *pb.User) (*pb.Message, error) {
	var newMsg *pb.Message

	for userName, msg := range userList.MessageMap {
		if msg.Text != "" {
			newMsg = msg
			userList.MessageMap[userName] = &pb.Message{User: &pb.User{Username: userName}, Text: ""}
			break
		}
	}

	return newMsg, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterChatServiceServer(s, &ChatServiceServer{})
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	userList = &pb.UserList{Users: make([]*pb.User, 5), MessageMap: make(map[string]*pb.Message)}
}
