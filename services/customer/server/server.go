package cServer

import (
	"context"
	"database/sql"
	qb "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "newsRepo.Create")
	defer span.Finish()

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
		Values(cs.Id, cs.FullName, cs.TimeCreated, cs.Login, cs.Password, cs.IsBlocked).
		PlaceholderFormat(qb.Question).
		RunWith(DB)

	_, err := query.ExecContext(ctx)
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
