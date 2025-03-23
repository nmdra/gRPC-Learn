package main

import (
	"io"
	"log"
	"net"
	"sync"

	pb "github.com/nmdra/gRPC-Learn/Bidirectional/pb"

	"google.golang.org/grpc"
)

type ChatServer struct {
	pb.UnimplementedChatServiceServer
	mu      sync.Mutex
	clients map[pb.ChatService_ChatStreamServer]bool
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		clients: make(map[pb.ChatService_ChatStreamServer]bool),
	}
}

func (s *ChatServer) ChatStream(stream pb.ChatService_ChatStreamServer) error {
	s.mu.Lock()
	s.clients[stream] = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, stream)
		s.mu.Unlock()
	}()

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			return err
		}

		log.Printf("[%s]: %s", msg.User, msg.Message)

		// Broadcast message to all connected clients
		s.mu.Lock()
		for client := range s.clients {
			if err := client.Send(msg); err != nil {
				log.Printf("Error sending message: %v", err)
			}
		}
		s.mu.Unlock()
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	chatServer := NewChatServer()
	pb.RegisterChatServiceServer(grpcServer, chatServer)

	log.Println("Chat server running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
