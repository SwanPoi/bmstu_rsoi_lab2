package repositories

import (
	"gorm.io/gorm"
)

type ICarRepo interface {

}

type Repository struct {
	ICarRepo
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		ICarRepo: NewCarPostgres(db),
	}
}