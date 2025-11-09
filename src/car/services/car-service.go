package services

import (
	repo "github.com/SwanPoi/bmstu_rsoi_lab2/src/car/repositories"
)

type CarService struct {
	repo *repo.Repository
}

func NewCarService(repo *repo.Repository) *CarService {
	return &CarService{repo: repo}
}