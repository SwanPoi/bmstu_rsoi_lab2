package services

import (
	"github.com/SwanPoi/bmstu_rsoi_lab2/src/car/converters"
	"github.com/SwanPoi/bmstu_rsoi_lab2/src/car/models"
	repo "github.com/SwanPoi/bmstu_rsoi_lab2/src/car/repositories"
)

type CarService struct {
	repo *repo.Repository
}

func NewCarService(repo *repo.Repository) *CarService {
	return &CarService{repo: repo}
}

func (s *CarService) GetCars(page int, size int, showAll bool) (*models.PaginationResponse, error) {
	offset := page * size

	cars, total, err := s.repo.GetCars(offset, size, showAll)

	if err != nil {
		return nil, err
	}

	carsResponse := converters.CarResponsesFromCars(cars)

	paginationResponse := &models.PaginationResponse{
		Page: page,
		PageSize: size,
		TotalElements: total,
		Items: carsResponse,
	}

	return paginationResponse, nil
}