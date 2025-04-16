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
	var car models.Car

	query := `SELECT c.id, c.name, c.year, c.brand, c.furl_type, c.engine_id, c.price, c.created_at, c.updated_at, e.id, e.displacement, e.no_of_cylinders, e.car_range FROM car c LEFT JOIN engine e ON c.engine_id = e.id WHERE c.id = $1`

	row := s.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&car.ID,
		&car.Name,
		&car.Year,
		&car.Brand,
		&car.FuelType,
		&car.Engine.EngineID,
		&car.Price,
		&car.CreatedAt,
		&car.UpdatedAt,
		&car.Engine.EngineID,
		&car.Engine.Displacement,
		&car.Engine.NoOfCylinders,
		&car.Engine.CarRange,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return car, nil
		}
		return car, err
	}

	return car, nil
}

func (s Store) GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error) {
	var cars []models.Car
	var query string
	if isEngine {
		query =  `SELECT c.id, c.name, c.year, c.brand, c.fuel_type, c.engie_id, c.price, c.created_at, c.updated_at, e.id, e.displacement, e.no_of_cylinders, e.car_range FROM car c LEFT JOIN engine ON c.engine_id = e.id WHERE c.brand = $1`
	} else {
		query = `SELECT id, name, year, brand, fuel_type, engie_id, price, created_at, updated_at FROM car WHERE brand = $1`
	}

	rows, err := s.db.QueryContext(ctx, query, brand)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next(){
		
	}
}

func (s Store) CreateCar(ctx context.Context, carReq *models.CarRequest) (models.Car, error) {

}

func (s Store) UpdateCar(ctx context.Context, id string, carReq *models.CarRequest) (models.Car, error) {

}

func (s Store) DeleteCar(ctx context.Context, id string, carReq *models.Engine) (models.Car, error) {

}
