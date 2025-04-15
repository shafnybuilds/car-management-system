package engine

import (
	"context"
	"database/sql"

	"github.com/shafnybuilds/car_management_sys/models"
)

type EngineStore struct {
	db *sql.DB
}

// constructor for initialization
func New(db *sql.DB) *EngineStore {
	return &EngineStore{db: db}
}

func (e EngineStore) GetEngineById(ctx context.Context, id string) (models.Engine, error) {

}

func (e EngineStore) CreatEngine(ctx context.Context, engineReq *models.EngineRequest) (models.Engine, error) {

}

func (e EngineStore) EngineUpdate(ctx context.Context, id string, engineReq *models.EngineRequest) (models.Engine, error) {

}

func (e EngineStore) DeleteEngine(ctx context.Context, id string) (models.Engine, error) {

}
