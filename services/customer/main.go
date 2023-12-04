package main

import (
	"database/sql"
	"fmt"
	cServer "github.com/selfscrfc/PetBank/customer/server"
	"github.com/selfscrfc/PetBank/utils"
	"github.com/selfscrfc/PetBankProtos/proto/Customers"
	"google.golang.org/grpc"
	"log"
	"net"
)

var db *sql.DB

var port = 50051

func main() {
	log.Println("Customer Service Starte on port: ", port)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db, err := utils.PostgreSQLConnection()

	grpcServer := grpc.NewServer()

	customers.RegisterCustomerServer(grpcServer, cServer.CustomerServer{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf(err.Error())
	}
}
