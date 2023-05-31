package pg

import (
	"context"
	"database/sql"
	"errors"

	"reby/domain/user"
)

type dbUser struct {
	id string `db:"id"`
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
	db *sql.DB
}

func NewUserDB(db *sql.DB) user.Repo {
	return &userDB{db: db}
}

func (db *userDB) GetByID(ctx context.Context, id string) (*user.User, error) {
	q := `SELECT id FROM "user" WHERE id=$1;`

	var u dbUser
	if err := db.db.QueryRowContext(ctx, q, id).Scan(&u.id); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, user.ErrNotFound
		}
		return nil, err
	}

	return u.toDomain(), nil
}

func (db *userDB) Create(ctx context.Context, u *user.User) error {
	uDB := toUserDB(u)
	q := `INSERT INTO "user" (id) VALUES ($1);`

	if _, err := db.db.ExecContext(ctx, q, uDB.id); err != nil {
		return err
	}

	return nil
}
