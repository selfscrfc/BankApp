package cServer

import (
	"context"
	"github.com/selfscrfc/PetBank/customer/postgres"
	"github.com/selfscrfc/PetBank/utils"
	customers "github.com/selfscrfc/PetBankProtos/proto/Customers"
)

type CustomerServer struct {
	customers.UnimplementedCustomerServer
}

func (c CustomerServer) Create(ctx context.Context, request *customers.CreateRequest) (*customers.CreateResponse, error) {
	db, err := utils.PostgreSQLConnection()
	if err != nil {
		return nil, err
	}

	cdb := postgres.NewCustomerRepo(db)

	cdb.Create(ctx, request.)
}
