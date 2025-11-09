package services

import (
	repo "github.com/SwanPoi/bmstu_rsoi_lab2/src/payment/repositories"
)

type IPaymentService interface {

}

type Services struct {
	IPaymentService
}

func NewServices(repo *repo.Repository) *Services {
	return &Services{
		IPaymentService: NewPaymentService(repo),
	}
}