package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	port = ":8008"
)

type Server struct {
	course.UnimplementedCourseProtoServer
}

func (s *Server) GetCourse(ctx context.Context, in *course.GetCourseRequest) (*course.GetCourseReply, error) {
	log.Printf("Received GetCourse request")
	return &course.GetCourseReply{Message: "WIP: *List of current courses*"}, nil
}

func (s *Server) PostCourse(ctx context.Context, in *course.Course) (*course.PostCourseReply, error) {
	log.Printf("Received request to add Course: %s", in.Name)
	complete := fmt.Sprintf("Course: %s added to current available courses!", in.Name)
	return &course.PostCourseReply{Message: complete}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	course.RegisterCourseProtoServer(s, &Server{})
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
