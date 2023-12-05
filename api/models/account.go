package models

import (
	"github.com/google/uuid"
)

type Currency int

const (
	RUBLES = iota
	DOLLARS
	EURO
)

type Account struct {
	Id       uuid.UUID
	UserId   uuid.UUID
	IsCredit bool
	Balance  int32
	Currency Currency
}
