package main

import (
	"context"
	pb "example/Mini_Project_2_Chitty-Chat/chat"
	"log"
	"net"
	"strconv"

	//"sync"

	cd "example/Mini_Project_2_Chitty-Chat/server/database"
	lamport "example/Mini_Project_2_Chitty-Chat/timestamp"

	"google.golang.org/grpc"
)

const (
	port = ":8008"
)

var chatData = cd.NewChatDatabase()
var lamp = lamport.NewClock()

type ChatServiceServer struct {
	pb.UnimplementedChatServiceServer
}

func (s *ChatServiceServer) Connect(ctx context.Context, rec *pb.Request) (*pb.Response, error) {
	lamp.Witness(rec.Timestamp)

	success := chatData.AddUser(rec.User)
	if !success {
		failMessage := "Join Failed: Username already taken! - Try another Username"
		resp := &pb.Response{Status: failMessage, Timestamp: lamp.Increment()}
		return resp, nil
	}

	defer chatData.InsertMessage(&pb.Message{User: rec.User, Text: rec.User.Username + " has joined the chat!", Timestamp: lamp.GetTimestamp()})

	resp := &pb.Response{Status: "Join Successful for user " + rec.User.Username, Timestamp: lamp.Increment()}
	log.Println(resp)
	return resp, nil
}

func (s *ChatServiceServer) Disconnect(ctx context.Context, rec *pb.Request) (*pb.Response, error) {
	lamp.Witness(rec.Timestamp)

	chatData.RemoveUser(rec.User)

	defer chatData.InsertMessage(&pb.Message{User: rec.User, Text: rec.User.Username + " has left the chat!", Timestamp: lamp.GetTimestamp()})

	resp := &pb.Response{Status: "Left Successful for user " + rec.User.Username, Timestamp: lamp.Increment()}
	log.Println(resp)
	return resp, nil
}

func (s *ChatServiceServer) Publish(ctx context.Context, msg *pb.Message) (*pb.Response, error) {
	lamp.Witness(msg.Timestamp)

	Broadcast(msg)

	return &pb.Response{Status: "Message Recieved", Timestamp: lamp.Increment()}, nil
}

func Broadcast(msg *pb.Message) {
	log.Println("Status at time " + strconv.FormatInt(lamp.GetTimestamp(), 10) + ": Broadcasting: " + msg.Text)
	chatData.InsertMessage(msg)
}

func (s *ChatServiceServer) Listen(ctx context.Context, rec *pb.Request) (*pb.Message, error) {
	lamp.Witness(rec.Timestamp)

	for {
		possibleMessage := chatData.PopMessage(rec.User)
		if possibleMessage != nil {
			log.Println("Status at time " + strconv.FormatInt(lamp.GetTimestamp(), 10) + ": Accesing message: " + possibleMessage.Text + " - for user: " + rec.User.Username)
			possibleMessage.Timestamp = lamp.Increment()
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
