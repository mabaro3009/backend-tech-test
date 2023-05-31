package ride

import (
	"context"
	"errors"

	"reby/domain/user"
	"reby/domain/vehicle"
	"reby/pkg/id"
	"reby/pkg/timenow"

	"github.com/stretchr/testify/mock"
)

type Starter interface {
	Start(ctx context.Context, params StartParams) (*Ride, error)
}

var (
	ErrUserIsRiding    = errors.New("ERR_USER_RIDING")
	ErrVehicleIsRiding = errors.New("ERR_VEHICLE_RIDING")
)

type StartParams struct {
	UserID    string
	VehicleID string
}

type starter struct {
	userRepo    user.Repo
	vehicleRepo vehicle.Repo
	rideRepo    Repo
	idGenerator id.Generator
	time        timenow.TimeNow
}

func NewStarter(userRepo user.Repo, vehicleRepo vehicle.Repo, rideRepo Repo, idGenerator id.Generator, time timenow.TimeNow) Starter {
	return &starter{userRepo: userRepo, vehicleRepo: vehicleRepo, rideRepo: rideRepo, idGenerator: idGenerator, time: time}
}

func (s *starter) Start(ctx context.Context, params StartParams) (*Ride, error) {
	u, err := s.userRepo.GetByID(ctx, params.UserID)
	if err != nil {
		return nil, err
	}

	v, err := s.vehicleRepo.GetByID(ctx, params.VehicleID)
	if err != nil {
		return nil, err
	}

	if err = s.startChecks(ctx, u.ID, v.ID); err != nil {
		return nil, err
	}

	r := &Ride{
		ID:         s.idGenerator.Generate(),
		VehicleID:  v.ID,
		UserID:     u.ID,
		StartedAt:  s.time.Now(),
		FinishedAt: nil,
		Price:      nil,
	}
	if err = s.rideRepo.Create(ctx, r); err != nil {
		return nil, err
	}

	return r, nil
}

func (s *starter) startChecks(ctx context.Context, userID string, vehicleID string) error {
	isUserRiding, err := s.rideRepo.IsUserRiding(ctx, userID)
	if err != nil {
		return err
	}
	if isUserRiding {
		return ErrUserIsRiding
	}

	isVehicleRiding, err := s.rideRepo.IsVehicleRiding(ctx, vehicleID)
	if err != nil {
		return err
	}
	if isVehicleRiding {
		return ErrVehicleIsRiding
	}

	return nil
}

type StarterMock struct {
	mock.Mock
}

func NewStarterMock() *StarterMock {
	return new(StarterMock)
}

func (m *StarterMock) Start(_ context.Context, params StartParams) (*Ride, error) {
	args := m.Mock.Called(params)
	return args.Get(0).(*Ride), args.Error(1)
}
