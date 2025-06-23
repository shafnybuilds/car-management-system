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

// constructor
func NewCarService(store store.CarStoreInterface) *CarService {
	return &CarService{
		store: store,
	}
}

// methods
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

func (s *CarService) UpdateCar(ctx context.Context, id string, carReq *models.CarRequest) (*models.Car, error) {
	if err := models.ValidationRequest(*carReq); err != nil {
		return nil, err
	}

	updatedCar, err := s.store.UpdateCar(ctx, id, carReq)
	if err != nil {
		return nil, err
	}

	return &updatedCar, nil
}

func (s *CarService) DeleteCar(ctx context.Context, id string) (*models.Car, error) {
	deletedCar, err := s.store.DeleteCar(ctx, id)
	if err != nil {
		return nil, err
	}

	return &deletedCar, nil
}