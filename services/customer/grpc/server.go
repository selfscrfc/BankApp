package grpccustomerserver

import (
	"context"
	"database/sql"
	qb "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/selfscrfc/PetBank/api/models"
	"github.com/selfscrfc/PetBank/utils"
	customers "github.com/selfscrfc/PetBankProtos/proto/Customers"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var DB *sql.DB

type CustomerServer struct {
	customers.UnimplementedCustomerServer
}

const (
	TableCustomers  = "customers"
	ColumnsCustomer = "id, fullname, time, login, password, isblocked"
)

type customerRepo struct {
	db *sql.DB
}

func (c CustomerServer) Create(ctx context.Context, request *customers.CreateRequest) (*customers.CreateResponse, error) {
	cs := &models.Customer{
		Id:          uuid.New(),
		TimeCreated: time.Now(),
		FullName:    request.FullName,
		Login:       request.Login,
		Password:    utils.GeneratePassword(request.Password),
		IsBlocked:   false,
	}

	query := qb.Insert(TableCustomers).
		Columns(ColumnsCustomer).
		Values(cs.Id.String(), cs.FullName, cs.TimeCreated.Unix(), cs.Login, cs.Password, cs.IsBlocked).
		PlaceholderFormat(qb.Dollar).
		RunWith(DB)

	_, err := query.Exec()
	if err != nil {
		return nil, err
	}

	return &customers.CreateResponse{
		Id:       cs.Id.String(),
		FullName: cs.FullName,
		Time:     cs.TimeCreated.Unix(),
		Login:    cs.Login,
		Password: "",
	}, nil
}

func (c CustomerServer) GetDetails(ctx context.Context, request *customers.GetDetailsRequest) (*customers.GetDetailsResponse, error) {
	query := qb.Select(ColumnsCustomer).
		From(TableCustomers).
		Where("id=$1", request.Id).
		PlaceholderFormat(qb.Dollar).
		RunWith(DB)

	res, err := query.Query()
	if err != nil {
		return nil, err
	}

	id_, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, err
	}
	cs := &models.Customer{
		Id: id_,
	}
	time_ := cs.TimeCreated.Unix()
	res.Next()
	err = res.Scan(&cs.Id, &cs.FullName, &time_, &cs.Login, &cs.Password, &cs.IsBlocked)

	if err != nil {
		return nil, err
	}

	cs.Password = ""

	return &customers.GetDetailsResponse{
		Id:        cs.Id.String(),
		FullName:  cs.FullName,
		Time:      time_,
		Login:     cs.Login,
		IsBlocked: cs.IsBlocked,
	}, nil
}

func (c CustomerServer) Block(ctx context.Context, request *customers.BlockRequest) (*customers.BlockResponse, error) {
	query := qb.Update(TableCustomers).
		Set("isblocked=$1", true).
		Where("id=$1", request.BlockId).
		PlaceholderFormat(qb.Dollar).
		RunWith(DB)

	if _, err := query.Exec(); err != nil {
		return nil, err
	}

	return &customers.BlockResponse{
		Success: true,
	}, nil
}

func (c CustomerServer) GetAll(ctx context.Context, request *customers.GetAllRequest) (*customers.GetAllResponse, error) {
	csa := make([]*customers.CustomerEntity, 0)
	mock := ""
	mockp := &mock

	query := qb.Select("*").
		From(TableCustomers).
		RunWith(DB)

	res, err := query.Query()
	defer res.Close()
	if err != nil {
		return nil, err
	}

	for res.Next() {
		cs := &customers.CustomerEntity{}
		err = res.Scan(&cs.Id, &cs.FullName, &cs.Time, &cs.Login, mockp, &cs.IsBlocked)
		if err != nil {
			return nil, err
		}
		csa = append(csa, cs)
	}
	*mockp = ""

	return &customers.GetAllResponse{
		Customers: csa,
	}, nil
}

func (c CustomerServer) SignIn(ctx context.Context, request *customers.SignInRequest) (*customers.SignInResponse, error) {
	query := qb.Select(ColumnsCustomer).
		From(TableCustomers).
		Where("login=$1", request.Login).
		PlaceholderFormat(qb.Dollar).
		RunWith(DB)

	res, err := query.Query()
	defer res.Close()

	if err != nil {
		return nil, err
	}

	resp := &customers.SignInResponse{}
	pass := ""

	if !res.Next() {
		return nil, errors.New("User not found")
	}

	if err = res.Scan(&resp.Id, &resp.FullName, &resp.Time, &resp.Login, &pass, &resp.IsBlocked); err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(request.Password)); err != nil {
		return nil, err
	}

	return resp, nil
}
