package pg

import (
	"context"
	"database/sql"
	"errors"

	"reby/domain/vehicle"
)

type dbVehicle struct {
	id string
}

func (v dbVehicle) toDomain() *vehicle.Vehicle {
	return &vehicle.Vehicle{
		ID: v.id,
	}
}

func toVehicleDB(v *vehicle.Vehicle) dbVehicle {
	return dbVehicle{id: v.ID}
}

type vehicleDB struct {
	db *sql.DB
}

func NewVehicleDB(db *sql.DB) vehicle.Repo {
	return &vehicleDB{db: db}
}

func (db *vehicleDB) GetByID(ctx context.Context, id string) (*vehicle.Vehicle, error) {
	q := `SELECT id FROM "vehicle" WHERE id=$1;`

	var v dbVehicle
	if err := db.db.QueryRowContext(ctx, q, id).Scan(&v.id); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, vehicle.ErrNotFound
		}
		return nil, err
	}

	return v.toDomain(), nil
}

func (db *vehicleDB) Create(ctx context.Context, v *vehicle.Vehicle) error {
	vDB := toVehicleDB(v)
	q := `INSERT INTO "vehicle" (id) VALUES ($1);`

	if _, err := db.db.ExecContext(ctx, q, vDB.id); err != nil {
		return err
	}

	return nil
}
