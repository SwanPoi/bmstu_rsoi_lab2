package repositories

import "gorm.io/gorm"

type PaymentPostgres struct {
	db *gorm.DB
}

func NewPaymentPostgres(db *gorm.DB) *PaymentPostgres {
	return &PaymentPostgres{db: db}
}