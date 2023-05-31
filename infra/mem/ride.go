package mem

import (
	"context"
	"time"

	"reby/domain/money"
	"reby/domain/ride"
)

type dbRide struct {
	id         string
	vehicleID  string
	userID     string
	startedAt  time.Time
	finishedAt *time.Time
	price      *money.Money
}

func (r *dbRide) toDomain() *ride.Ride {
	return &ride.Ride{
		ID:         r.id,
		VehicleID:  r.vehicleID,
		UserID:     r.userID,
		StartedAt:  r.startedAt,
		FinishedAt: r.finishedAt,
		Price:      r.price,
	}
}

func rideToDB(r *ride.Ride) *dbRide {
	return &dbRide{
		id:         r.ID,
		vehicleID:  r.VehicleID,
		userID:     r.UserID,
		startedAt:  r.StartedAt,
		finishedAt: r.FinishedAt,
		price:      r.Price,
	}
}

type rideDB struct {
	rides map[string]*dbRide
}

func NewRideDB() ride.Repo {
	return &rideDB{rides: make(map[string]*dbRide)}
}

func (m *rideDB) GetByID(_ context.Context, id string) (*ride.Ride, error) {
	r, ok := m.rides[id]
	if !ok {
		return nil, ride.ErrNotFound
	}

	return r.toDomain(), nil
}

// Update updates some predefined fields of ride.
func (m *rideDB) Update(_ context.Context, r *ride.Ride) (*ride.Ride, error) {
	oldRide, ok := m.rides[r.ID]
	if !ok {
		return nil, ride.ErrNotFound
	}

	// We only want the possibility of updating some fields
	oldRide.finishedAt = r.FinishedAt
	oldRide.price = r.Price
	m.rides[r.ID] = oldRide

	return oldRide.toDomain(), nil
}

func (m *rideDB) Create(_ context.Context, r *ride.Ride) error {
	_, ok := m.rides[r.ID]
	if ok {
		return ride.ErrAlreadyExists
	}

	m.rides[r.ID] = rideToDB(r)

	return nil
}

func (m *rideDB) IsUserRiding(_ context.Context, userID string) (bool, error) {
	for _, r := range m.rides {
		if r.userID == userID && r.finishedAt == nil {
			return true, nil
		}
	}

	return false, nil
}

func (m *rideDB) IsVehicleRiding(_ context.Context, vehicleID string) (bool, error) {
	for _, r := range m.rides {
		if r.vehicleID == vehicleID && r.finishedAt == nil {
			return true, nil
		}
	}

	return false, nil
}
