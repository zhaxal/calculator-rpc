package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"
)
import "calculator/calculator/calculatorpb"

func main() {
	fmt.Println("client")
	cc, err := grpc.Dial("localhost:4000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	doPrime(c)
	doAverage(c)

}


func doPrime(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("starting PrimeDecomposition")
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 120,
	}
	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error PrimeDecomposition RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		fmt.Printf("prime factor: %v\n", res.GetPrimeFactor())
	}
}

func doAverage(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("starting ComputeAverage")
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error stream: %v", err)
	}

	numbers := []int32{1, 2, 3, 4, 5}

	for _, number := range numbers {
		fmt.Printf("sent: %v\n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error response: %v", err)
	}

	fmt.Printf("average: %v\n", res.GetAverage())
}

