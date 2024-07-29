package main

import (
	"context"
	"log"
	"time"

	pb "github.com/khiangte85/grpc-go/greet/proto"
)

func doGreet(c pb.GreetServiceClient) {
	t := time.Now()

	log.Println("do gree was invoked")

	res, err := c.Greet(context.Background(), &pb.GreetRequest{
		FirstName: "Khiangte",
	})

	if err != nil {
		log.Fatalf("Greet call failed: %v\n", err)
	}

	log.Println(res)
	log.Println("time taken is ", time.Since(t))
}
