package services

import (
	"github.com/SwanPoi/bmstu_rsoi_lab2/src/payment/models"
	repo "github.com/SwanPoi/bmstu_rsoi_lab2/src/payment/repositories"
)

type PaymentService struct {
	repo *repo.Repository
}

func NewPaymentService(repo *repo.Repository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) GetPaymentByUid(uid string) (*models.PaymentResponse, error) {
	return s.repo.GetPaymentByUid(uid)
}

func (s *PaymentService) GetPaymentsByUids(uids []string) ([]models.PaymentResponse, error) {
	return s.repo.GetPaymentsByUids(uids)
}
