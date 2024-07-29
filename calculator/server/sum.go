package main

import (
	"context"
	"log"

	pb "github.com/khiangte85/grpc-go/calculator/proto"
)

func (s *Server) Sum(c context.Context, in *pb.SumRequest) (*pb.SumResponse, error) {
	log.Printf("Sum invoked: %v, %v", in.FirstNumber, in.SecondNumber)

	return &pb.SumResponse{
		Result: in.FirstNumber + in.SecondNumber,
	}, nil
}
