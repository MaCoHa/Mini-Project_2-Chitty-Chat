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

var messageList []string = make([]string, 1)
var userToMesageMap map[*pb.User]*pb.Message = make(map[*pb.User]*pb.Message)

type ChatServiceServer struct {
	pb.UnimplementedChatServiceServer
}

func (s *ChatServiceServer) Connect(ctx context.Context, user *pb.User) (*pb.Response, error) {
	if messageList == nil {
		log.Println("this list is nil")
	}
	if userToMesageMap == nil {
		log.Println("the map is nil")
	}
	messageList = append(messageList, user.Username+" is connected")
	userToMesageMap[user] = nil
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
	messageList = append(messageList, msg.Text)
	userToMesageMap[msg.User] = msg
	return &pb.Response{Status: "Message Recieved"}, nil
}

func (s *ChatServiceServer) Listen(ctx context.Context, user *pb.User) (*pb.Message, error) {
	var newMsg *pb.Message

	for {
		if userToMesageMap[user] != nil {
			newMsg = userToMesageMap[user]
			userToMesageMap[user] = nil
			break
		}
	}

	/*for {
		if len(messageList) > 0 {
			newMsg = &pb.Message{Text: messageList[0]}
			messageList = messageList[1:]
			break
		}
	}*/

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
}
