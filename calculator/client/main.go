package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

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

	// doSum(c)
	// doPrimes(c)
	// doAverage(c)
	// doMax(c)
	doSqrt(c, -2358)
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
	nums := []float32{2, 5, 17, 3, 4, 9, 34, 24}

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

func doMax(c pb.CalculatorServiceClient) {
	reqs := []*pb.MaxRequest{
		{Number: 1},
		{Number: 6},
		{Number: 4},
		{Number: 17},
		{Number: 23},
		{Number: 14},
		{Number: 12},
		{Number: 40},
		{Number: 100},
		{Number: 10},
	}

	stream, err := c.Max(context.Background())

	if err != nil {
		log.Fatalf("error while connecting server: %v", err)
	}

	waitc := make(chan string)

	go func() {
		for _, req := range reqs {
			log.Printf("Send : %v", req.Number)

			err := stream.Send(req)
			time.Sleep(1 * time.Second)

			if err != nil {
				log.Fatalf("error while sending stream: %v\n", err)
			}

		}

		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				log.Fatalf("all data received: %v", err)
			}

			if err != nil {
				log.Fatalf("error while reading stream: %v", err)
			}

			log.Printf("Max: %v\n", res.Max)
		}
	}()

	<-waitc
}

func doSqrt(c pb.CalculatorServiceClient, n int32) {
	res, err := c.Sqrt(context.Background(), &pb.SqrtRequest{Number: n})

	if err != nil {
		e, ok := status.FromError(err)

		if ok {
			log.Printf("gRPC error code: %v\n", e.Proto().GetCode())
			log.Printf("gRPC error message: %v\n", e.Proto().GetMessage())
		} else {
			 log.Fatalf("non gRPC error: %v", err)
		}

		return 
	}

	log.Printf("Square root of %d is %f", n, res.Result)
}
