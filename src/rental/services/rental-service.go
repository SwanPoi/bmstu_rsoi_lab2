package services

import (
	"time"

	"github.com/SwanPoi/bmstu_rsoi_lab2/src/rental/models"
	repo "github.com/SwanPoi/bmstu_rsoi_lab2/src/rental/repositories"
	"github.com/google/uuid"
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

func (s *RentalService) CreateRental(rentalReq models.RentCreation) (*models.RentalResponse, error) {
	dateFrom, err := time.Parse(time.RFC3339, rentalReq.DateFrom)
    if err != nil {
        return nil, err
    }

    dateTo, err := time.Parse(time.RFC3339, rentalReq.DateTo)
    if err != nil {
        return nil, err
    }

	rental := models.Rental{
		RentalUID: uuid.New().String(),
		Username: rentalReq.Username,
		CarUID: rentalReq.CarUID,
		PaymentUID: rentalReq.PaymentUID,
		Status: "IN_PROGRESS",
		DateFrom: dateFrom,
		DateTo: dateTo,
	}

	if err := s.repo.CreateRental(rental); err == nil {
		response := models.RentalResponse{
			RentalUID: rental.RentalUID,
			PaymentUID: rental.PaymentUID,
			CarUID: rental.CarUID,
			DateFrom: rental.DateFrom,
			DateTo: rental.DateTo,
			Status: rental.Status,
		}

		return &response, nil
	} else {
		return nil, err
	}
}

func (s *RentalService) UpdateRental(rental models.RentalUpsert, uid string, username string) (error) {
	validStatuses := map[string]bool{
        "IN_PROGRESS": true,
        "FINISHED":    true,
        "CANCELED":    true,
    }

	if !validStatuses[rental.Status] {
        return models.InvalidStatus
    }

	return s.repo.UpdateRental(rental, uid, username)
}