package main

import (
	"context"
	pb "example/Mini_Project_2_Chitty-Chat/chat"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

const (
	port = ":8008"
)

type ChatServiceServer struct {
	pb.UnimplementedChatServiceServer
}

func (s *ChatServiceServer) Connect(ctx context.Context, user *pb.User) (*pb.Response, error) {
	chatData.AddUser(user)

	//messageList = append(messageList, user.Username+" is connected")

	return &pb.Response{Status: "Join Successful"}, nil
}

func (s *ChatServiceServer) Publish(ctx context.Context, msg *pb.Message) (*pb.Response, error) {
	resp, err := s.Broadcast(ctx, msg)
	if err != nil {
		log.Fatalf("Could not send message: %v", err)
	}
	return resp, nil
}

func (s *ChatServiceServer) Broadcast(ctx context.Context, msg *pb.Message) (*pb.Response, error) {
	log.Printf("Broadcasting:" + msg.Text)
	chatData.InsertMessage(msg)
	return &pb.Response{Status: "Message Recieved"}, nil
}

func (s *ChatServiceServer) Listen(ctx context.Context, user *pb.User) (*pb.Message, error) {
	var newMsg *pb.Message

	for {
		possibleMessage := chatData.PopMessage()
		if possibleMessage != nil {
			newMsg = possibleMessage
			break
		}
	}

	/*for {
		if userToMesageMap[user] != nil {
			newMsg = userToMesageMap[user]
			userToMesageMap[user] = nil
			break
		}
	}*/

	/*for {
		if len(messageList) > 0 {
			newMsg = &pb.Message{Text: messageList[0]}
			messageList = messageList[1:]
			break
		}
	}*/

	return newMsg, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterChatServiceServer(s, &ChatServiceServer{})
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

//here
//begins
//something
//special

type chatDatabase struct {
	connectedUsers  []*pb.User
	messageList     []string
	userToMesageMap map[*pb.User]*pb.Message
	mu              sync.Mutex
}

var chatData *chatDatabase = &chatDatabase{
	connectedUsers:  make([]*pb.User, 0),
	messageList:     make([]string, 0),
	userToMesageMap: make(map[*pb.User]*pb.Message)}

func (cd *chatDatabase) GetMessage(user *pb.User) *pb.Message {
	cd.mu.Lock()
	defer cd.mu.Unlock()
	return cd.userToMesageMap[user]
}

func (cd *chatDatabase) InsertMessage(msg *pb.Message) {
	cd.mu.Lock()
	defer cd.mu.Unlock()
	if cd.userToMesageMap[msg.User] != nil {
		log.Println("Message overwritten: " + msg.Text)
	}
	cd.userToMesageMap[msg.User] = msg

	cd.messageList = append(cd.messageList, msg.Text)
}

func (cd *chatDatabase) PopMessage() *pb.Message {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	if len(cd.messageList) == 0 {
		return nil
	}

	msg := &pb.Message{Text: cd.messageList[0]}
	cd.messageList = cd.messageList[1:]
	return msg
}

func (cd *chatDatabase) AddUser(user *pb.User) {
	cd.mu.Lock()
	defer cd.mu.Unlock()
	cd.connectedUsers = append(cd.connectedUsers, user)
}

func (cd *chatDatabase) RemoveUser(user *pb.User) {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	var newConnectedUsers []*pb.User = make([]*pb.User, 0)
	for i := range cd.connectedUsers {
		if cd.connectedUsers[i].Username != user.Username {
			newConnectedUsers = append(newConnectedUsers, cd.connectedUsers[i])
		}
	}
	cd.connectedUsers = newConnectedUsers
}
