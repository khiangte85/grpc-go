package main

import (
	"io"
	"log"
	"net"

	pb "github.com/khiangte85/grpc-go/calculator/proto"
	"google.golang.org/grpc"
)

var addr string = "0.0.0.0:50051"

type Server struct {
	pb.CalculatorServiceServer
}

func main() {
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("Failed to listen on %v\n", err)
	}

	log.Printf("Listen on %v\n", addr)

	s := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(s, &Server{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}

}

func (s *Server) Primes(in *pb.PrimeRequest, stream pb.CalculatorService_PrimesServer) error {
	k := int64(2)
	n := in.Number

	for n > 1 {
		if n%k == 0 {
			stream.Send(&pb.PrimeResponse{Prime: k})
			n = n / k
		} else {
			k++
		}
	}

	return nil
}

func (s *Server) Average(stream pb.CalculatorService_AverageServer) error {
	total := float32(0.0)
	count := 0

	for {
		reg, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&pb.AverageResponse{
				Result: total / float32(count),
			})
		}

		if err != nil {
			log.Fatalf("stream receive error: %v \n", err)
		}

		log.Printf("Received %f from client\n", reg.Value)

		total += reg.Value
		count++
	}
}
