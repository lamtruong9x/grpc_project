package main

import (
	"errors"
	"log"
	"time"

	pb "github.com/lamtruong9x/grpc_project/calculator/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const addr string = "localhost:50123"

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			log.Println("Cannot dial to server - server may currently offline.")
			return
		default:
			log.Fatal("Error: ", err)
		}
	}
	c := pb.NewCalculatorServiceClient(conn)

	if err := doSum(c); err != nil {
		log.Println("Sum request error: ", err)
	}

	if err := doPrime(c, 120); err != nil {
		log.Println("Sum request error: ", err)
	}

}