package models

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("key not found")
)

type Models struct {
	KeyValue KeyValueModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		KeyValue: KeyValueModel{DB: db},
	}
}