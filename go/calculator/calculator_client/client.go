package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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

	// doServerStreaming(c)

	doClientStreaming(c)
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

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("[calculator] Starting to do a ComputeAverage Server Streaming RPC...")
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("[calculator] error while calling ComputeAverage RPC: %v\n", err)
	}
	numbers := []float64{1.0, 2.0, 3.0, 4.0}
	// we iterate over our slice and send each message individually
	for _, number := range numbers {
		fmt.Printf("[calculator] Sending number: %v\n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
		time.Sleep(100 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("[calculator] error while receiving response from ComputeAverage: %v", err)
	}
	fmt.Printf("[calculator] Response from ComputeAverage: %v\n", res.GetAverage())
}
