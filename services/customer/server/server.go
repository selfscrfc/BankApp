package cServer

import (
	"context"
	customers "github.com/selfscrfc/PetBankProtos/proto/Customers"
)

type CustomerServer struct {
	customers.UnimplementedCustomerServer
}

func (c CustomerServer) Create(ctx context.Context, request *customers.CreateRequest) (*customers.CreateResponse, error) {

}
