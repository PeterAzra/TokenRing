package pinger_mocks

import (
	"errors"
	"net/url"

	"github.com/google/uuid"
)

type SuccessfulPingMock struct {
}

func NewSuccessfulPingMock() *SuccessfulPingMock {
	return &SuccessfulPingMock{}
}

func (svc *SuccessfulPingMock) Ping(url *url.URL) (uuid.UUID, error) {
	toRtn, _ := uuid.NewUUID()
	return toRtn, nil
}

type ErrorReturningPingMock struct {
}

func NewErrorReturningPingMock() *ErrorReturningPingMock {
	return &ErrorReturningPingMock{}
}

func (svc *ErrorReturningPingMock) Ping(url *url.URL) (uuid.UUID, error) {
	return uuid.Nil, errors.New("test")
}
