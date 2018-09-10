package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"../greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello World! This is Client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()
	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Printf("Created client: %f", c)

	// doUnary(c)

	// doServerStreaming(c)

	// doClientStreaming(c)

	doBiDirectionalStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Tuyen",
			LastName:  "Le Chau",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Tuyen",
			LastName:  "Le Chau",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we're reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet RPC: %v", err)
	}
	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tuyen 1",
				LastName:  "Le Chau",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tuyen 2",
				LastName:  "Le Chau",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tuyen 3",
				LastName:  "Le Chau",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tuyen 4",
				LastName:  "Le Chau",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tuyen 5",
				LastName:  "Le Chau",
			},
		},
	}
	// we iterate over our slice and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGreet: %v", err)
	}
	fmt.Printf("Response from LongGreet: %v\n", res)
}

func doBiDirectionalStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Bi Directional Streaming RPC...")

	// we create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while calling GreetingEveryone RPC: %v", err)
		return
	}
	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tuyen 1",
				LastName:  "Le Chau",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tuyen 2",
				LastName:  "Le Chau",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tuyen 3",
				LastName:  "Le Chau",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tuyen 4",
				LastName:  "Le Chau",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tuyen 5",
				LastName:  "Le Chau",
			},
		},
	}

	waitc := make(chan struct{})
	// we send a bunch of messages to the client (go routine)
	go func() {
		// function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// we receive a bunch of messages from the client (go routine)
	go func() {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while receiving message: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}
