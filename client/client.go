package main

import (
	//"bufio"
	"context"
	pb "example/Mini_Project_2_Chitty-Chat/chat"
	"log"

	//"os"
	"strings"
	"time"

	tui "example/Mini_Project_2_Chitty-Chat/tui"

	"google.golang.org/grpc"
)

var Uimessage chan string = make(chan string)
var UIuserName chan string = make(chan string)

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

	tui.StartChatview(Uimessage)

	user = connect()
	defer disconnect()

	go listen()
	read()
}

func read() {
	//reader := bufio.NewReader(os.Stdin)
	for {
		var line string
		select {
		case l := <-Uimessage:
			line = l

		}
		//line, _ := reader.ReadString('\n')
		//line := tui.ReadFromChan()
		if strings.Contains(line, "/quit") {
			break
		}

		line = strings.Replace(line, "\n", "", 1)
		line = strings.Replace(line, "\r", "", 1)

		if len(line) > 128 {
			tui.Println("Message to big! Max 128 characters!")
			continue
		}

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
		tui.Println(msg.User.Username + ": " + msg.Text)
	}
}

func connect() *pb.User {
	tui.Println("Login with Username:")
	//reader := bufio.NewReader(os.Stdin)
	var tryUser *pb.User

	for {
		username := "mads"
		//username, _ := reader.ReadString('\n')
		//username := tui.ReadFromChan()
		username = strings.Replace(username, "\n", "", 1)
		username = strings.Replace(username, "\r", "", 1)
		tryUser = &pb.User{Username: username}

		resp, err := client.Connect(ctx, tryUser)
		if err != nil {
			log.Fatalf("connection problem: %v", err)
		}

		if strings.Contains(resp.Status, "Failed") {
			tui.Println(resp.Status)
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
	tui.Println(resp.Status)
}
