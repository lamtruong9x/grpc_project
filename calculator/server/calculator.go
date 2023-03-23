package main

import (
	"context"

	pb "github.com/lamtruong9x/grpc_project/calculator/proto"
)

func (c *Calculator) Sum(ctx context.Context, in *pb.SumRequest) (*pb.SumResponse, error) {
	sum := pb.SumResponse{
		Sum: in.GetA() + in.GetB(),
	}

	return &sum, nil
}

func (c *Calculator) Prime(in *pb.PrimeRequest, stream pb.CalculatorService_PrimeServer) error {
	N := int(in.GetNum())
	outputStream := factorToPrime(N)
	for {
		p, ok := <-outputStream
		if !ok {
			break
		}
		stream.Send(&pb.PrimeResponse{
			Prime: int32(p),
		})
	}
	return nil
}

func factorToPrime(N int) <-chan int {
	result := make(chan int)
	go func() {
		defer close(result)
		k := 2
		for N > 1 {
			if N%k == 0 {
				result <- k
				N = N / k
			} else {
				// divide N by k so that we have the rest of the number left.
				k = k + 1
			}
		}
	}()
	return result
}
