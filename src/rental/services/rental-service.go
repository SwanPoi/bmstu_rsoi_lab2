package services

import (
	repo "github.com/SwanPoi/bmstu_rsoi_lab2/src/rental/repositories"
)

type RentalService struct {
	repo *repo.Repository
}

func NewRentalService(repo *repo.Repository) *RentalService {
	return &RentalService{repo: repo}
}