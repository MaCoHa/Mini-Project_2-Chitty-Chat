package main

import (
	"context"
	pb "example/Mini_Project_2_Chitty-Chat/chat"
	"log"
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp := client.Publish(ctx, &pb.Msg{Message: "Test bish"})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if resp != nil {
		log.Fatalf("success, %v", resp)
	} else {
		log.Fatalf("problem, %v", err)
	}

}

func (client *ChatServiceClient) Publish(ctx context.Context, msg *pb.Msg) *pb.Response {
	resp, err := client.Broadcast(ctx, msg)
	if err != nil {
		log.Fatalf("Could not send message: %v", err)
	}
	return resp, nil
}
