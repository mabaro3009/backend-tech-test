package ride

import (
	"errors"
	"time"

	"reby/domain/money"
)

var (
	ErrAlreadyExists = errors.New("ERR_RIDE_ALREADY_EXISTS")
	ErrNotFound      = errors.New("ERR_RIDE_NOT_FOUND")
)

type Ride struct {
	ID         string       `json:"id"`
	VehicleID  string       `json:"vehicle_id"`
	UserID     string       `json:"user_id"`
	StartedAt  time.Time    `json:"started_at"`
	FinishedAt *time.Time   `json:"finished_at"`
	Price      *money.Money `json:"price"`
}
