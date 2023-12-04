package main

import (
	"fmt"
	cServer "github.com/selfscrfc/PetBank/customer/server"
	"github.com/selfscrfc/PetBankProtos/proto/Customers"
	"google.golang.org/grpc"
	"log"
	"net"
)

var port = 50051

func main() {
	log.Println("Customer Service Starte on port: ", port)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	customers.RegisterCustomerServer(grpcServer, cServer.CustomerServer{})
	grpcServer.Serve(lis)
}
