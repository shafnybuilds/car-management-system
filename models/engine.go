package models

import (
	"errors"

	"github.com/google/uuid"
)

type Engine struct {
	EngineID      uuid.UUID `json: "engine_id"`
	Displacement  int64     `json: "displacement"`
	NoOfCylinders int64     `json: "no_of_cylinders`
	CarRange      int64     `json: "car_range"`
}

type EngineRequest struct {
	Displacement  int64 `json: "displacement"`
	NoOfCylinders int64 `json: "no_of_cylinders`
	CarRange      int64 `json: "car_range"`
}

// parent validation func
func ValidateEngineRequest(EngReq EngineRequest) error {
	if err := validateDisplacement(EngReq.Displacement); err != nil {
		return err
	}
	if err := validateNoOfCylinders(EngReq.NoOfCylinders); err != nil {
		return err
	}
	if err := validateDisplacement(EngReq.CarRange); err != nil {
		return err
	}
	return nil
}

func validateDisplacement(displacement int64) error {
	if displacement <= 0 {
		return errors.New("Displacement must be greater than zero")
	}
	return nil
}

func validateNoOfCylinders(NoOfCylinders int64) error {
	if NoOfCylinders <= 0 {
		return errors.New("No of Cylinders must be greater thatn zero")
	}
	return nil
}

func validateCarRange(CarRange int64) error {
	if CarRange <= 0 {
		return errors.New("Car Range must be greater than zero")
	}
	return nil
}
