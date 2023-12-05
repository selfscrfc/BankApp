package main

import (
	"fmt"
	"github.com/pressly/goose"
	aServer "github.com/selfscrfc/PetBank/proto/Accounts/grpc"
	"github.com/selfscrfc/PetBank/utils"
	account "github.com/selfscrfc/PetBankProtos/proto/Accounts"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	port := os.Getenv("SERVER_PORT")

	log.Println("Account Service Starts on port:", port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	aServer.DB, err = utils.PostgreSQLConnection()
	if err != nil {
		log.Fatalf("DB connection fail: " + err.Error())
	}
	if err = goose.Up(aServer.DB, "migrations"); err != nil {
		log.Fatalf(err.Error())
	}

	grpcServer := grpc.NewServer()

	account.RegisterAccountServiceServer(grpcServer, aServer.AccountsServer{})

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf(err.Error())
	}
}
