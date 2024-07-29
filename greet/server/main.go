package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/khiangte85/grpc-go/greet/proto"
	"google.golang.org/grpc"
)

var addr string = "0.0.0.0:50052"

type Server struct {
	pb.GreetServiceServer
}

func main() {
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("Failed to listen on %v\n", err)
	}

	log.Printf("Listen on %v\n", addr)

	s := grpc.NewServer()
	pb.RegisterGreetServiceServer(s, &Server{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}

func (s *Server) Greet(ctx context.Context, reg *pb.GreetRequest) (*pb.GreetResponse, error) {
	log.Printf("Greet invoked with: %v\n", reg)

	return &pb.GreetResponse{
		Result: "Hello " + reg.FirstName,
	}, nil
}

func (s *Server) GreetManyTimes(in *pb.GreetRequest, stream pb.GreetService_GreetManyTimesServer) error {
	log.Printf("GreetManyTimes invoked with: %v\n", in.FirstName)

	for i := 0; i < 10; i++ {
		res := &pb.GreetResponse{
			Result: fmt.Sprintf("Hello %s, number %d", in.FirstName, i),
		}

		if err := stream.Send(res); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) LongGreet(stream pb.GreetService_LongGreetServer) error {
	res := ""

	for {
		req, err := stream.Recv()


		if err == io.EOF {
			return stream.SendAndClose(&pb.GreetResponse{Result: res})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		log.Printf("received %v\n", req.FirstName)
		res += fmt.Sprintf("Hello %s\n", req.FirstName)
	}
}
