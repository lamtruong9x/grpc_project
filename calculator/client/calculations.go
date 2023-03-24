package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"sync"

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

func doAvg(c pb.CalculatorServiceClient, nums ...int32) error {
	cl, err := c.Avg(context.Background())
	if err != nil {
		log.Println("Avg request error:", err)
		return err
	}
	for _, num := range nums {
		err = cl.Send(&pb.AvgRequest{
			Num: num,
		})
		if err != nil {
			log.Println("Avg request error:", err)
			return err
		}
	}

	res, err := cl.CloseAndRecv()
	if err != nil {
		return err
	}
	fmt.Println("The avg is:", res.GetAvg())
	return nil
}

func doMax(c pb.CalculatorServiceClient, nums ...int32) error {
	var wg = &sync.WaitGroup{}

	wg.Add(2)
	cl, err := c.Max(context.Background())
	if err != nil {
		return err
	}

	go func() {
		defer wg.Done()
		for _, num := range nums {
			err := cl.Send(&pb.MaxRequest{
				Num: num,
			})
			if err != nil {
				log.Println("cannot send request due to", err)
				return
			}
		}
		if err := cl.CloseSend(); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		defer wg.Done()
		for {
			res, err := cl.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				log.Fatal(err)
			}
			if res != nil {
				fmt.Println("Current max is:", res.GetMax())

			}
		}
	}()

	wg.Wait()
	return nil
}
