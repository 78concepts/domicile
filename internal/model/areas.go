package model

import (
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"time"
)

type Area struct {
	Id uint64 `json:"id"`
	Uuid uuid.UUID `json:"uuid"`
	DateCreated time.Time `json:"dateCreated"`
	Name string `json:"id"`
}
