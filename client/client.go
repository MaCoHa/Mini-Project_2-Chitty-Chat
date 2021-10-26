package main

import (
	pb "example/Mini_Project_2_Chitty-Chat/chat"
	"log"

	"google.golang.org/grpc"
)

const (
	serverAddr = "localhost:8008"
)

func main() {
	//var opts []grpc.DialOption
	// Set up a connection to the server.
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	client := pb.NewChatServiceClient(conn)

	SendPublish(client)
}

func SendPublish(client pb.ChatServiceClient) {
	chat := pb.Msg{
		Message: "TestMessage",
	}
	log.Printf("Message: %s", chat.Message)
}
