package main

import (
	"context"
	pb "example/Mini_Project_2_Chitty-Chat/chat"
	"log"
	"net"

	//"sync"

	cd "example/Mini_Project_2_Chitty-Chat/server/database"

	"google.golang.org/grpc"
)

const (
	port = ":8008"
)

var chatData = cd.NewChatDatabase()

type ChatServiceServer struct {
	pb.UnimplementedChatServiceServer
}

func (s *ChatServiceServer) Connect(ctx context.Context, user *pb.User) (*pb.Response, error) {
	success := chatData.AddUser(user)
	if !success {
		return &pb.Response{Status: "Join Failed: Username already taken! - Try another Username"}, nil
	}

	chatData.InsertMessage(&pb.Message{User: user, Text: user.Username + " has joined the chat!"})
	log.Println("Status: " + user.Username + " has joined the chat!")
	return &pb.Response{Status: "Join Successful"}, nil
}

func (s *ChatServiceServer) Disconnect(ctx context.Context, user *pb.User) (*pb.Response, error) {
	chatData.RemoveUser(user)
	chatData.InsertMessage(&pb.Message{User: user, Text: user.Username + " has left the chat!"})
	log.Println("Status: " + user.Username + " has left the chat!")
	return &pb.Response{Status: "Left Successful"}, nil
}

func (s *ChatServiceServer) Publish(ctx context.Context, msg *pb.Message) (*pb.Response, error) {
	resp, err := s.Broadcast(ctx, msg)
	if err != nil {
		log.Fatalf("Could not send message: %v", err)
	}
	return resp, err
}

func (s *ChatServiceServer) Broadcast(ctx context.Context, msg *pb.Message) (*pb.Response, error) {
	log.Println("Status: Broadcasting: " + msg.Text)
	chatData.InsertMessage(msg)
	return &pb.Response{Status: "Message Recieved"}, nil
}

func (s *ChatServiceServer) Listen(ctx context.Context, user *pb.User) (*pb.Message, error) {
	for {
		possibleMessage := chatData.PopMessage(user)
		if possibleMessage != nil {
			return possibleMessage, nil
		}
	}
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
}
