package services

import (
	repo "github.com/SwanPoi/bmstu_rsoi_lab2/src/rental/repositories"
)

type IRentalService interface {

}

type Services struct {
	IRentalService
}

func NewServices(repo *repo.Repository) *Services {
	return &Services{
		IRentalService: NewRentalService(repo),
	}
}