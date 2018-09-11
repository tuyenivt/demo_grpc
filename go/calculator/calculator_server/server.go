package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"../calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("[calculator] Sum function was invoked with %v", req)
	first := req.GetFirstNumber()
	second := req.GetSecondNumber()
	result := first + second
	res := &calculatorpb.SumResponse{
		SumResult: result,
	}
	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("[calculator] PrimeNumberDecomposition function was invoked with %v\n", req)
	number := req.GetNumber()
	divisor := int64(2)
	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			number = number / divisor
		} else {
			divisor++
			fmt.Printf("[calculator] Divisor has increased to %v\n", divisor)
		}
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("[calculator] ComputeAverage function was invoked with a streaming request")
	sum := float64(0)
	count := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: sum / float64(count),
			})
		}
		if err != nil {
			log.Fatalf("[calculator] error while reading client stream: %v", err)
		}
		sum += req.GetNumber()
		count++
	}

}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Println("[calculator] FindMaximum function was invoked with a streaming request")
	max := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("[calculator] error while reading client stream: %v", err)
			return err
		}
		number := req.GetNumber()
		if number > max {
			max = number
			sendErr := stream.Send(&calculatorpb.FindMaximumResponse{
				Max: max,
			})
			if sendErr != nil {
				log.Fatalf("[calculator] error while sending data to client: %v", err)
				return err
			}
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Printf("[calculator] SquareRoot function was invoked with %v\n", req)
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("[calculator] Received a negative number: %v", number),
		)
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	fmt.Println("[calculator] Hello! This is server")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("[calculator] Failed to listen: %v\n", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("[calculator] Failed to serve: %v\n", err)
	}
}
