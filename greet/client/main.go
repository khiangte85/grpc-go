package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/khiangte85/grpc-go/greet/proto"
)

var addr string = "localhost:50052"

func main() {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("Failed to connect: %v", err)
	}

	defer conn.Close()

	c := pb.NewGreetServiceClient(conn)

	// doGreet(c)
	// doGreetManyTimes(c)
	doLongGreet(c)
}

func doGreet(c pb.GreetServiceClient) {
	res, err := c.Greet(context.Background(), &pb.GreetRequest{
		FirstName: "Khiangte",
	})

	if err != nil {
		log.Fatalf("Greet call failed: %v\n", err)
	}

	log.Println(res)
}

func doGreetManyTimes(c pb.GreetServiceClient) {
	stream, err := c.GreetManyTimes(context.Background(), &pb.GreetRequest{
		FirstName: "Lalmuanawma",
	})

	if err != nil {
		log.Fatalf("rpc call error: %v\n", err)
	}

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("stream recv error: %v\n", err)
		}

		fmt.Println(msg)
	}
}

func doLongGreet(c pb.GreetServiceClient) {
	reqs := []*pb.GreetRequest{
		{FirstName: "Awmtea"},
		{FirstName: "Tei"},
		{FirstName: "Sawmsawmi"},
	}

	stream, err := c.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("error calling long greet: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 2)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Stream error: %v", err)
	}

	fmt.Println(res)
}
