package main

import (
	"context"
	"errors"
	"io"
	"log"
	"math"

	pb "github.com/lamtruong9x/grpc_project/calculator/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Unary
func (c *Calculator) Sum(ctx context.Context, in *pb.SumRequest) (*pb.SumResponse, error) {
	sum := pb.SumResponse{
		Sum: in.GetA() + in.GetB(),
	}

	return &sum, nil
}

// Server Streaming
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

// Client Streaming
func (c *Calculator) Avg(stream pb.CalculatorService_AvgServer) error {
	var intStream = make(chan int, 5)
	result := avg(intStream)
	for {
		req, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		intStream <- int(req.GetNum())
	}
	close(intStream)
	stream.SendAndClose(&pb.AvgResponse{
		Avg: <-result,
	})
	return nil
}

func avg(intStream <-chan int) <-chan float64 {
	var result = make(chan float64)
	go func() {
		var sum int
		var i int
		defer close(result)
		for {
			v, ok := <-intStream
			if !ok {
				break
			}
			i++
			sum += v
		}
		if i == 0 {
			result <- 0
			return
		}
		result <- (float64(sum) / float64(i))
	}()
	return result
}

func (c *Calculator) Max(stream pb.CalculatorService_MaxServer) error {
	var max = math.Inf(-1)
	for {
		req, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		temp := float64(req.GetNum())
		if temp > max {
			max = temp
		}

		err = stream.Send(&pb.MaxResponse{Max: int32(max)})
		if err != nil {
			log.Println("Cannot send response")
			return err
		}
	}
}

func (c *Calculator) Sqrt(ctx context.Context, in *pb.SqrtRequest) (*pb.SqrtResponse, error) {
	num := in.GetNum()
	if num < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "cannot sqare root negative number: %d", num)
	}
	return &pb.SqrtResponse{
		Sqrt: math.Sqrt(float64(num)),
	}, nil
}
