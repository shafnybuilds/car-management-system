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

	DeleteCar(ctx context.Context, id string) (models.Car, error)
}

type EngineStoreInterface interface {
	GetEngineById(ctx context.Context, id string) (models.Engine, error)

	EngineCreate(ctx context.Context, engineReq *models.EngineRequest) (models.Engine, error)

	EngineUpdate(ctx context.Context, id string, engineReq *models.EngineRequest) (models.Engine, error)

	DeleteEngine(ctx context.Context, id string) (models.Engine, error)
}
