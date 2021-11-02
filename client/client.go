package client

import (
	"context"
	pb "example/Mini_Project_2_Chitty-Chat/chat"
	"log"

	"strings"
	"time"

	"github.com/marcusolsson/tui-go"
	"google.golang.org/grpc"
)

var Uimessage chan string = make(chan string)

const (
	serverAddr = "localhost:8008"
)

var client pb.ChatServiceClient
var ctx context.Context
var user *pb.User

type ChatServiceClient struct {
	pb.UnimplementedChatServiceServer
}

func Startclint() {
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
	StartChatview()

}

func read(message string) {

	message = strings.Replace(message, "\n", "", 1)
	message = strings.Replace(message, "\r", "", 1)

	if len(message) > 128 {
		ReciveMessage("Message to big! Max 128 characters!")

	}

	msg := &pb.Message{User: user, Text: message}
	client.Publish(ctx, msg)

}

//protoc go types
//https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#google.protobuf.Any

func listen() {
	for {

		// msg, err := client.Listen(ctx, user)
		// if err != nil {
		// 	log.Fatalf("listening problem: %v", err)
		// }
		// var message string
		// message = " " + msg.User.Username + ": " + msg.Text
		// chatview.messageHistory.Append(tui.NewLabel("message"))

	}
}

func connect(username string) *pb.User {

	var tryUser *pb.User

	for {
		tryUser = &pb.User{Username: username}

		resp, err := client.Connect(ctx, tryUser)
		if err != nil {
			log.Fatalf("connection problem: %v", err)
		}

		if strings.Contains(resp.Status, "Failed") {
			ReciveMessage(resp.Status)
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
	ReciveMessage(resp.Status)
}

var chatview *ChatView
var Ui tui.UI
var c chan string

// needs to get a client varible that can be used
func StartChatview() {

	chatview := NewChatView()
	chatlogin := NewChatLogin()

	ui, err := tui.New(chatlogin.view)
	if err != nil {
		log.Fatal(err)
	}
	exit := func() {
		// send user logout message to server
		defer disconnect()
		ui.Quit()
	}
	ui.SetKeybinding("Esc", exit)
	ui.SetKeybinding("Ctrl+c", exit)
	Ui = ui
	chatlogin.Login(func(username string) {
		user = connect(username)
		//username is the new user joining the chat. call
		//the server with the name

		ui.SetWidget(chatview.view)

		go listen()
		//defer cancel()
		chatview.messageHistory.Append(tui.NewLabel("willcome to Chitty chat"))
	})

	chatview.SendMessage(func(message string) {
		read(message)
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
func ReciveMessage(msg string) {
	Ui.Update(func() { chatview.ReciveMessage(msg) })
}
