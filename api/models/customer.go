package models

import (
	"github.com/google/uuid"
	"time"
)

type Customer struct {
	Id          uuid.UUID `json:"id" validate:"uuid"`
	TimeCreated time.Time `json:"timecreated"`
	FullName    string    `json:"fullname"`
	IsBlocked   bool      `json:"blocked"`
	Login       string    `json:"login"`
	Password    string    `json:"password"`
}
