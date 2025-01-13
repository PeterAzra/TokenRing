package connect

import (
	"errors"
	"net/url"
	"testing"
	"tokenRing/pkg/models"
	"tokenRing/pkg/node"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_ConnectLeftAdjacentNode_ReturnsTrue_ForSuccessfulRequest(t *testing.T) {
	nodeUrl, _ := url.Parse("http://localhost:8080")
	adjNodeUrl, _ := url.Parse("http://localhost:8081")
	node := node.NewNode(nodeUrl)
	result, err := ConnectLeftAdjacentNode(node, adjNodeUrl, &MockSuccessfulLinkRequest{})
	assert.True(t, result, "ConnectLeftAdjacentNode result was not true for successful request")
	assert.Nil(t, err, "ConnectLeftAdjacentNode error was not nil for successful request")
}

func Test_ConnectRightAdjacentNode_ReturnsTrue_ForSuccessfulRequest(t *testing.T) {
	nodeUrl, _ := url.Parse("http://localhost:8080")
	adjNodeUrl, _ := url.Parse("http://localhost:8081")
	node := node.NewNode(nodeUrl)
	result, err := ConnectRightAdjacentNode(node, adjNodeUrl, &MockSuccessfulLinkRequest{})
	assert.True(t, result, "ConnectRightAdjacentNode result was not true for successful request")
	assert.Nil(t, err, "ConnectRightAdjacentNode error was not nil for successful request")
}

func Test_ConnectLeftAdjacentNode_ReturnsFalse_ForUnSuccessfulRequest(t *testing.T) {
	nodeUrl, _ := url.Parse("http://localhost:8080")
	adjNodeUrl, _ := url.Parse("http://localhost:8081")
	node := node.NewNode(nodeUrl)
	result, err := ConnectLeftAdjacentNode(node, adjNodeUrl, &MockUnSuccessfulLinkRequest{})
	assert.False(t, result, "ConnectLeftAdjacentNode result was not false for bad request")
	assert.NotNil(t, err, "ConnectLeftAdjacentNode error was nil for bad request")
}

func Test_ConnectRightAdjacentNode_ReturnsFalse_ForUnSuccessfulRequest(t *testing.T) {
	nodeUrl, _ := url.Parse("http://localhost:8080")
	adjNodeUrl, _ := url.Parse("http://localhost:8081")
	node := node.NewNode(nodeUrl)
	result, err := ConnectRightAdjacentNode(node, adjNodeUrl, &MockUnSuccessfulLinkRequest{})
	assert.False(t, result, "ConnectRightAdjacentNode result was not false for bad request")
	assert.NotNil(t, err, "ConnectRightAdjacentNode error was nil for bad request")
}

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
