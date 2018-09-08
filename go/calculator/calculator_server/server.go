package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"../calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum function was invoked with %v", req)
	first := req.GetAddition().GetFirstAddition()
	second := req.GetAddition().GetSecondAddition()
	result := first + second
	res := &calculatorpb.SumResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	fmt.Printf("[calculator] Hello! This is server")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("[calculator] Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("[calculator] Failed to serve: %v", err)
	}
}
