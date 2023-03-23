package main

import (
	"log"
	"net"

	pb "github.com/lamtruong9x/grpc_project/calculator/proto"
	"google.golang.org/grpc"
)

const addr string = "localhost:50123"

type Calculator struct {
	pb.CalculatorServiceServer
}

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on port %v", addr)
	s := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(s, &Calculator{})
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
