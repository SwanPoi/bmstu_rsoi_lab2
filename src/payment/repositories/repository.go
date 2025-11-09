package repositories

import (
	"gorm.io/gorm"
)

type IPaymentRepo interface {

}

type Repository struct {
	IPaymentRepo
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		IPaymentRepo: NewPaymentPostgres(db),
	}
}