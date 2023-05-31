package mem

import (
	"context"

	"reby/domain/user"
)

type dbUser struct {
	id string
}

func (u dbUser) toDomain() *user.User {
	return &user.User{
		ID: u.id,
	}
}

func toUserDB(u *user.User) dbUser {
	return dbUser{id: u.ID}
}

type userDB struct {
	users map[string]dbUser
}

func NewUserDB() user.Repo {
	return &userDB{users: map[string]dbUser{"1": {"1"}, "2": {"2"}}}
}

func (m *userDB) GetByID(_ context.Context, id string) (*user.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, user.ErrNotFound
	}

	return u.toDomain(), nil
}

func (m *userDB) Create(_ context.Context, u *user.User) error {
	uDB := toUserDB(u)
	if _, ok := m.users[uDB.id]; ok {
		return user.ErrAlreadyExists
	}

	m.users[uDB.id] = uDB
	return nil
}
