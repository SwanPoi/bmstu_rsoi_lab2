package services

import (
	repo "github.com/SwanPoi/bmstu_rsoi_lab2/src/car/repositories"
)

type ICarService interface {

}

type Services struct {
	ICarService
}

func NewServices(repo *repo.Repository) *Services {
	return &Services{
		ICarService: NewCarService(repo),
	}
}