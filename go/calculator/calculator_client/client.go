package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"../calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("[calculator] Hello! This is client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[calculator] Could not connect: %v\n", err)
	}
	defer cc.Close()
	c := calculatorpb.NewCalculatorServiceClient(cc)

	// doUnary(c)

	doServerStreaming(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("[calculator] Starting to do a Unary RPC...")
	req := &calculatorpb.SumRequest{
		FirstNumber:  10,
		SecondNumber: 15,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("[calculator] error while calling Sum RPC: %v\n", err)
	}
	log.Printf("[calculator] Response from Sum: %v\n", res.SumResult)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("[calculator] Starting to do a PrimeNumberDecomposition Server Streaming RPC...")
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 120,
	}
	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("[calculator] error while calling PrimeNumberDecomposition RPC: %v\n", err)
	}
	for {
		res, err := resStream.Recv()
		if err == io.EOF {
			// we're reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("[calculator] error while reading stream: %v\n", err)
		}
		fmt.Println(res.GetPrimeFactor())
	}

}
