package services

import (
	"github.com/SwanPoi/bmstu_rsoi_lab2/src/car/models"
	repo "github.com/SwanPoi/bmstu_rsoi_lab2/src/car/repositories"
)

type ICarService interface {
	GetCars(page int, size int, showAll bool) (*models.PaginationResponse, error)
}

type Services struct {
	ICarService
}

func NewServices(repo *repo.Repository) *Services {
	return &Services{
		ICarService: NewCarService(repo),
	}
}