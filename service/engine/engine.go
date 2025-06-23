package engine

import (
	"context"

	"github.com/shafnybuilds/car_management_sys/models"
	"github.com/shafnybuilds/car_management_sys/store"
)

type EngineService struct {
	store store.EngineStoreInterface
}

// constructor
func NewEngineService(store store.EngineStoreInterface) *EngineService {
	return &EngineService{
		store: store,
	}
}

// methods
func (s *EngineService) GetEngineById(ctx context.Context, id string) (*models.Engine, error) {
	engine, err := s.store.GetEngineById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &engine, nil
}

func (s *EngineService) CreateEngine(ctx context.Context, engineReq *models.EngineRequest) (*models.Engine, error) {
	if err := models.ValidateEngineRequest(*engineReq); err != nil {
		return nil, err
	}

	createdEngine, err := s.store.EngineCreate(ctx, engineReq)
	if err != nil {
		return nil, err
	}

	return &createdEngine, nil
}

func (s *EngineService) UpdateEngine(ctx context.Context, id string, engineReq *models.EngineRequest) (*models.Engine, error) {
	if err := models.ValidateEngineRequest(*engineReq); err != nil {
		return nil, err
	}

	updatedEngine, err := s.store.EngineUpdate(ctx, id, engineReq)
	if err != nil {
		return nil, err
	}

	return &updatedEngine, nil

}

func (s *EngineService) DeleteEngine(ctx context.Context, id string) (*models.Engine, error) {
	deletedEngine, err := s.store.DeleteEngine(ctx, id)
	if err != nil {
		return nil, err
	}

	// DeleteEngine func returns two params, models.Engine and error
	// since we already handled the error we can return nil
	return &deletedEngine, nil
}
