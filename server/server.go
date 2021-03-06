package main

import (
	"context"
	pb "example/Mini_Project_2_Chitty-Chat/chat"
	"net"

	"google.golang.org/grpc"

	"os"

	"fmt"
	"log"

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
	//Setup the file for log outputs
	LOG_FILE := "./systemlogs/server.log"
	// open log file

	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Panic(err)
	}

	log.SetOutput(logFile)
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
	defer logFile.Close()
}

var chatData = cd.NewChatDatabase()
var localLamport = lamport.NewClock()

func (s *ChatServiceServer) Connect(ctx context.Context, rec *pb.Request) (*pb.Response, error) {
	localLamport.Witness(rec.Timestamp)
	rec.Timestamp = localLamport.GetTimestamp()
	log.Print(rec)
	log.Println(" - wants to connect")

	success := chatData.AddUser(rec.User)
	if !success {
		status := "Join Failed: Username already taken! - Try another Username"
		resp := &pb.Response{Status: status, Timestamp: localLamport.Increment()}
		log.Println(resp)
		return resp, nil
	}

	status := rec.User.Username + " - Join Successful"
	resp := &pb.Response{Status: status, Timestamp: localLamport.Increment()}
	log.Println(resp)

	defer chatData.InsertMessage(&pb.Message{User: rec.User, Text: rec.User.Username + " has joined the chat!", Timestamp: localLamport.GetTimestamp()})
	return resp, nil
}

func (s *ChatServiceServer) Disconnect(ctx context.Context, rec *pb.Request) (*pb.Response, error) {
	localLamport.Witness(rec.Timestamp)
	rec.Timestamp = localLamport.GetTimestamp()
	log.Print(rec)
	log.Println(" - wants to disconnect")

	chatData.RemoveUser(rec.User)

	status := rec.User.Username + " - Left Successful"
	resp := &pb.Response{Status: status, Timestamp: localLamport.Increment()}
	log.Println(resp)

	defer chatData.InsertMessage(&pb.Message{User: rec.User, Text: rec.User.Username + " has left the chat!", Timestamp: localLamport.GetTimestamp()})
	return resp, nil
}

func (s *ChatServiceServer) Publish(ctx context.Context, msg *pb.Message) (*pb.Response, error) {
	localLamport.Witness(msg.Timestamp)
	msg.Timestamp = localLamport.GetTimestamp()
	log.Print(msg)
	log.Println(" - Message Recieved")

	chatData.InsertMessage(msg)

	status := msg.User.Username + " - Broadcasting Message: " + msg.Text
	resp := &pb.Response{Status: status, Timestamp: localLamport.Increment()}
	log.Println(resp)
	return resp, nil
}

func (s *ChatServiceServer) Listen(ctx context.Context, rec *pb.Request) (*pb.Message, error) {
	localLamport.Witness(rec.Timestamp)
	rec.Timestamp = localLamport.GetTimestamp()
	log.Print(rec)
	log.Println(" - is listening")

	for {
		possibleMessage := chatData.PopMessage(rec.User)
		if possibleMessage != nil {
			possibleMessage.Timestamp = localLamport.Increment()
			log.Print(possibleMessage)
			log.Println(" - retrieved by " + rec.User.Username)
			return possibleMessage, nil
		}
	}
}
