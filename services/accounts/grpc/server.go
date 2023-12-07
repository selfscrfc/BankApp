package grpcaccountsserver

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/selfscrfc/PetBank/api/models"
	"github.com/selfscrfc/PetBank/proto/Accounts/repository"
	Accounts "github.com/selfscrfc/PetBankProtos/proto/Accounts"
)

type AccountsServer struct {
	Accounts.UnimplementedAccountServiceServer
}

var AR = repository.AccountsRepo{DB: &sql.DB{}}

func (a AccountsServer) Create(ctx context.Context, request *Accounts.CreateRequest) (*Accounts.CreateResponse, error) {
	id, err := uuid.Parse(request.UserId)
	if err != nil {
		return nil, err
	}
	ac := &models.Account{
		Id:       uuid.New(),
		UserId:   id,
		IsCredit: request.IsCredit,
		Balance:  request.Balance,
		Currency: models.Currency(request.Currency),
	}

	ac, err = AR.CreateAccountQuery(ac)

	if err != nil {
		return nil, err
	}

	resp := &Accounts.CreateResponse{
		Id:       ac.Id.String(),
		UserId:   ac.UserId.String(),
		IsCredit: ac.IsCredit,
		Balance:  ac.Balance,
		Currency: Accounts.Currency(ac.Currency),
	}

	return resp, nil
}

func (a AccountsServer) GetDetail(ctx context.Context, request *Accounts.GetDetailsRequest) (*Accounts.GetDetailsResponse, error) {
	ac := &Accounts.Account{
		Id:     request.Id,
		UserId: request.UserId,
	}

	resp, err := AR.GetAccountDetailsQuery(ac)

	if err != nil {
		return nil, err
	}

	return &Accounts.GetDetailsResponse{
		Id:        resp.Id,
		UserId:    resp.UserId,
		IsCredit:  resp.IsCredit,
		Balance:   resp.Balance,
		Currency:  resp.Currency,
		IsBlocked: resp.IsBlocked,
	}, nil
}

func (a AccountsServer) Block(ctx context.Context, request *Accounts.BlockRequest) (*Accounts.BlockResponse, error) {
	ac := &Accounts.Account{Id: request.Id, UserId: request.UserId}

	resp, err := AR.BlockAccountQuery(ac)

	if err != nil {
		return nil, err
	}

	return &Accounts.BlockResponse{Success: resp}, nil
}

func (a AccountsServer) GetAll(ctx context.Context, request *Accounts.GetAllRequest) (*Accounts.GetAllResponse, error) {
	acs := make([]*Accounts.Account, 0)

	acs, err := AR.GetAllAccountsQuery(acs)

	if err != nil {
		return nil, err
	}

	return &Accounts.GetAllResponse{
		Accounts: acs,
	}, nil
}

func (a AccountsServer) RW(ctx context.Context, request *Accounts.RWRequest) (*Accounts.RWResponse, error) {
	resp, err := AR.RWQuery(request.AId, request.UId, request.Amount)

	if err != nil {
		return nil, err
	}

	return &Accounts.RWResponse{
		Amount:    int32(request.Amount),
		Id:        resp.Id,
		UserId:    resp.UserId,
		IsCredit:  resp.IsCredit,
		Balance:   resp.Balance,
		Currency:  resp.Currency,
		IsBlocked: resp.IsBlocked,
	}, nil
}
