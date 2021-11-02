package main

import (
	"context"
	pb "example/Mini_Project_2_Chitty-Chat/chat"
	"net"

	"google.golang.org/grpc"

	"fmt"
	"log"
	"strconv"

	cd "example/Mini_Project_2_Chitty-Chat/server/database"
	lamport "example/Mini_Project_2_Chitty-Chat/timestamp"
)

const (
	port = ":8008"
)

type ChatServiceServer struct {
	pb.UnimplementedChatServiceServer
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterChatServiceServer(s, &ChatServiceServer{})
	fmt.Printf("server listening at %v\n", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

var chatData = cd.NewChatDatabase()
var localLamport = lamport.NewClock()

func (s *ChatServiceServer) Connect(ctx context.Context, rec *pb.Request) (*pb.Response, error) {
	localLamport.Witness(rec.Timestamp)

	success := chatData.AddUser(rec.User)
	if !success {
		status := "Join Failed: Username already taken! - Try another Username"
		resp := &pb.Response{Status: status, Timestamp: localLamport.Increment()}
		return resp, nil
	}

	defer chatData.InsertMessage(&pb.Message{User: rec.User, Text: rec.User.Username + " has joined the chat!", Timestamp: localLamport.GetTimestamp()})

	status := "for user " + rec.User.Username + " - Join Successful"
	resp := &pb.Response{Status: status, Timestamp: localLamport.Increment()}
	fmt.Println(resp)
	return resp, nil
}

func (s *ChatServiceServer) Disconnect(ctx context.Context, rec *pb.Request) (*pb.Response, error) {
	localLamport.Witness(rec.Timestamp)

	chatData.RemoveUser(rec.User)

	defer chatData.InsertMessage(&pb.Message{User: rec.User, Text: rec.User.Username + " has left the chat!", Timestamp: localLamport.GetTimestamp()})

	status := "for user " + rec.User.Username + " - Left Successful"
	resp := &pb.Response{Status: status, Timestamp: localLamport.Increment()}
	fmt.Println(resp)
	return resp, nil
}

func (s *ChatServiceServer) Publish(ctx context.Context, msg *pb.Message) (*pb.Response, error) {
	localLamport.Witness(msg.Timestamp)

	chatData.InsertMessage(msg)

	status := "for user " + msg.User.Username + " - Message Recieved: " + msg.Text
	resp := &pb.Response{Status: status, Timestamp: localLamport.Increment()}
	fmt.Println(resp)
	return resp, nil
}

func (s *ChatServiceServer) Listen(ctx context.Context, rec *pb.Request) (*pb.Message, error) {
	localLamport.Witness(rec.Timestamp)

	for {
		possibleMessage := chatData.PopMessage(rec.User)
		if possibleMessage != nil {
			possibleMessage.Timestamp = localLamport.Increment()
			status := "Timestamp:" + strconv.FormatInt(localLamport.GetTimestamp(), 10) +
				"  Status:\"for user " + rec.User.Username +
				" - Accesing Message: " + possibleMessage.Text + "\""
			fmt.Println(status)
			return possibleMessage, nil
		}
	}
}
