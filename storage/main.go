package main

import (
	pb "censys_take_home/grpc"
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGRPCServer
	mutex   sync.RWMutex
	storage map[string]string
}

/*
gRPC setValue. Locks mutex for writing, and adds/updates key value pair
*/
func (s *server) SetValue(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	s.mutex.Lock()         // Do a write lock
	defer s.mutex.Unlock() // unlock on exit

	prev, exists := s.storage[req.Key] // storage[key] gives value, exists
	s.storage[req.Key] = req.Value

	if exists {
		result := fmt.Sprintf("Key %s existed with value %s. Now updated with value %s", req.Key, prev, req.Value)
		return &pb.SetResponse{Result: result}, nil
	}
	result := fmt.Sprintf("New key %s added with value %s", req.Key, req.Value)
	return &pb.SetResponse{Result: result}, nil
}

/*
gRPC getValue. Locks mutex for reading, and retrieves value for key.
Returns value "" if key does not exist. (Could probably have a smarter does not exist indication, but we can't store empty values anyways)
*/
func (s *server) GetValue(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	s.mutex.RLock()         // Do a read lock
	defer s.mutex.RUnlock() // unlock on exit

	fmt.Println("In GetValue for key:", req.Key)

	value, exists := s.storage[req.Key] // storage[key] gives value, exists

	if exists {
		return &pb.GetResponse{Key: req.Key, Value: value}, nil
	}
	return &pb.GetResponse{Key: req.Key, Value: ""}, nil
}

/*
gRPC deleteValue. Locks mutex for writing, and deletes key value pair if it exists.
*/
func (s *server) DeleteValue(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	s.mutex.Lock()         // Do a write lock
	defer s.mutex.Unlock() // unlock on exit

	_, exists := s.storage[req.Key] // storage[key] gives value, exists

	if exists {
		delete(s.storage, req.Key)
		result := fmt.Sprintf("Key %s existed and is now deleted", req.Key)
		return &pb.DeleteResponse{Result: result}, nil
	}
	result := fmt.Sprintf("Key %s does not exist", req.Key)
	return &pb.DeleteResponse{Result: result}, nil
}

func main() {
	fmt.Println("Hello from storage")

	//Following guide to set up GRPC

	lis, err := net.Listen("tcp", ":50000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	server := &server{ // Initialize just storage, mutex is zero value
		storage: make(map[string]string),
	}

	pb.RegisterGRPCServer(grpcServer, server)
	fmt.Println("gRPC server listening on port 50000")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
