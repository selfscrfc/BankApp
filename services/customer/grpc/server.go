package grpccustomerserver

import (
	"context"
	"database/sql"
	qb "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/selfscrfc/PetBank/api/models"
	"github.com/selfscrfc/PetBank/utils"
	customers "github.com/selfscrfc/PetBankProtos/proto/Customers"
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
	id_, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, err
	}
	cs := &models.Customer{
		Id: id_,
	}

	query := qb.Select(ColumnsCustomer).
		From(TableCustomers).
		Where("id = $?", cs.Id.String()).
		RunWith(DB)

	res, err := query.Query()
	if err != nil {
		return nil, err
	}

	res.Next()
	err = res.Scan(&cs.Id, &cs.FullName, &cs.TimeCreated, &cs.Login, &cs.Password, &cs.IsBlocked)

	if err != nil {
		return nil, err
	}

	cs.Password = ""

	return &customers.GetDetailsResponse{
		Id:        cs.Id.String(),
		FullName:  cs.FullName,
		Time:      cs.TimeCreated.Unix(),
		Login:     cs.Login,
		IsBlocked: cs.IsBlocked,
	}, nil
}

func (c CustomerServer) Block(ctx context.Context, request *customers.BlockRequest) (*customers.BlockResponse, error) {
	id_, err := uuid.Parse(request.BlockId)
	if err != nil {
		return nil, err
	}
	cs := &models.Customer{
		Id: id_,
	}

	query := qb.Update(TableCustomers).
		Set("isblocked", true).
		Where("id = $?", cs.Id.String()).
		RunWith(DB)

	res, err := query.Query()
	if err != nil {
		return nil, err
	}

	res.Next()
	err = res.Scan(&cs.Id, &cs.FullName, &cs.TimeCreated, &cs.Login, &cs.Password, &cs.IsBlocked)
	if err != nil {
		return nil, err
	}

	cs.Password = ""

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
