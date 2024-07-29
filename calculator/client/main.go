package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/khiangte85/grpc-go/calculator/proto"
)

var addr string = "localhost:50051"

func main() {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect: %v\n", err)
	}

	defer conn.Close()

	c := pb.NewCalculatorServiceClient(conn)

	doSum(c)
}

func doSum(c pb.CalculatorServiceClient) {
	log.Println("Invoke Sum")

	res, err := c.Sum(context.Background(), &pb.SumRequest{
		FirstNumber:  20,
		SecondNumber: 24,
	})

	if err != nil {
		log.Fatalf("Client invoke error: %v\n", err)
	}

	log.Println(res)
}
