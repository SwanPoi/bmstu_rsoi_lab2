package repositories

import (
	"errors"

	"github.com/SwanPoi/bmstu_rsoi_lab2/src/payment/models"
	"gorm.io/gorm"
)

type PaymentPostgres struct {
	DB *gorm.DB
}

func NewPaymentPostgres(db *gorm.DB) *PaymentPostgres {
	return &PaymentPostgres{DB: db}
}

func (r *PaymentPostgres) GetPaymentByUid(uid string) (*models.PaymentResponse, error) {
	var payment models.PaymentResponse

	if err := r.DB.Select("payment_uid", "status", "price").Where("payment_uid = ?", uid).First(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrorNotFound
		}

		return nil, err
	}

	return &payment, nil
}

func (r *PaymentPostgres) GetPaymentsByUids(uids []string) ([]models.PaymentResponse, error) {
	var payments []models.PaymentResponse

	if err := r.DB.Select("payment_uid", "status", "price").Where("payment_uid IN ?", uids).Find(&payments).Error; err != nil {
		return nil, err
	}

	return payments, nil
}
