package test_utils

import (
	"errors"
	"net/url"
	"tokenRing/pkg/models"
	"tokenRing/pkg/node"

	"github.com/google/uuid"
)

type MockSuccessfulLinkRequest struct{}

func (mock *MockSuccessfulLinkRequest) PingNode(url *url.URL) (uuid.UUID, error) {
	return uuid.New(), nil
}
func (mock *MockSuccessfulLinkRequest) Join(url *url.URL, n *node.Node) (*models.JoinResponse, error) {
	return nil, nil
}
func (mock *MockSuccessfulLinkRequest) LinkNode(url *url.URL, request *models.LinkRequest) (bool, error) {
	return true, nil
}
func (mock *MockSuccessfulLinkRequest) SendToken(from *node.Node, to *node.Node) error {
	return nil
}

type MockUnSuccessfulLinkRequest struct{}

func (mock *MockUnSuccessfulLinkRequest) PingNode(url *url.URL) (uuid.UUID, error) {
	return uuid.New(), nil
}
func (mock *MockUnSuccessfulLinkRequest) Join(url *url.URL, n *node.Node) (*models.JoinResponse, error) {
	return nil, nil
}
func (mock *MockUnSuccessfulLinkRequest) LinkNode(url *url.URL, request *models.LinkRequest) (bool, error) {
	return false, errors.New("Fail")
}
func (mock *MockUnSuccessfulLinkRequest) SendToken(from *node.Node, to *node.Node) error {
	return nil
}
