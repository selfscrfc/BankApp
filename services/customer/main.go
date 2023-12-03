package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/selfscrfc/PetBank/proto/Customers"
	"google.golang.org/grpc"
	"log"
	"net"
)

type CustomerServer struct {
	Customers.UnimplementedCustomerServer
}

func (c CustomerServer) Create(ctx context.Context, request *Customers.CreateRequest) (*Customers.CreateResponse, error) {
	resp, err := toPSQL()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

var port = 50051

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	Customers.RegisterCustomerServer(grpcServer, CustomerServer{})
	grpcServer.Serve(lis)
}

func toPSQL() (*Customers.CreateResponse, error) {
	return &Customers.CreateResponse{
		Id:      1,
		Success: true,
		Error:   "",
	}, nil
}
