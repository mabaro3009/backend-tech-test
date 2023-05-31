package vehicle

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Repo interface {
	GetByID(ctx context.Context, id string) (*Vehicle, error)
	Create(ctx context.Context, v *Vehicle) error
}

type RepoMock struct {
	mock.Mock
}

func NewRepoMock() *RepoMock {
	return new(RepoMock)
}

func (m *RepoMock) GetByID(_ context.Context, id string) (*Vehicle, error) {
	args := m.Mock.Called(id)
	return args.Get(0).(*Vehicle), args.Error(1)
}

func (m *RepoMock) Create(_ context.Context, v *Vehicle) error {
	args := m.Mock.Called(v)
	return args.Error(0)
}
