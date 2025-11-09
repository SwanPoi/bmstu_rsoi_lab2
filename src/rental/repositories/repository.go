package repositories

import (
	"gorm.io/gorm"
)

type IRentalRepo interface {

}

type Repository struct {
	IRentalRepo
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		IRentalRepo: NewRentalPostgres(db),
	}
}