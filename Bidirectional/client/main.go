package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/nmdra/gRPC-Learn/Bidirectional/pb"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewChatServiceClient(conn)
	stream, err := client.ChatStream(context.Background())
	if err != nil {
		log.Fatalf("Error opening stream: %v", err)
	}

	// Receiving messages from server
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Fatalf("Error receiving message: %v", err)
			}
			fmt.Printf("\n%s: %s\n", msg.User, msg.Message)
		}
	}()

	// Sending messages
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter your name: ")
	scanner.Scan()
	username := scanner.Text()

	fmt.Println("Start chatting! Type your messages below:")

	for scanner.Scan() {
		message := scanner.Text()
		if message == "exit" {
			break
		}

		err := stream.Send(&pb.ChatMessage{
			User:    username,
			Message: message,
		})
		if err != nil {
			log.Fatalf("Error sending message: %v", err)
		}

		time.Sleep(time.Millisecond * 200)
	}

	stream.CloseSend()
}
