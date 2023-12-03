package models

type Customer struct {
	Id              int
	FullName        string
	YearCreated     int
	DatetimeCreated int
	IsBlocked       bool
	Login           string
	Password        string
}
