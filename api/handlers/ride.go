package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"reby/api"
	"reby/domain/ride"
	"reby/domain/user"
	"reby/domain/vehicle"

	"github.com/go-chi/chi/v5"
)

type RideHandlers struct {
	Start  http.Handler
	Finish http.Handler
}

func NewRideHandlers(starter ride.Starter, finisher ride.Finisher) RideHandlers {
	return RideHandlers{
		Start:  Start(starter),
		Finish: Finish(finisher),
	}
}

func AddRideEndpoints(mx *chi.Mux, rh RideHandlers) {
	mx.Method(http.MethodPost, "/rides", rh.Start)
	mx.Method(http.MethodPost, "/rides/{rideID}/finish", rh.Finish)
}

func Start(starter ride.Starter) http.Handler {
	handleError := func(w http.ResponseWriter, err error) {
		switch {
		case errors.Is(err, user.ErrNotFound) || errors.Is(err, vehicle.ErrNotFound):
			api.RespondError(w, api.Error{
				Err:        err,
				HTTPStatus: http.StatusNotFound,
				Reason:     api.InvalidParameter,
			})
		case errors.Is(err, ride.ErrUserIsRiding) || errors.Is(err, ride.ErrVehicleIsRiding):
			api.RespondError(w, api.Error{
				Err:        err,
				HTTPStatus: http.StatusLocked,
				Reason:     api.Locked,
			})
		default:
			api.RespondError(w, api.Error{
				Err:        err,
				HTTPStatus: http.StatusInternalServerError,
				Reason:     api.Internal,
			})
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			UserID    string `json:"user_id"`
			VehicleID string `json:"vehicle_id"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.RespondError(w, api.Error{
				Err:        err,
				HTTPStatus: http.StatusBadRequest,
				Reason:     api.InvalidJSON,
			})
			return
		}

		startedRide, err := starter.Start(r.Context(), ride.StartParams{
			UserID:    req.UserID,
			VehicleID: req.VehicleID,
		})
		if err != nil {
			handleError(w, err)
			return
		}

		api.RespondOK(w, startedRide)
	})
}

func Finish(finisher ride.Finisher) http.Handler {
	handleError := func(w http.ResponseWriter, err error) {
		switch {
		case errors.Is(err, ride.ErrNotFound):
			api.RespondError(w, api.Error{
				Err:        err,
				HTTPStatus: http.StatusNotFound,
				Reason:     api.InvalidParameter,
			})
		case errors.Is(err, ride.ErrAlreadyFinished):
			api.RespondError(w, api.Error{
				Err:        err,
				HTTPStatus: http.StatusConflict,
				Reason:     api.Conflict,
			})
		default:
			api.RespondError(w, api.Error{
				Err:        err,
				HTTPStatus: http.StatusInternalServerError,
				Reason:     api.Internal,
			})
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rideID, err := api.GetStringURLParam(r, "rideID")
		if err != nil {
			api.RespondError(w, api.Error{
				Err:        err,
				HTTPStatus: http.StatusBadRequest,
				Reason:     api.InvalidParameter,
			})
			return
		}

		finishedRide, err := finisher.Finish(r.Context(), rideID)
		if err != nil {
			handleError(w, err)
			return
		}

		api.RespondOK(w, finishedRide)
	})
}
