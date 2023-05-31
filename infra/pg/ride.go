package pg

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"reby/domain/money"
	"reby/domain/ride"
)

type dbRide struct {
	id            string     `db:"id"`
	vehicleID     string     `db:"vehicle_id"`
	userID        string     `db:"user_id"`
	startedAt     time.Time  `db:"started_at"`
	finishedAt    *time.Time `db:"finished_at"`
	priceValue    *int       `db:"price_value"`
	priceCurrency *string    `db:"price_currency"`
}

func (r *dbRide) toDomain() *ride.Ride {
	var price *money.Money
	if r.priceValue != nil && r.priceCurrency != nil {
		p := money.NewMoney(*r.priceValue, *r.priceCurrency)
		price = &p
	}
	return &ride.Ride{
		ID:         r.id,
		VehicleID:  r.vehicleID,
		UserID:     r.userID,
		StartedAt:  r.startedAt,
		FinishedAt: r.finishedAt,
		Price:      price,
	}
}

func toRideDB(r *ride.Ride) *dbRide {
	rd := &dbRide{
		id:            r.ID,
		vehicleID:     r.VehicleID,
		userID:        r.UserID,
		startedAt:     r.StartedAt,
		finishedAt:    r.FinishedAt,
		priceValue:    nil,
		priceCurrency: nil,
	}
	if r.Price != nil {
		pv := r.Price.Value.Int()
		rd.priceValue = &pv

		pc := r.Price.Currency.String()
		rd.priceCurrency = &pc
	}

	return rd
}

type rideDB struct {
	db *sql.DB
}

func NewRideDB(db *sql.DB) ride.Repo {
	return &rideDB{db: db}
}

func (db *rideDB) GetByID(ctx context.Context, id string) (*ride.Ride, error) {
	q := `SELECT id, vehicle_id, user_id, started_at, finished_at, price_value, price_currency FROM "ride" WHERE id=$1;`

	var r dbRide
	if err := db.db.QueryRowContext(ctx, q, id).Scan(
		&r.id, &r.vehicleID, &r.userID, &r.startedAt, &r.finishedAt, &r.priceValue, &r.priceCurrency,
	); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, ride.ErrNotFound
		}
		return nil, err
	}

	return r.toDomain(), nil
}

func (db *rideDB) Create(ctx context.Context, r *ride.Ride) error {
	rDB := toRideDB(r)
	q := `INSERT INTO "ride" (id, vehicle_id, user_id, started_at) VALUES ($1, $2, $3, $4);`

	if _, err := db.db.ExecContext(ctx, q, rDB.id, rDB.vehicleID, rDB.userID, rDB.startedAt); err != nil {
		return err
	}

	return nil
}

func (db *rideDB) IsUserRiding(ctx context.Context, userID string) (bool, error) {
	q := `SELECT 1 FROM "ride" WHERE user_id=$1 AND finished_at is null;`

	var result int
	if err := db.db.QueryRowContext(ctx, q, userID).Scan(&result); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (db *rideDB) IsVehicleRiding(ctx context.Context, vehicleID string) (bool, error) {
	q := `SELECT 1 FROM "ride" WHERE vehicle_id=$1 AND finished_at is null;`

	var result int
	if err := db.db.QueryRowContext(ctx, q, vehicleID).Scan(&result); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (db *rideDB) Update(ctx context.Context, r *ride.Ride) (*ride.Ride, error) {
	q := `UPDATE "ride" SET finished_at=$1, price_value=$2, price_currency=$3 WHERE ID=$4;`
	rDB := toRideDB(r)

	if _, err := db.db.ExecContext(ctx, q, rDB.finishedAt, rDB.priceValue, rDB.priceCurrency, rDB.id); err != nil {
		return nil, err
	}

	return db.GetByID(ctx, rDB.id)
}
