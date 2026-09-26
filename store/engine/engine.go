package engine

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
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
	var engine models.Engine

	err := e.db.QueryRowContext(ctx, "SELECT id, displacement, no_of_cylinders, car_range FROM engine WHERE id=$1", id).Scan(
		&engine.EngineID,
		&engine.Displacement,
		&engine.NoOfCylinders,
		&engine.CarRange,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return engine, nil
		}
		return engine, err
	}

	return engine, err
}

func (e EngineStore) EngineCreate(ctx context.Context, engineReq *models.EngineRequest) (models.Engine, error) {

	engineID := uuid.New()

	_, err := e.db.ExecContext(ctx,
		"INSERT INTO engine (id, displacement, no_of_cylinders, car_range) VALUES ($1, $2, $3, $4)",
		engineID, engineReq.Displacement, engineReq.NoOfCylinders, engineReq.CarRange)

	if err != nil {
		return models.Engine{}, err
	}

	engine := models.Engine{
		EngineID:      engineID,
		Displacement:  engineReq.Displacement,
		NoOfCylinders: engineReq.NoOfCylinders,
		CarRange:      engineReq.CarRange,
	}

	return engine, nil

}

func (e EngineStore) EngineUpdate(ctx context.Context, id string, engineReq *models.EngineRequest) (models.Engine, error) {

	engineID, err := uuid.Parse(id)
	if err != nil {
		return models.Engine{}, fmt.Errorf("Invalied Engine ID: %w", err)
	}

	result, err := e.db.ExecContext(ctx,
		"UPDATE engine SET displacement = $1, no_of_cylinders = $2, car_range = $3 WHERE id = $4", engineReq.Displacement, engineReq.NoOfCylinders, engineReq.CarRange, engineID)

	if err != nil {
		return models.Engine{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Engine{}, err
	}
	if rowsAffected == 0 {
		return models.Engine{}, errors.New("No rows were modified!")
	}

	engine := models.Engine{
		EngineID:      engineID,
		Displacement:  engineReq.Displacement,
		NoOfCylinders: engineReq.NoOfCylinders,
		CarRange:      engineReq.CarRange,
	}

	return engine, nil

}

func (e EngineStore) DeleteEngine(ctx context.Context, id string) (models.Engine, error) {
	var engine models.Engine
	err := e.db.QueryRowContext(ctx, "DELETE FROM engine WHERE id = $1 RETURNING id, displacement, no_of_cylinders, car_range", id).Scan(
		&engine.EngineID, &engine.Displacement, &engine.NoOfCylinders, &engine.CarRange,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return engine, nil
	}
	return engine, err
}
