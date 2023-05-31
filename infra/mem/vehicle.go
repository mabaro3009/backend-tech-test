package mem

import (
	"context"

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
	vehicles map[string]dbVehicle
}

func NewVehicleDB() vehicle.Repo {
	return &vehicleDB{vehicles: map[string]dbVehicle{"1": {"1"}, "2": {"2"}}}
}

func (m *vehicleDB) GetByID(_ context.Context, id string) (*vehicle.Vehicle, error) {
	v, ok := m.vehicles[id]
	if !ok {
		return nil, vehicle.ErrNotFound
	}

	return v.toDomain(), nil
}

func (m *vehicleDB) Create(_ context.Context, v *vehicle.Vehicle) error {
	vDB := toVehicleDB(v)
	if _, ok := m.vehicles[vDB.id]; ok {
		return vehicle.ErrAlreadyExists
	}

	m.vehicles[vDB.id] = vDB
	return nil
}
