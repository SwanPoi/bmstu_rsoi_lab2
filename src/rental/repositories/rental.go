package repositories

import "gorm.io/gorm"

type RentalPostgres struct {
	db *gorm.DB
}

func NewRentalPostgres(db *gorm.DB) *RentalPostgres {
	return &RentalPostgres{db: db}
}