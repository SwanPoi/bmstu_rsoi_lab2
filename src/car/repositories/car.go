package repositories

import "gorm.io/gorm"

type CarPostgres struct {
	db *gorm.DB
}

func NewCarPostgres(db *gorm.DB) *CarPostgres {
	return &CarPostgres{db: db}
}