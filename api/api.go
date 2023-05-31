package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Reason string

const (
	InvalidJSON      Reason = "INVALID_JSON"
	InvalidParameter Reason = "INVALID_PARAMETER"
	Internal         Reason = "INTERNAL"
	Locked           Reason = "LOCKED"
	Conflict         Reason = "CONFLICT"
)

var (
	ErrInvalidPathParam = errors.New("ERR_INVALID_PATH_PARAM")
)

type Error struct {
	Err        error
	HTTPStatus int
	Reason     Reason
}

type ErrorDetail struct {
	Reason string `json:"reason"`
	Detail string `json:"detail"`
}

func RespondError(w http.ResponseWriter, err Error) {
	var detail string
	if err.Err != nil {
		detail = err.Err.Error()
	}

	w.WriteHeader(err.HTTPStatus)
	_ = json.NewEncoder(w).Encode(ErrorDetail{
		Reason: string(err.Reason),
		Detail: detail,
	})
}

func RespondOK(w http.ResponseWriter, resp interface{}) {
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		RespondError(w, Error{
			Err:        err,
			HTTPStatus: http.StatusInternalServerError,
			Reason:     Internal,
		})
	}
}

func GetStringURLParam(r *http.Request, key string) (string, error) {
	param := chi.URLParam(r, key)
	if param == "" {
		return "", ErrInvalidPathParam
	}

	return param, nil
}
