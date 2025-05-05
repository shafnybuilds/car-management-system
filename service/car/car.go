package car

import (
	"context"

	"github.com/shafnybuilds/car_management_sys/models"
	"github.com/shafnybuilds/car_management_sys/store"
)

// in this case we are actually using dependency injection here
type CarService struct {
	store store.CarStoreInterface
}

func NewCarService(store store.CarStoreInterface) *CarService {
	return &CarService{
		store: store,
	}
}

func (s *CarService) GetCarByID(ctx context.Context, id string) (*models.Car, error) {
	car, err := s.store.GetCarById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &car, nil
}

func (s *CarService) GetCarsByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error) {
	cars, err := s.store.GetCarByBrand(ctx, brand, isEngine)
	if err != nil {
		return nil, err
	}

	return cars, nil
}

func (s *CarService) CreateCar(ctx context.Context, car *models.CarRequest) (*models.Car, error) {
	if err := models.ValidationRequest(*car); err != nil {
		return nil, err
	}

	createdCar, err := s.store.CreateCar(ctx, car)
	if err != nil {
		return nil, err
	}

	return &createdCar, nil
}
