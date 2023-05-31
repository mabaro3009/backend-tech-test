package ride

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Repo interface {
	GetByID(ctx context.Context, id string) (*Ride, error)
	Create(ctx context.Context, ride *Ride) error
	IsUserRiding(ctx context.Context, userID string) (bool, error)
	IsVehicleRiding(ctx context.Context, vehicleID string) (bool, error)
	Update(ctx context.Context, ride *Ride) (*Ride, error)
}

type RepoMock struct {
	mock.Mock
}

func NewRepoMock() *RepoMock {
	return new(RepoMock)
}

func (m *RepoMock) GetByID(_ context.Context, id string) (*Ride, error) {
	args := m.Mock.Called(id)
	return args.Get(0).(*Ride), args.Error(1)
}

func (m *RepoMock) Create(_ context.Context, ride *Ride) error {
	args := m.Mock.Called(ride)
	return args.Error(0)
}

func (m *RepoMock) IsUserRiding(_ context.Context, userID string) (bool, error) {
	args := m.Mock.Called(userID)
	return args.Get(0).(bool), args.Error(1)
}

func (m *RepoMock) IsVehicleRiding(_ context.Context, vehicleID string) (bool, error) {
	args := m.Mock.Called(vehicleID)
	return args.Get(0).(bool), args.Error(1)
}

func (m *RepoMock) Update(_ context.Context, ride *Ride) (*Ride, error) {
	args := m.Mock.Called(ride)
	return args.Get(0).(*Ride), args.Error(1)
}
