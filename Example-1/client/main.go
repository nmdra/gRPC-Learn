package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/nmdra/gRPC-Learn/Example-1/pb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	// Create a new user
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	})
	if err != nil {
		log.Fatalf("Could not create user: %v", err)
	}
	fmt.Printf("User created: %v\n", res.User)

	// Fetch user details
	userRes, err := client.GetUser(ctx, &pb.GetUserRequest{Id: res.User.Id})
	if err != nil {
		log.Fatalf("Could not get user: %v", err)
	}
	fmt.Printf("User details: %v\n", userRes.User)
}
