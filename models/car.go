package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Car struct {
	ID        uuid.UUID `json "id"`
	Name      string    `json: "name"`
	Year      string    `json: "year`
	Brand     string    `json: "brand"`
	FuelType  string    `json: "fuel_type"`
	Engine    Engine    `json: "engine"`
	Price     float64   `json: "price"`
	CreatedAt time.Time `json: "created_at"`
	UpdatedAt time.Time `json: "Updated_at"`
}

// Client request body
type CarRequest struct {
	Name     string  `json: "name"`
	Year     string  `json: "year`
	Brand    string  `json: "brand"`
	FuelType string  `json: "fuel_type"`
	Engine   Engine  `json: "engine"`
	Price    float64 `json: "price"`
}

// validation functions for the car
func validateName(name string) error {
	if name == "" {
		return errors.New("Name is Required")
	}
	return nil
}

func validateYear(year string) error {
	if year == "" {
		return errors.New("Year is Required")
	}
	_, err := strconv.Atoi(year)
	if err != nil {
		return errors.New("year must be a valid number")
	}

	currentYear := time.Now().Year()
	YearInt, _ := strconv.Atoi(year)
	if YearInt < 1886 || YearInt > currentYear {
		return errors.New("Year must be a value between 1886 and current year")
	}

	return nil
}

func validateBrand(brand string) error {
	if brand == "" {
		return errors.New("Brand is Required")
	}
	return nil
}
