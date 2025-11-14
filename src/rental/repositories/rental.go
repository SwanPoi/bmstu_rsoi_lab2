package repositories

import (
	"errors"

	"github.com/SwanPoi/bmstu_rsoi_lab2/src/rental/models"
	"gorm.io/gorm"
)

type RentalPostgres struct {
	DB *gorm.DB
}

func NewRentalPostgres(db *gorm.DB) *RentalPostgres {
	return &RentalPostgres{DB: db}
}

func (r *RentalPostgres) GetRentalByUid(uid string) (*models.Rental, error) {
	var rental models.Rental

	if err := r.DB.Where("rental_uid = ?", uid).First(&rental).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrorNotFound
		}

		return nil, err
	}

	return &rental, nil
}

func (r *RentalPostgres) GetUserRentals(username string) ([]models.RentalResponse, error) {
	var rentals []models.RentalResponse

	if err := r.DB.Omit("id", "username").Where("username = ?", username).Find(&rentals).Error; err != nil {
		return nil, err
	}

	return rentals, nil
}
