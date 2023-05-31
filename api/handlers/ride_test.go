package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reby/api"
	"reby/api/handlers"
	"reby/domain/money"
	"reby/domain/ride"
	"reby/domain/user"
	"reby/domain/vehicle"
)

func TestRideStart(t *testing.T) {
	var starterMock *ride.StarterMock
	var hd handlers.RideHandlers
	ctx := context.Background()

	setup := func() {
		starterMock = ride.NewStarterMock()
		hd = handlers.NewRideHandlers(starterMock, nil)
	}

	doReq := func() *httptest.ResponseRecorder {
		body := struct {
			UserID    string `json:"user_id"`
			VehicleID string `json:"vehicle_id"`
		}{
			UserID:    "1",
			VehicleID: "1",
		}
		var buf bytes.Buffer
		require.NoError(t, json.NewEncoder(&buf).Encode(body))
		req, err := http.NewRequest(http.MethodPost, "/rides", &buf)
		require.NoError(t, err)

		resp := httptest.NewRecorder()
		hd.Start.ServeHTTP(resp, req.WithContext(ctx))

		return resp
	}

	testCases := []struct {
		description    string
		starterErr     error
		expectedCode   int
		expectedReason string
		expectedDetail string
	}{
		{
			description:    "user not found",
			starterErr:     user.ErrNotFound,
			expectedCode:   http.StatusNotFound,
			expectedReason: string(api.InvalidParameter),
			expectedDetail: "ERR_USER_NOT_FOUND",
		},
		{
			description:    "vehicle not found",
			starterErr:     vehicle.ErrNotFound,
			expectedCode:   http.StatusNotFound,
			expectedReason: string(api.InvalidParameter),
			expectedDetail: "ERR_VEHICLE_NOT_FOUND",
		},
		{
			description:    "user riding",
			starterErr:     ride.ErrUserIsRiding,
			expectedCode:   http.StatusLocked,
			expectedReason: string(api.Locked),
			expectedDetail: "ERR_USER_RIDING",
		},
		{
			description:    "vehicle riding",
			starterErr:     ride.ErrVehicleIsRiding,
			expectedCode:   http.StatusLocked,
			expectedReason: string(api.Locked),
			expectedDetail: "ERR_VEHICLE_RIDING",
		},
		{
			description:    "internal error",
			starterErr:     errors.New("ERR_RANDOM_ERROR"),
			expectedCode:   http.StatusInternalServerError,
			expectedReason: string(api.Internal),
			expectedDetail: "ERR_RANDOM_ERROR",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			setup()
			starterMock.On("Start", ride.StartParams{
				UserID:    "1",
				VehicleID: "1",
			}).Return(&ride.Ride{}, tc.starterErr)

			resp := doReq()
			assert.Equal(t, tc.expectedCode, resp.Code)

			var errorDetail api.ErrorDetail
			require.NoError(t, json.NewDecoder(resp.Body).Decode(&errorDetail))
			assert.Equal(t, tc.expectedDetail, errorDetail.Detail)
			assert.Equal(t, tc.expectedReason, errorDetail.Reason)
		})
	}

	t.Run("ok", func(t *testing.T) {
		setup()
		r := &ride.Ride{
			ID:         "1",
			VehicleID:  "1",
			UserID:     "1",
			StartedAt:  time.Now(),
			FinishedAt: nil,
			Price:      nil,
		}

		starterMock.On("Start", ride.StartParams{
			UserID:    "1",
			VehicleID: "1",
		}).Return(r, nil)

		resp := doReq()
		assert.Equal(t, http.StatusOK, resp.Code)

		var respRide *ride.Ride
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&respRide))
		assert.NotEmpty(t, respRide)
	})
}

func TestRideFinish(t *testing.T) {
	var finisherMock *ride.FinisherMock
	var hd handlers.RideHandlers
	rideID := "r_1"

	setup := func() {
		finisherMock = ride.NewFinisherMock()
		hd = handlers.NewRideHandlers(nil, finisherMock)
	}

	doReq := func(rideID string) *httptest.ResponseRecorder {
		path := fmt.Sprintf("/rides/%s/finish", rideID)
		req, err := http.NewRequest(http.MethodPost, path, nil)
		require.NoError(t, err)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("rideID", rideID)

		resp := httptest.NewRecorder()
		hd.Finish.ServeHTTP(resp, req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx)))

		return resp
	}

	testCases := []struct {
		description    string
		rideID         string
		finisherErr    error
		expectedCode   int
		expectedReason string
		expectedDetail string
	}{
		{
			description:    "ride path param invalid",
			rideID:         "",
			finisherErr:    nil,
			expectedCode:   http.StatusBadRequest,
			expectedReason: string(api.InvalidParameter),
			expectedDetail: "ERR_INVALID_PATH_PARAM",
		},
		{
			description:    "ride not found",
			rideID:         rideID,
			finisherErr:    ride.ErrNotFound,
			expectedCode:   http.StatusNotFound,
			expectedReason: string(api.InvalidParameter),
			expectedDetail: "ERR_RIDE_NOT_FOUND",
		},
		{
			description:    "ride finished",
			rideID:         rideID,
			finisherErr:    ride.ErrAlreadyFinished,
			expectedCode:   http.StatusConflict,
			expectedReason: string(api.Conflict),
			expectedDetail: "ERR_ALREADY_FINISHED",
		},
		{
			description:    "internal",
			rideID:         rideID,
			finisherErr:    errors.New("ERR_RANDOM"),
			expectedCode:   http.StatusInternalServerError,
			expectedReason: string(api.Internal),
			expectedDetail: "ERR_RANDOM",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			setup()
			finisherMock.On("Finish", tc.rideID).Return(&ride.Ride{}, tc.finisherErr)

			resp := doReq(tc.rideID)
			assert.Equal(t, tc.expectedCode, resp.Code)

			var errorDetail api.ErrorDetail
			require.NoError(t, json.NewDecoder(resp.Body).Decode(&errorDetail))
			assert.Equal(t, tc.expectedDetail, errorDetail.Detail)
			assert.Equal(t, tc.expectedReason, errorDetail.Reason)
		})
	}

	t.Run("ok", func(t *testing.T) {
		setup()
		now := time.Now()
		finishedRide := &ride.Ride{
			ID:         rideID,
			VehicleID:  "1",
			UserID:     "1",
			StartedAt:  now.Add(-5 * time.Minute),
			FinishedAt: &now,
			Price: &money.Money{
				Value:    100,
				Currency: "EUR",
			},
		}
		finisherMock.On("Finish", rideID).Return(finishedRide, nil)

		resp := doReq(rideID)
		assert.Equal(t, http.StatusOK, resp.Code)

		var respRide *ride.Ride
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&respRide))
		assert.NotEmpty(t, respRide)
	})
}
