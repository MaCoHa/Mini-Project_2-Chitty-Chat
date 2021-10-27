package main

import (
	"bufio"
	"context"
	pb "example/Mini_Project_2_Chitty-Chat/chat"
	"log"
	"os"
	"time"
	"strings"

	"google.golang.org/grpc"
)

const (
	serverAddr = "localhost:8008"
)

type ChatServiceClient struct {
	pb.UnimplementedChatServiceServer
}

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
			log.Fatalf("connection problem: %v", err)
		}
	}(conn)

	client := pb.NewChatServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	go updateNewsfeed(ctx, client, &pb.User{})
	read(ctx, client)
}

func read(ctx context.Context, client pb.ChatServiceClient) {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, _ := reader.ReadString('\n')
		if strings.Contains(line, "/quit") {
			//code to leave chat here
			break
		}

		msg := &pb.Message{Text: line}

		client.Publish(ctx, msg)
	}
}

//protoc go types
//https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#google.protobuf.Any

func updateNewsfeed(ctx context.Context, client pb.ChatServiceClient, user *pb.User) {
	for {
		msg, err := client.Listen(ctx, user)
		if err != nil {
			log.Fatalf("listening problem: %v", err)
		}
		log.Println(msg.Text)
	}
}
