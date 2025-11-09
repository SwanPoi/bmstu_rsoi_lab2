package services

import (
	repo "github.com/SwanPoi/bmstu_rsoi_lab2/src/payment/repositories"
)

type PaymentService struct {
	repo *repo.Repository
}

func NewPaymentService(repo *repo.Repository) *PaymentService {
	return &PaymentService{repo: repo}
}