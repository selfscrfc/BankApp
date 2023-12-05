package grpcaccountsserver

import (
	"context"
	"database/sql"
	qb "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	Accounts "github.com/selfscrfc/PetBankProtos/proto/Accounts"
)

var DB *sql.DB

type AccountsServer struct {
	Accounts.UnimplementedAccountServiceServer
}

const (
	TableAccounts   = "accounts"
	ColumnsAccounts = "id, userid, iscredit, balance, currency, isblocked"
)

func (a AccountsServer) Create(ctx context.Context, request *Accounts.CreateRequest) (*Accounts.CreateResponse, error) {
	ac := &Accounts.CreateResponse{
		Id:       uuid.New().String(),
		UserId:   request.UserId,
		IsCredit: request.IsCredit,
		Balance:  request.Balance,
		Currency: request.Currency,
	}

	query := qb.Insert(TableAccounts).
		Columns(ColumnsAccounts).
		Values(ac.Id, ac.UserId, ac.IsCredit, ac.Balance, ac.Currency, false).
		PlaceholderFormat(qb.Dollar).
		RunWith(DB)

	_, err := query.Exec()
	if err != nil {
		return nil, err
	}

	return ac, nil
}

func (a AccountsServer) GetDetails(ctx context.Context, request *Accounts.GetDetailsRequest) (*Accounts.GetDetailsResponse, error) {
	ac := &Accounts.GetDetailsResponse{
		Id:     request.Id,
		UserId: request.UserId,
	}

	query := qb.Select(ColumnsAccounts).
		From(TableAccounts).
		Where("(id = $?) AND (userid = $?)", ac.Id, ac.UserId).
		RunWith(DB)

	res, err := query.Query()
	if err != nil {
		return nil, err
	}

	res.Next()
	err = res.Scan(&ac.Id, &ac.UserId, &ac.IsCredit, &ac.Balance, &ac.Currency, &ac.IsBlocked)

	if err != nil {
		return nil, err
	}

	return ac, nil
}

func (a AccountsServer) Block(ctx context.Context, request *Accounts.BlockRequest) (*Accounts.BlockResponse, error) {
	query := qb.Update(TableAccounts).
		Set("isblocked=$1", true).
		Where("(id=$2)AND(userid=$3)", request.Id, request.UserId).
		PlaceholderFormat(qb.Dollar).
		RunWith(DB)

	_, err := query.Exec()
	if err != nil {
		return nil, err
	}
	return &Accounts.BlockResponse{Success: true}, nil
}

func (a AccountsServer) GetAll(ctx context.Context, request *Accounts.GetAllRequest) (*Accounts.GetAllResponse, error) {
	acs := make([]*Accounts.Account, 0)

	query := qb.Select("*").
		From(TableAccounts).
		RunWith(DB)

	res, err := query.Query()
	if err != nil {
		return nil, err
	}

	for res.Next() {
		ac := &Accounts.Account{}
		err = res.Scan(&ac.Id, &ac.UserId, &ac.IsCredit, &ac.Balance, &ac.IsBlocked, &ac.IsBlocked)
		if err != nil {
			return nil, err
		}
		acs = append(acs, ac)
	}

	return &Accounts.GetAllResponse{
		Accounts: acs,
	}, nil
}
