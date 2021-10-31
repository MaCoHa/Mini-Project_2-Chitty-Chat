package main

import (
	"bufio"
	"context"
	pb "example/Mini_Project_2_Chitty-Chat/chat"
	"log"
	"math/rand"
	"os"
	"strconv"
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
	ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var username string = strconv.Itoa(rand.Intn(1000))
	user = &pb.User{Username: username}

	connect()
	go updateNewsfeed()
	read()
}

func read() {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, _ := reader.ReadString('\n')
		if strings.Split(line, " ")[0] == "/quit" {
			//code to leave chat here
			break
		}

		msg := &pb.Message{User: user, Text: line}

		client.Publish(ctx, msg)
	}
}

//protoc go types
//https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#google.protobuf.Any

func updateNewsfeed() {
	for {
		msg, err := client.Listen(ctx, user)
		if err != nil {
			log.Fatalf("listening problem: %v", err)
		}
		log.Println(msg.Text)
	}
}

func connect() {
	resp, err := client.Connect(ctx, user)
	if err != nil {
		log.Fatalf("listening problem: %v", err)
	}
	log.Println(resp)
}
