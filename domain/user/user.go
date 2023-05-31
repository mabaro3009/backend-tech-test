package user

import "errors"

var (
	ErrNotFound      = errors.New("ERR_USER_NOT_FOUND")
	ErrAlreadyExists = errors.New("ERR_USER_ALREADY_EXISTS")
)

type User struct {
	ID string `json:"id"`
}
