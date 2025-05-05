package store

import (
	"context"

	"github.com/shafnybuilds/car_management_sys/models"
)

type CarStoreInterface interface {
	GetCarById(ctx context.Context, id string) (models.Car, error)

	GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error)

	CreateCar(ctx context.Context, carReq *models.CarRequest) (models.Car, error)

	UpdateCar(ctx context.Context, id string, carReq *models.CarRequest) (models.Car, error)

	DeleteCar(ctx context.Context, id string, carReq *models.Engine) (models.Car, error)
}
