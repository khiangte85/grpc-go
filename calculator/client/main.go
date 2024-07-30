package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
	doAverage(c)
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

func doAverage(c pb.CalculatorServiceClient) {
	nums := []float32{2,5,17,3,4,9,34,24}

	stream, err := c.Average(context.Background())

	if err != nil {
		log.Fatalf("Server call error: %v \n", err)
	}

	for _, num := range nums {
		log.Printf("Send %f to server\n", num)

		stream.Send(&pb.AverageRequest{Value: num})

		time.Sleep(500 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Receive error from server: %v", err)
	}

	fmt.Println("Average is: ", res.Result)
}
