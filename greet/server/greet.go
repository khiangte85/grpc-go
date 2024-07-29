package main

import (
	"context"
	"log"

	pb "github.com/khiangte85/grpc-go/greet/proto"
)

func (s *Server) Greet(ctx context.Context, reg *pb.GreetRequest) (*pb.GreetResponse, error) {
	log.Printf("Greet function call with: %v\n", reg)

	return &pb.GreetResponse{
		Result: "Hello " + reg.FirstName,
	}, nil
}
