package ride

import (
	"context"
	"errors"

	"reby/pkg/timenow"

	"github.com/stretchr/testify/mock"
)

type Finisher interface {
	Finish(ctx context.Context, id string) (*Ride, error)
}

var (
	ErrAlreadyFinished = errors.New("ERR_ALREADY_FINISHED")
)

type finisher struct {
	rideRepo        Repo
	priceCalculator PriceCalculator
	time            timenow.TimeNow
}

func NewFinisher(rideRepo Repo, priceCalculator PriceCalculator, time timenow.TimeNow) Finisher {
	return &finisher{rideRepo: rideRepo, priceCalculator: priceCalculator, time: time}
}

func (f *finisher) Finish(ctx context.Context, id string) (*Ride, error) {
	r, err := f.rideRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if r.FinishedAt != nil {
		return nil, ErrAlreadyFinished
	}

	price, err := f.priceCalculator.Calculate(*r)
	if err != nil {
		return nil, err
	}

	now := f.time.Now()
	r.FinishedAt = &now
	r.Price = &price

	return f.rideRepo.Update(ctx, r)
}

type FinisherMock struct {
	mock.Mock
}

func NewFinisherMock() *FinisherMock {
	return new(FinisherMock)
}

func (m *FinisherMock) Finish(_ context.Context, id string) (*Ride, error) {
	args := m.Mock.Called(id)
	return args.Get(0).(*Ride), args.Error(1)
}
