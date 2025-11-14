package services

import (
	"github.com/SwanPoi/bmstu_rsoi_lab2/src/rental/models"
	repo "github.com/SwanPoi/bmstu_rsoi_lab2/src/rental/repositories"
)

type RentalService struct {
	repo *repo.Repository
}

func NewRentalService(repo *repo.Repository) *RentalService {
	return &RentalService{repo: repo}
}

func (s *RentalService) GetUserRentalByUid(uid string, username string) (*models.RentalResponse, error) {
	rental, err := s.repo.GetRentalByUid(uid)

	if err != nil {
		return nil, err
	}

	if rental.Username != username {
		return nil, models.Forbidden
	}

	rentalResponse := models.RentalResponse{
		RentalUID: rental.RentalUID,
		PaymentUID: rental.PaymentUID,
		CarUID: rental.CarUID,
		DateFrom: rental.DateFrom,
		DateTo: rental.DateTo,
		Status: rental.Status,
	}

	return &rentalResponse, nil
}

func (s *RentalService) GetUserRentals(username string) ([]models.RentalResponse, error) {
	return s.repo.GetUserRentals(username)
}