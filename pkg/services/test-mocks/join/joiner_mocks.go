package joiner_mock

import (
	"errors"
	"net/url"
	"tokenRing/pkg/models"
	"tokenRing/pkg/node"
)

type SuccessfulJoinMock struct {
	leftUrl, rightUrl string
}

func NewSuccessfulJoinMock(leftUrl string, rightUrl string) *SuccessfulJoinMock {
	return &SuccessfulJoinMock{
		leftUrl:  leftUrl,
		rightUrl: rightUrl,
	}
}

func (s *SuccessfulJoinMock) Join(url *url.URL, n *node.Node) (*models.JoinResponse, error) {
	toRtn := models.NewJoinResponse(s.leftUrl, s.rightUrl)
	return toRtn, nil
}

type ErrorReturningJoinMock struct {
}

func NewErrorReturningJoinMock() *ErrorReturningJoinMock {
	return &ErrorReturningJoinMock{}
}

func (s *ErrorReturningJoinMock) Join(url *url.URL, n *node.Node) (*models.JoinResponse, error) {
	return nil, errors.New("test")
}
