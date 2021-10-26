package main

import (
	"bufio"
	"context"
	pb "example/Mini_Project_2_Chitty-Chat/chat"
	"log"
	"os"
	"time"

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

		}
	}(conn)

	client := pb.NewChatServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	/*resp, err := client.Publish(ctx, &pb.Msg{Message: "Test bish"})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if resp != nil {
		log.Printf("success, %v", resp)
	} else {
		log.Printf("problem, %v", err)
	}*/

	read(ctx, client)
}

func read(ctx context.Context, client pb.ChatServiceClient) {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, _ := reader.ReadString('\n')
		if line == "quit" {
			//code to leave chat here
			break
		}

		msg := &pb.Msg{Message: line}

		client.Publish(ctx, msg)
	}
}

func updateNewsfeed(ctx context.Context, client ChatServiceClient) {
	//call client and wait for response from server
	//response should contain the broadcasted messages
}
