package mygrpc

import (
	"flag"
	"github.com/selfscrfc/PetBank/config"
	account "github.com/selfscrfc/PetBankProtos/proto/Accounts"
	Customers "github.com/selfscrfc/PetBankProtos/proto/Customers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewCustomerClient(cfg *config.Config) (*Customers.CustomerClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	serverAddr := flag.String("customerServerAddr", "localhost:50051", "The grpc address in the format of host:port")
	customerService, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		return nil, err
	}
	client := Customers.NewCustomerClient(customerService)

	return &client, nil
}

func NewAccountsClient(cfg *config.Config) (*account.AccountServiceClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	serverAddr := flag.String("accountsServerAddr", "localhost:50052", "The grpc address in the format of host:port")
	accountsService, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		return nil, err
	}
	client := account.NewAccountServiceClient(accountsService)

	return &client, nil
}
