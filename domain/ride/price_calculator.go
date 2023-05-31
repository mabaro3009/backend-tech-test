package ride

import (
	"errors"
	"math"

	"reby/domain/money"
	"reby/pkg/timenow"

	"github.com/stretchr/testify/mock"
)

type PriceCalculator interface {
	Calculate(ride Ride) (money.Money, error)
}

const (
	DefaultUnlockValue   = 100
	DefaultMinuteValue   = 18
	defaultPriceCurrency = "EUR" // To simplify things for the task, all rides are in EUR
)

var (
	ErrInvalidRideMinutes = errors.New("ERR_INVALID_RIDE_MINUTES")
)

type basePriceCalculator struct {
	unlockValue int
	minuteValue int
	time        timenow.TimeNow
}

func NewBasePriceCalculator(unlockValue int, minuteValue int, time timenow.TimeNow) PriceCalculator {
	return &basePriceCalculator{
		unlockValue: unlockValue,
		minuteValue: minuteValue,
		time:        time,
	}
}

func (c *basePriceCalculator) Calculate(ride Ride) (money.Money, error) {
	minutes, err := c.getMinutesFromRide(ride)
	if err != nil {
		return money.Money{}, err
	}
	value := c.unlockValue + minutes*c.minuteValue
	return money.NewMoney(value, defaultPriceCurrency), nil
}

func (c *basePriceCalculator) getMinutesFromRide(ride Ride) (int, error) {
	finishedAt := c.time.Now()
	if ride.FinishedAt != nil {
		finishedAt = *ride.FinishedAt
	}

	minutes := int(math.Ceil(finishedAt.Sub(ride.StartedAt).Minutes()))
	if minutes < 0 {
		return 0, ErrInvalidRideMinutes
	}

	return minutes, nil
}

type PriceCalculatorMock struct {
	mock.Mock
}

func NewPriceCalculatorMock() *PriceCalculatorMock {
	return new(PriceCalculatorMock)
}

func (m *PriceCalculatorMock) Calculate(ride Ride) (money.Money, error) {
	args := m.Mock.Called(ride)
	return args.Get(0).(money.Money), args.Error(1)
}
