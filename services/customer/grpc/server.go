package grpccustomerserver

import (
	"context"
	"database/sql"
	"github.com/selfscrfc/PetBank/customer/repository"
	customers "github.com/selfscrfc/PetBankProtos/proto/Customers"
)

var CR = &repository.CustomerRepo{DB: &sql.DB{}}

type CustomerServer struct {
	customers.UnimplementedCustomerServer
}

func (c CustomerServer) Create(ctx context.Context, request *customers.CreateRequest) (*customers.CreateResponse, error) {
	cs := &customers.CustomerEntity{
		FullName:  request.FullName,
		Login:     request.Login,
		IsBlocked: false,
	}

	cs, err := CR.CreateCustomerQuery(cs, request.Password)

	if err != nil {
		return nil, err
	}

	return &customers.CreateResponse{
		Id:       cs.Id,
		FullName: cs.FullName,
		Time:     cs.Time,
		Login:    cs.Login,
		Password: "",
	}, nil
}

func (c CustomerServer) GetDetails(ctx context.Context, request *customers.GetDetailsRequest) (*customers.GetDetailsResponse, error) {
	cs := &customers.CustomerEntity{
		Id: request.Id,
	}

	cs, err := CR.GetCustomerDetailsQuery(cs)
	if err != nil {
		return nil, err
	}

	return &customers.GetDetailsResponse{
		Id:        cs.Id,
		FullName:  cs.FullName,
		Time:      cs.Time,
		Login:     cs.Login,
		IsBlocked: cs.IsBlocked,
	}, nil
}

func (c CustomerServer) Block(ctx context.Context, request *customers.BlockRequest) (*customers.BlockResponse, error) {
	_, err := CR.BlockAccountQuery(request.BlockId)

	if err != nil {
		return nil, err
	}

	return &customers.BlockResponse{
		Success: true,
	}, nil
}

func (c CustomerServer) GetAll(ctx context.Context, request *customers.GetAllRequest) (*customers.GetAllResponse, error) {
	csa := make([]*customers.CustomerEntity, 0)

	csa, err := CR.GetAllCustomersQuery(csa)

	if err != nil {
		return nil, err
	}

	return &customers.GetAllResponse{
		Customers: csa,
	}, nil
}

func (c CustomerServer) SignIn(ctx context.Context, request *customers.SignInRequest) (*customers.SignInResponse, error) {
	resp, err := CR.SignInQuery(request.Login, request.Password)

	if err != nil {
		return nil, err
	}

	return &customers.SignInResponse{
		Id:        resp.Id,
		FullName:  resp.FullName,
		Time:      resp.Time,
		Login:     resp.Login,
		IsBlocked: resp.IsBlocked,
	}, nil
}
