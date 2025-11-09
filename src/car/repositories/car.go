package repositories

import (
	"github.com/SwanPoi/bmstu_rsoi_lab2/src/car/models"
	"gorm.io/gorm"
)

type CarPostgres struct {
	db *gorm.DB
}

func NewCarPostgres(db *gorm.DB) *CarPostgres {
	return &CarPostgres{db: db}
}

func (r *CarPostgres) GetCars(offset int, limit int, showAll bool) ([]models.Car, int, error) {
	var total int64
	var cars []models.Car

	query := r.db.Model(&models.Car{})

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