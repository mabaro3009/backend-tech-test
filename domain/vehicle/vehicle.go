package vehicle

import "errors"

var (
	ErrNotFound      = errors.New("ERR_VEHICLE_NOT_FOUND")
	ErrAlreadyExists = errors.New("ERR_VEHICLE_ALREADY_EXISTS")
)

type Vehicle struct {
	ID string `json:"id"`
}
