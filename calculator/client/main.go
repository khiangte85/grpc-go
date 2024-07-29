package main

import (
	"context"
	"fmt"
	"io"
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
	doPrimes(c)
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

func doPrimes(c pb.CalculatorServiceClient) {
	stream, err := c.Primes(context.Background(), &pb.PrimeRequest{
		Number: 123456789021548,
	})

	if err != nil {
		log.Fatalf("error calling primes: %v", err)
	}

	for {
		prime, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}

		fmt.Println(prime.Prime)
	}
}
