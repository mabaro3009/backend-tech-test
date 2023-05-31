package id

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type Generator interface {
	Generate() string
}

type uuidGenerator struct{}

func NewUUIDGenerator() Generator {
	return &uuidGenerator{}
}

func (g *uuidGenerator) Generate() string {
	return uuid.NewString()
}

type GeneratorMock struct {
	mock.Mock
}

func NewGeneratorMock() *GeneratorMock {
	return new(GeneratorMock)
}

func (m *GeneratorMock) Generate() string {
	args := m.Mock.Called()
	return args.Get(0).(string)
}
