package repositories

import (
	"errors"

	"github.com/SwanPoi/bmstu_rsoi_lab2/src/car/models"
	"gorm.io/gorm"
)

type CarPostgres struct {
	DB *gorm.DB
}

func NewCarPostgres(db *gorm.DB) *CarPostgres {
	return &CarPostgres{DB: db}
}

func (r *CarPostgres) GetCars(offset int, limit int, showAll bool) ([]models.Car, int, error) {
	var total int64
	var cars []models.Car

	query := r.DB.Model(&models.Car{})

	if !showAll {
		query = query.Where("availability = ?", true)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&cars).Error; err != nil {
		return nil, 0, err
	}

	return cars, int(total), nil
}

func (r *CarPostgres) GetCarByUid(uid string) (*models.Car, error) {
	var car models.Car

	if err := r.DB.Where("car_uid = ?", uid).First(&car).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrorNotFound
		}

		return nil, err
	}

	return &car, nil
}

func (r *CarPostgres) GetCarsByUids(uids []string) ([]models.Car, error) {
	var cars []models.Car

	if err := r.DB.Where("car_uid IN ?", uids).Find(&cars).Error; err != nil {
		return nil, err
	}

	return cars, nil
}