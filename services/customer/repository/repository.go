package repository

import (
	"database/sql"
	qb "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/selfscrfc/PetBank/utils"
	customers "github.com/selfscrfc/PetBankProtos/proto/Customers"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type CustomerRepo struct {
	DB *sql.DB
}

const (
	TableCustomers  = "customers"
	ColumnsCustomer = "id, fullname, time, login, password, isblocked"
)

func (db CustomerRepo) CreateCustomerQuery(cs *customers.CustomerEntity, pass string) (*customers.CustomerEntity, error) {
	cs.Id = uuid.New().String()
	cs.Time = time.Now().Unix()

	query := qb.Insert(TableCustomers).
		Columns(ColumnsCustomer).
		Values(cs.Id, cs.FullName, cs.Time, cs.Login, utils.GeneratePassword(pass), cs.IsBlocked).
		PlaceholderFormat(qb.Dollar).
		RunWith(db.DB)

	_, err := query.Exec()
	if err != nil {
		return nil, err
	}

	return cs, nil
}

func (db CustomerRepo) GetAllCustomersQuery(csa []*customers.CustomerEntity) ([]*customers.CustomerEntity, error) {
	mock := ""
	mockp := &mock

	query := qb.Select("*").
		From(TableCustomers).
		RunWith(db.DB)

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
		*mockp = ""
	}

	return csa, nil
}

func (db CustomerRepo) GetCustomerDetailsQuery(cs *customers.CustomerEntity) (*customers.CustomerEntity, error) {
	query := qb.Select(ColumnsCustomer).
		From(TableCustomers).
		Where("id=$1", cs.Id).
		PlaceholderFormat(qb.Dollar).
		RunWith(db.DB)

	res, err := query.Query()
	if err != nil {
		return nil, err
	}
	pass := ""
	res.Next()
	err = res.Scan(&cs.Id, &cs.FullName, &cs.Time, &cs.Login, &pass, &cs.IsBlocked)
	pass = ""

	return cs, err
}

func (db CustomerRepo) BlockAccountQuery(id string) (bool, error) {
	query := qb.Update(TableCustomers).
		Set("isblocked", true).
		Where("id=$2", id).
		PlaceholderFormat(qb.Dollar).
		RunWith(db.DB)

	if _, err := query.Exec(); err != nil {
		return false, err
	}

	return true, nil
}

func (db CustomerRepo) SignInQuery(log, pass string) (*customers.CustomerEntity, error) {
	query := qb.Select(ColumnsCustomer).
		From(TableCustomers).
		Where("login=$1", log).
		PlaceholderFormat(qb.Dollar).
		RunWith(db.DB)

	res, err := query.Query()
	defer res.Close()

	if err != nil {
		return nil, err
	}

	cs := &customers.CustomerEntity{}

	hashPass := ""

	if !res.Next() {
		return nil, errors.New("User not found")
	}

	if err = res.Scan(&cs.Id, &cs.FullName, &cs.Time, &cs.Login, &hashPass, &cs.IsBlocked); err != nil {
		return nil, err
	}

	if cs.IsBlocked == true {
		return nil, errors.New("account is blocked")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(pass)); err != nil {
		return nil, err
	}

	return cs, err
}
