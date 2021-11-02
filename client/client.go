package main

import (
	"bufio"
	"context"
	pb "example/Mini_Project_2_Chitty-Chat/chat"
	"fmt"
	"log"

	lamport "example/Mini_Project_2_Chitty-Chat/timestamp"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

const (
	serverAddr = "localhost:8008"
)

var client pb.ChatServiceClient
var ctx context.Context
var user *pb.User
var lamp *lamport.Clock

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

	lamp = lamport.NewClock()
	SetupCloseHandler()

	user = connect()
	defer disconnect()

	go listen()
	read()
}

func read() {
	reader := bufio.NewReader(os.Stdin)

	for {
		line, _ := reader.ReadString('\n')
		line = strings.Replace(line, "\n", "", 1)
		line = strings.Replace(line, "\r", "", 1)

		if len(line) > 128 {
			log.Println("Message to big! Max 128 characters!")
			continue
		}

		if strings.Contains(line, "/quit") {
			break
		}

		msg := &pb.Message{User: user, Text: line, Timestamp: lamp.Increment()}
		resp, err := client.Publish(ctx, msg)
		if err != nil {
			log.Fatalf("Broadcasting problem: %v", err)
		}

		lamp.Witness(resp.Timestamp)
	}
}

//protoc go types
//https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#google.protobuf.Any

func listen() {
	for {
		rec := &pb.Request{User: user, Timestamp: lamp.Increment()}

		msg, err := client.Listen(ctx, rec)
		if err != nil {
			log.Fatalf("listening problem: %v", err)
		}

		lamp.Witness(msg.Timestamp)
		msg.Timestamp = lamp.GetTimestamp()

		log.Printf("%d: %s: %s", msg.Timestamp, msg.User.Username, msg.Text)
	}
}

func connect() *pb.User {
	fmt.Println("Login with Username:")
	reader := bufio.NewReader(os.Stdin)

	for {
		username, _ := reader.ReadString('\n')
		username = strings.Replace(username, "\n", "", 1)
		username = strings.Replace(username, "\r", "", 1)

		tryUser := &pb.User{Username: username}
		rec := &pb.Request{User: tryUser, Timestamp: lamp.Increment()}

		resp, err := client.Connect(ctx, rec)
		if err != nil {
			log.Fatalf("connection problem: %v", err)
		}

		lamp.Witness(resp.Timestamp)
		resp.Timestamp = lamp.GetTimestamp()

		log.Println(resp)
		if strings.Contains(resp.Status, "Failed") {
			continue
		}
		return tryUser
	}
}

func disconnect() {
	rec := &pb.Request{User: user, Timestamp: lamp.Increment()}

	resp, err := client.Disconnect(ctx, rec)
	if err != nil {
		log.Fatalf("disconnection problem: %v", err)
	}

	lamp.Witness(resp.Timestamp)
	resp.Timestamp = lamp.GetTimestamp()

	log.Println(resp)
}

func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		disconnect()
		os.Exit(0)
	}()
}
