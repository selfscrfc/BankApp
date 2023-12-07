package repository

import (
	"database/sql"
	qb "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/selfscrfc/PetBank/api/models"
	Accounts "github.com/selfscrfc/PetBankProtos/proto/Accounts"
	"strconv"
)

type AccountsRepo struct {
	DB *sql.DB
}

var MAX_CREDIT int32 = 50000

const (
	TableAccounts   = "accounts"
	ColumnsAccounts = "id, userid, iscredit, balance, currency, isblocked"
)

func (db AccountsRepo) CreateAccountQuery(ac *models.Account) (*models.Account, error) {
	query := qb.Insert(TableAccounts).
		Columns(ColumnsAccounts).
		Values(ac.Id, ac.UserId, ac.IsCredit, ac.Balance, ac.Currency, false).
		PlaceholderFormat(qb.Dollar).
		RunWith(db.DB)

	_, err := query.Exec()

	if err != nil {
		return nil, err
	}

	return ac, nil
}

func (db AccountsRepo) GetAllAccountsQuery(acs []*Accounts.Account) ([]*Accounts.Account, error) {
	query := qb.Select("*").
		From(TableAccounts).
		RunWith(db.DB)

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

	return acs, nil
}

func (db AccountsRepo) GetAccountDetailsQuery(ac *Accounts.Account) (*Accounts.Account, error) {
	query := qb.Select(ColumnsAccounts).
		From(TableAccounts).
		Where("(id = $1) AND (userid = $2)", ac.Id, ac.UserId).
		PlaceholderFormat(qb.Dollar).
		RunWith(db.DB)

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

	return ac, err
}

func (db AccountsRepo) BlockAccountQuery(ac *Accounts.Account) (bool, error) {
	query := qb.Update(TableAccounts).
		Set("isblocked", true).
		Where("(id=$2)AND(userid=$3)", ac.Id, ac.UserId).
		PlaceholderFormat(qb.Dollar).
		RunWith(db.DB)

	_, err := query.Exec()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (db AccountsRepo) RWQuery(aId, uId string, amount int64) (*Accounts.Account, error) {
	query := qb.Select(ColumnsAccounts).
		From(TableAccounts).
		Where("(id=$1)AND(userid=$2)", aId, uId).
		PlaceholderFormat(qb.Dollar).
		RunWith(db.DB)

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

	if ac.IsBlocked == true {
		return nil, errors.New("account is blocked")
	}

	nb := int32(amount) + ac.Balance

	if nb < 0 && ac.IsCredit == false {
		return nil, errors.New("insufficient funds")
	} else if ac.IsCredit == true && nb < -MAX_CREDIT {
		return nil, errors.New("credit cant overlap " + strconv.Itoa(int(MAX_CREDIT)))
	}

	query2 := qb.Update(TableAccounts).
		Set("balance", nb).
		Where("(id=$2)AND(userid=$3)", aId, uId).
		PlaceholderFormat(qb.Dollar).
		RunWith(db.DB)

	_, err = query2.Exec()
	if err != nil {
		return nil, err
	}

	return ac, err
}
