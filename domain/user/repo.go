package user

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Repo interface {
	GetByID(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, u *User) error
}

type RepoMock struct {
	mock.Mock
}

func NewRepoMock() *RepoMock {
	return new(RepoMock)
}

func (m *RepoMock) GetByID(_ context.Context, id string) (*User, error) {
	args := m.Mock.Called(id)
	return args.Get(0).(*User), args.Error(1)
}

func (m *RepoMock) Create(_ context.Context, u *User) error {
	args := m.Mock.Called(u)
	return args.Error(0)
}
