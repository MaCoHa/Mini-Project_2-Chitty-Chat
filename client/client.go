package main

import (
	"context"
	pb "gRPC/chat"
	"log"

	"google.golang.org/grpc"
)

const (
	serverAddr = "localhost:8008"
)

func main() {
	//var opts []grpc.DialOption
	// Set up a connection to the server.
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewCourseProtoClient(conn)

	SendPostCourseRequest(client)
	SendGetCourseRequest(client)
}

func SendPostCourseRequest(client pb.CourseProtoClient) {
	course := pb.Course{
		CourseId:          1,
		Name:              "DØSYS",
		Workload:          15,
		SatisfactoryScore: 10,
		Teachers:          []string{"Tom"},
		Students:          []string{"Bente"},
	}
	log.Printf("Course: %s added to courses", course.Name)
}

func SendGetCourseRequest(client pb.CourseProtoClient) {
	msg := pb.GetCourseRequest{}

	resp, err := client.GetCourse(context.Background(), &msg)
	if err != nil {
		log.Fatalf("Error when using GetCourse: %s", err)
	}

	log.Printf("Server responeded with: %s \n", resp)
}
