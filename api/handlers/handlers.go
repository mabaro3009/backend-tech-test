package handlers

import (
	"log"

	"reby/infra/pg"

	_ "github.com/lib/pq" // Postgres driver

	"reby/app/config"
	"reby/domain/ride"
	"reby/domain/user"
	"reby/domain/vehicle"
	"reby/infra"
	"reby/infra/mem"
	"reby/pkg/id"
	"reby/pkg/timenow"
)

type repos struct {
	user    user.Repo
	vehicle vehicle.Repo
	ride    ride.Repo
}

type services struct {
	starter  ride.Starter
	finisher ride.Finisher
}

type Handlers struct {
	Ride RideHandlers
}

func initRepos(conf *config.Config) repos {
	switch conf.DBType {
	case infra.Postgres:
		db := pg.InitDB(conf)
		return repos{
			user:    pg.NewUserDB(db),
			vehicle: pg.NewVehicleDB(db),
			ride:    pg.NewRideDB(db),
		}
	case infra.InMemory:
		return repos{
			user:    mem.NewUserDB(),
			vehicle: mem.NewVehicleDB(),
			ride:    mem.NewRideDB(),
		}
	default:
		log.Fatalf("unrecognized %s memory system", conf.DBType)
		return repos{}
	}
}

func initServices(repos repos) services {
	idGenerator := id.NewUUIDGenerator()
	time := timenow.NewRealTime()
	starter := ride.NewStarter(
		repos.user,
		repos.vehicle,
		repos.ride,
		idGenerator,
		time,
	)

	priceCalculator := ride.NewBasePriceCalculator(
		ride.DefaultUnlockValue,
		ride.DefaultMinuteValue,
		time,
	)

	finisher := ride.NewFinisher(
		repos.ride,
		priceCalculator,
		time,
	)

	return services{
		starter:  starter,
		finisher: finisher,
	}
}

func InitHandlers(conf *config.Config) Handlers {
	r := initRepos(conf)
	svc := initServices(r)
	return Handlers{
		Ride: NewRideHandlers(svc.starter, svc.finisher),
	}
}
