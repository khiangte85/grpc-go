package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	pb "github.com/khiangte85/grpc-go/greet/proto"
)

var addr string = "localhost:50052"

func main() {
	tls := true
	opts := []grpc.DialOption{}

	if tls {
		certFile := "ssl/ca.crt"
		creds, err := credentials.NewClientTLSFromFile(certFile, "")

		if err != nil {
			log.Fatalf("failed to load cert file: %v\n", certFile)
		}

		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	conn, err := grpc.NewClient(addr, opts...)

	if err != nil {
		log.Printf("Failed to connect: %v", err)
	}

	defer conn.Close()

	c := pb.NewGreetServiceClient(conn)

	// doGreet(c)
	// doGreetManyTimes(c)
	// doLongGreet(c)
	doGreetEveryone(c)
	// doGreetWithDeadline(c, 2*time.Second)
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

func doGreetEveryone(c pb.GreetServiceClient) {
	reqs := []string{
		"Awmtea",
		"Tei",
		"Sawmsawmi",
	}

	stream, err := c.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("Error while connecting server: %v\n", err)
	}

	waitc := make(chan string)

	go func() {

		for _, name := range reqs {

			log.Println("sending ", name)

			err := stream.Send(&pb.GreetRequest{FirstName: name})
			time.Sleep(1 * time.Second)

			if err != nil {
				log.Fatalf("error while sending stream to server: %v\n", err)
			}
		}

		stream.CloseSend()
	}()

	go func() {
		for {

			res, err := stream.Recv()

			if err == io.EOF {
				log.Printf("Received all message from server: %v", err)
				os.Exit(0)
			}

			if err != nil {
				log.Fatalf("error while reading stream from server: %v\n", err)
			}

			log.Println("Received: ", res.Result)

		}

		close(waitc)
	}()

	<-waitc
}

func doGreetWithDeadline(c pb.GreetServiceClient, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	res, err := c.GreetWithDeadline(ctx, &pb.GreetRequest{FirstName: "Lalmuanawma"})

	if err != nil {
		e, ok := status.FromError(err)

		if ok {
			log.Fatalf("gRPC error: %v", e)
		} else {
			log.Fatalf("other error: %v", err)
		}
	}

	log.Println(res.Result)
}
