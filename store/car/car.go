package car

import (
	"context"
	"database/sql"

	"github.com/shafnybuilds/car_management_sys/models"
)

type Store struct {
	db *sql.DB
}

// constructor function for creating a new instance
func new(db *sql.DB) Store {
	return Store{db: db}
}

// methods
func (s Store) GetCarById(ctx context.Context, id string) (models.Car, error) {

}

func (s Store) GetCarByBrand(ctx context.Context, brand string, isEngine bool) {

}

func (s Store) CreateCar(ctx context.Context, carReq *models.CarRequest) (models.Car, error) {

}

func (s Store) UpdateCar(ctx context.Context, id string, carReq *models.CarRequest) (models.Car, error) {

}

func (s Store) DeleteCar(ctx context.Context, id string, carReq *models.Engine) (models.Car, error) {

}
