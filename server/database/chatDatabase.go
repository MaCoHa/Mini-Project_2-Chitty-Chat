package database

import (
	pb "example/Mini_Project_2_Chitty-Chat/chat"
	
	"sync"
)

type chatDatabase struct {
	connectedUsers  []*pb.User
	userToMesageMap map[string][]*pb.Message
	mu              sync.Mutex
}

func NewChatDatabase() *chatDatabase {
	return &chatDatabase{
		connectedUsers:  make([]*pb.User, 0),
		userToMesageMap: make(map[string][]*pb.Message)}
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

func (cd *chatDatabase) InsertMessage(msg *pb.Message) {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	for _, user := range cd.connectedUsers {
		/*if msg.User.Username == user.Username { //do not send message to the user who wrote it
			continue
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
	cd.userToMesageMap[user.Username] = cd.userToMesageMap[user.Username][1:]
	return msg
}
