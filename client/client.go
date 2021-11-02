package main

import (
	"bufio"
	"context"
	pb "example/Mini_Project_2_Chitty-Chat/chat"
	"fmt"
	"log"

	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
)

const (
	serverAddr = "localhost:8008"
)

var client pb.ChatServiceClient
var ctx context.Context
var user *pb.User

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

	client = pb.NewChatServiceClient(conn)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	
	user = connect()
	defer disconnect()

	go listen()
	read()
}

func read() {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, _ := reader.ReadString('\n')
		if strings.Contains(line, "/quit") {
			break
		}

		line = strings.Replace(line, "\n", "", 1)
		line = strings.Replace(line, "\r", "", 1)
		msg := &pb.Message{User: user, Text: line}

		client.Publish(ctx, msg)
	}
}

//protoc go types
//https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#google.protobuf.Any

func listen() {
	for {
		msg, err := client.Listen(ctx, user)
		if err != nil {
			log.Fatalf("listening problem: %v", err)
		}
		//log.Println(msg)
		log.Println(msg.User.Username + ": " + msg.Text)
	}
}

func connect() *pb.User {
	fmt.Println("Login with Username:")
	reader := bufio.NewReader(os.Stdin)
	var tryUser *pb.User

	for {
		username, _ := reader.ReadString('\n')
		username = strings.Replace(username, "\n", "", 1)
		username = strings.Replace(username, "\r", "", 1)
		tryUser = &pb.User{Username: username}

		resp, err := client.Connect(ctx, tryUser)
		if err != nil {
			log.Fatalf("connection problem: %v", err)
		}

		if strings.Contains(resp.Status, "Failed") {
			log.Println(resp)
			continue
		}
		break
	}

	return tryUser
}

func disconnect() {
	resp, err := client.Disconnect(ctx, user)
	if err != nil {
		log.Fatalf("disconnection problem: %v", err)
	}
	log.Println(resp)
}
