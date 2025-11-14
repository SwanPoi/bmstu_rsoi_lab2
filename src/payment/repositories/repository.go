package repositories

import (
	"gorm.io/gorm"

	"github.com/SwanPoi/bmstu_rsoi_lab2/src/payment/models"
)

type IPaymentRepo interface {
	GetPaymentByUid(string) (*models.PaymentResponse, error)
	GetPaymentsByUids(uids []string) ([]models.PaymentResponse, error)
}

type Repository struct {
	IPaymentRepo
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		IPaymentRepo: NewPaymentPostgres(db),
	}
}