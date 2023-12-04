package postgres

import (
	"context"
	"database/sql"
	qb "github.com/Masterminds/squirrel"
	"github.com/opentracing/opentracing-go"
	"github.com/selfscrfc/PetBank/api/models"
)

const (
	TableCustomers  = "customers"
	ColumnsCustomer = "id, fullname, time, login, password, isblocked"
)

type customerRepo struct {
	db *sql.DB
}

func NewCustomerRepo(db *sql.DB) *customerRepo {
	return &customerRepo{db: db}
}

func (cr *customerRepo) Create(ctx context.Context, c *models.Customer) (*models.Customer, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "newsRepo.Create")
	defer span.Finish()

	query := qb.Insert(TableCustomers).
		Columns(ColumnsCustomer).
		Values(c.Id, c.FullName, c.TimeCreated, c.Login, c.Password, "false").
		PlaceholderFormat(qb.Question).
		RunWith(cr.db)

	_, err := query.ExecContext(ctx)
	if err != nil {
		return nil, err
	}

	return c, nil
}
