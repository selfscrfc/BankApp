package main

import (
	"fmt"
	"github.com/pressly/goose/v3"
	cServer "github.com/selfscrfc/PetBank/customer/grpc"
	"github.com/selfscrfc/PetBank/utils"
	"github.com/selfscrfc/PetBankProtos/proto/Customers"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {

	port := os.Getenv("SERVER_PORT")

	log.Println("Customer Service Starts on port:", port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	cServer.CR.DB, err = utils.PostgreSQLConnection()
	if err != nil {
		log.Fatalf("DB connection fail: " + err.Error())
	}
	if err = goose.Up(cServer.CR.DB, "migrations"); err != nil {
		log.Fatalf(err.Error())
	}

	grpcServer := grpc.NewServer()

	customers.RegisterCustomerServer(grpcServer, cServer.CustomerServer{})

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf(err.Error())
	}
}
