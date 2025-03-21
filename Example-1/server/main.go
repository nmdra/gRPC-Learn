package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/nmdra/gRPC-Learn/Example-1/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedUserServiceServer
	users map[string]*pb.User
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	id := uuid.New().String()
	user := &pb.User{
		Id:    id,
		Name:  req.Name,
		Email: req.Email,
	}
	s.users[id] = user

	return &pb.CreateUserResponse{User: user}, nil
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, exists := s.users[req.Id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return &pb.GetUserResponse{User: user}, nil
}

func (s *server) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	var userList []*pb.User
	for _, user := range s.users {
		userList = append(userList, user)
	}
	return &pb.ListUsersResponse{Users: userList}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &server{
		users: make(map[string]*pb.User),
	})

	// **Enable reflection**
	reflection.Register(grpcServer)

	log.Println("gRPC server started on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
