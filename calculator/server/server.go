package main

import (
	"calculator/calculator/calculatorpb"
	"fmt"

	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type server struct{
	calculatorpb.UnimplementedCalculatorServiceServer
}


func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("received PrimeNumberDecomposition RPC: %v\n", req)
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
			log.Printf("divisor increased to %v", divisor)
		}
	}

	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("received ComputeAverage RPC")

	sum := int32(0)
	count := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := float64(sum) / float64(count)
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average,
			})
		}
		if err != nil {
			log.Fatalf("error client stream: %v", err)
		}

		sum += req.GetNumber()
		count++
	}
}


func main() {
	fmt.Println("running gRPC on port 4000...")
	lis, err := net.Listen("tcp", "0.0.0.0:4000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("error serve: %v", err)
	}
}
