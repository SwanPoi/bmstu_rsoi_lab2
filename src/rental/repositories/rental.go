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

func (r *RentalPostgres) CreateRental(rental models.Rental) (error) {
	return r.DB.Create(rental).Error
}

func (r *RentalPostgres) UpdateRental(rental models.RentalUpsert, uid string, username string) (error) {
	result := r.DB.Model(&models.Rental{}).
					Where("rental_uid = ? AND username = ?", uid, username).
					Update("status", rental.Status)
	
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return models.ErrorNotFound
	}

	return nil
}