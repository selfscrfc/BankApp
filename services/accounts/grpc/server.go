package grpcaccountsserver

import (
	"context"
	"database/sql"
	qb "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/pkg/errors"
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

var MAX_CREDIT int32 = 50000

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

func (a AccountsServer) GetDetail(ctx context.Context, request *Accounts.GetDetailsRequest) (*Accounts.GetDetailsResponse, error) {
	ac := &Accounts.GetDetailsResponse{
		Id:     request.Id,
		UserId: request.UserId,
	}

	query := qb.Select(ColumnsAccounts).
		From(TableAccounts).
		Where("(id = $1) AND (userid = $2)", ac.Id, ac.UserId).
		PlaceholderFormat(qb.Dollar).
		RunWith(DB)

	res, err := query.Query()
	defer res.Close()

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
		Set("isblocked=$1", "true").
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
		err = res.Scan(&ac.Id, &ac.UserId, &ac.IsCredit, &ac.Balance, &ac.Currency, &ac.IsBlocked)
		if err != nil {
			return nil, err
		}
		acs = append(acs, ac)
	}

	return &Accounts.GetAllResponse{
		Accounts: acs,
	}, nil
}

func (a AccountsServer) RW(ctx context.Context, request *Accounts.RWRequest) (*Accounts.RWResponse, error) {
	query := qb.Select("*").
		From(TableAccounts).
		Where("(id=$1)AND(userid=$2)", request.AId, request.UId).
		RunWith(DB)

	res, err := query.Query()
	if err != nil {
		return nil, err
	}

	res.Next()
	defer res.Close()

	ac := &Accounts.Account{}
	err = res.Scan(&ac.Id, &ac.UserId, &ac.IsCredit, &ac.Balance, &ac.Currency, &ac.IsBlocked)
	if err != nil {
		return nil, err
	}

	if ac.IsBlocked {
		return nil, errors.New("account is blocked")
	}
	nb := int32(request.Amount) + ac.Balance
	if nb < 0 && ac.IsCredit == false {
		return nil, errors.New("insufficient funds")
	} else if ac.IsCredit == true && nb < -MAX_CREDIT {
		return nil, errors.New("credit cant overlap 50000")
	}

	query2 := qb.Update(TableAccounts).
		Set("balance", nb).
		Where("(id=$2)AND(userid=$3)", request.AId, request.UId).
		PlaceholderFormat(qb.Dollar).
		RunWith(DB)

	_, err = query2.Exec()
	if err != nil {
		return nil, err
	}

	return &Accounts.RWResponse{
		Amount:    int32(request.Amount),
		Id:        ac.Id,
		UserId:    ac.UserId,
		IsCredit:  ac.IsCredit,
		Balance:   ac.Balance,
		Currency:  ac.Currency,
		IsBlocked: ac.IsBlocked,
	}, nil
}
