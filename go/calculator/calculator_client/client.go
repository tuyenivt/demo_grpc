package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

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

	// doClientStreaming(c)

	// doBiDirectionalStreaming(c)

	doErrorUnary(c)
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

func doBiDirectionalStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("[calculator] Starting to do a Bi Directional Streaming RPC...")
	// we create a stream by invoking the client
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("[calculator] error while calling FindMaximum RPC: %v\n", err)
	}
	waitc := make(chan struct{})
	// we send a bunch of numbers to the client (go routine)
	go func() {
		// function to send a bunch of numbers
		numbers := []int32{1, 5, 3, 6, 2, 20}
		for _, number := range numbers {
			fmt.Printf("[calculator] Sending number: %v\n", number)
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	// we receive a bunch of numbers from the client (go routine)
	go func() {
		// function to receive a bunch of numbers
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("[calculator] error while receiving max of numbers: %v", err)
				break
			}
			fmt.Printf("[calculator] Received new maximum: %v\n", res.GetMax())
		}
		close(waitc)
	}()
	// block until everything is done
	<-waitc
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("[calculator] Starting to do a SquareRoot Unary RPC...")
	// correct call
	doErrorCall(c, 10)
	// error call
	doErrorCall(c, -1)
}

func doErrorCall(c calculatorpb.CalculatorServiceClient, number int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: number})
	if err != nil {
		resErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC (user error)
			fmt.Println(resErr.Message())
			fmt.Printf("[calculator] Error code from server: %v\n", resErr.Code())
			if resErr.Code() == codes.InvalidArgument {
				fmt.Println("[calculator] We probably sent a negative number!")
			}
		} else {
			log.Fatalf("[calculator] Error while calling SquareRoot RPC: %v\n", err)
		}
	} else {
		fmt.Printf("[calculator] Result of square root of %v: %v\n", number, res.GetNumberRoot())
	}
}
