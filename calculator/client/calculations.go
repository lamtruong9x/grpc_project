package main

import (
	"errors"
	"fmt"
	"io"
	"log"

	pb "github.com/lamtruong9x/grpc_project/calculator/proto"
	"golang.org/x/net/context"
)

func doSum(c pb.CalculatorServiceClient) error {
	req, err := c.Sum(context.Background(), &pb.SumRequest{
		A: 1,
		B: 5,
	})
	if err != nil {
		return err
	}
	fmt.Println("Sum:", req.GetSum())
	return nil
}

func doPrime(c pb.CalculatorServiceClient, n int32) error {
	cl, err := c.Prime(context.Background(), &pb.PrimeRequest{
		Num: n,
	})
	if err != nil {
		log.Println("Prime request error:", err)
		return err
	}
	fmt.Print("Primes: ")
	for {
		res, err := cl.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		fmt.Print(res.Prime)
	}
	fmt.Println()
	return nil
}