package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/selfscrfc/PetBank/api/models"
)

type customerRepo struct {
	db *sqlx.DB
}

func newCustomerRepo(db *sqlx.DB) *customerRepo {
	return &customerRepo{db: db}
}

func Create(ctx context.Context, customer *models.Customer) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "newsRepo.Create")
	defer span.Finish()

	customer := models.Customer{}

}
