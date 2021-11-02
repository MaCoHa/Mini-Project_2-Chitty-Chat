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
	success := chatData.AddUser(user)
	if !success {
		return &pb.Response{Status: "Join Failed: Username already taken! - Try another Username"}, nil
	}

	chatData.InsertMessage(&pb.Message{User: user, Text: user.Username + " has joined the chat!"})
	log.Println("Status: " + user.Username + " has joined the chat!")
	return &pb.Response{Status: "Join Successful"}, nil
}

func (s *ChatServiceServer) Disconnect(ctx context.Context, user *pb.User) (*pb.Response, error) {
	chatData.RemoveUser(user)
	chatData.InsertMessage(&pb.Message{User: user, Text: user.Username + " has left the chat!"})
	log.Println("Status: " + user.Username + " has left the chat!")
	return &pb.Response{Status: "Leaved Successful"}, nil
}

func (s *ChatServiceServer) Publish(ctx context.Context, msg *pb.Message) (*pb.Response, error) {
	resp, err := s.Broadcast(ctx, msg)
	if err != nil {
		log.Fatalf("Could not send message: %v", err)
	}
	return resp, err
}

func (s *ChatServiceServer) Broadcast(ctx context.Context, msg *pb.Message) (*pb.Response, error) {
	log.Println("Status: Broadcasting: " + msg.Text)
	chatData.InsertMessage(msg)
	return &pb.Response{Status: "Message Recieved"}, nil
}

func (s *ChatServiceServer) Listen(ctx context.Context, user *pb.User) (*pb.Message, error) {
	for {
		possibleMessage := chatData.PopMessage(user)
		if possibleMessage != nil {
			return possibleMessage, nil
		}
	}
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
	userToMesageMap map[string][]*pb.Message
	mu              sync.Mutex
}

var chatData *chatDatabase = &chatDatabase{
	connectedUsers:  make([]*pb.User, 0),
	userToMesageMap: make(map[string][]*pb.Message)}

func (cd *chatDatabase) InsertMessage(msg *pb.Message) {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	for _, user := range cd.connectedUsers {
		if msg.User.Username == user.Username { //do not send message to the user who wrote it
			continue
		}

		/*if cd.userToMesageMap[user.Username] != nil {
			log.Println("Status: Message overwritten: " + msg.Text + " - for user: " + user.Username)
		}*/
		cd.userToMesageMap[user.Username] = append(cd.userToMesageMap[user.Username], msg)
	}
}

func (cd *chatDatabase) PopMessage(user *pb.User) *pb.Message {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	if len(cd.userToMesageMap[user.Username]) < 1 {
		return nil
	}

	msg := cd.userToMesageMap[user.Username][0]
	log.Println("Status: Accesing message: " + msg.Text + " - for user: " + user.Username)

	cd.userToMesageMap[user.Username] = cd.userToMesageMap[user.Username][1:]
	return msg
}

func (cd *chatDatabase) AddUser(user *pb.User) bool {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	for i := range cd.connectedUsers {
		if cd.connectedUsers[i].Username == user.Username {
			return false
		}
	}

	cd.connectedUsers = append(cd.connectedUsers, user)
	cd.userToMesageMap[user.Username] = make([]*pb.Message, 0)

	return true
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

	delete(cd.userToMesageMap, user.Username)
}
