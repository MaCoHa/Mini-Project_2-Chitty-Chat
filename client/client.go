package main

import (
	"context"
	pb "example/Mini_Project_2_Chitty-Chat/chat"
	"time"

	"google.golang.org/grpc"

	"os"
	"os/signal"
	"syscall"

	"bufio"
	"log"
	"strings"

	lamport "example/Mini_Project_2_Chitty-Chat/timestamp"
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
	//Setup the file for log outputs
	LOG_FILE := "./systemlogs/client.log"
	// open log file
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

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
	publish()
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

func connect() *pb.User {
	log.Println("Login with Username:")
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

func publish() {
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

func listen() {
	for {
		rec := &pb.Request{User: user, Timestamp: lamp.Increment()}

		msg, err := client.Listen(ctx, rec)
		if err != nil {
			log.Fatalf("listening problem: %v", err)
		}

		lamp.Witness(msg.Timestamp)
		msg.Timestamp = lamp.GetTimestamp()

		log.Printf("[%d - %s]: %s\n", msg.Timestamp, msg.User.Username, msg.Text)
	}
}
