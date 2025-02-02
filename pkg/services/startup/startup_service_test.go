package startup_service

// import (
// 	"errors"
// 	"net/url"
// 	"testing"
// 	"tokenRing/pkg/models"
// 	"tokenRing/pkg/node"

// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// )

// func Test_StartupBaseNode_ReturnsNil_OnBaseNodeStartup(t *testing.T) {
// 	mock := &MockBaseNodeHttp{}
// 	sut := NewStartupService(mock)
// 	baseNodeUrl, _ := url.Parse("http://localhost:8080")
// 	_, ok := sut.StartUpBaseNode(baseNodeUrl)
// 	assert.Equal(t, ok, true, "StartupBaseNode response is not true for base node")
// }

// func Test_StartupBaseNode_ReturnsNew_ForJoiningNode(t *testing.T) {
// 	mock := &MockJoiningNodeStartup{}
// 	sut := NewStartupService(mock)
// 	baseNodeUrl, _ := url.Parse("http://localhost:8080")
// 	_, ok := sut.StartUpBaseNode(baseNodeUrl)
// 	assert.Equal(t, ok, false, "StartupBaseNode response is true for joining node")
// }

// func Test_JoinNodeRing_LinksNodes_ForJoiningNode(t *testing.T) {
// 	baseNodeUrl, _ := url.Parse("http://localhost:8080")
// 	thisNodeUrl, _ := url.Parse("http://localhost:8081")
// 	leftNodeUrl := "http://localhost:8090"
// 	rightNodeUrl := "http://localhost:8091"

// 	mock := &MockJoiningNodeJoinRequest{
// 		LeftUrl:  leftNodeUrl,
// 		RightUrl: rightNodeUrl,
// 	}
// 	sut := NewStartupService(mock)

// 	baseNode := node.NewNode(baseNodeUrl)
// 	thisNode, err := sut.JoinNodeRing(baseNode, thisNodeUrl)

// 	assert.Truef(t, err == nil, "Expected nil err but found: %w", err)
// 	assert.Equal(t, leftNodeUrl, thisNode.Left.Url.String())
// 	assert.Equal(t, rightNodeUrl, thisNode.Right.Url.String())
// }

// func Test_JoinNodeRing_ReturnsFalse_ForJoiningNode(t *testing.T) {
// 	mock := &MockRejectedJoinRequest{}
// 	sut := NewStartupService(mock)
// 	baseNodeUrl, _ := url.Parse("http://localhost:8080")
// 	thisNodeUrl, _ := url.Parse("http://localhost:8081")
// 	baseNode := node.NewNode(baseNodeUrl)
// 	_, err := sut.JoinNodeRing(baseNode, thisNodeUrl)
// 	assert.NotNil(t, err)
// }

// type MockBaseNodeHttp struct{}

// func (mock *MockBaseNodeHttp) PingNode(url *url.URL) (uuid.UUID, error) {
// 	return uuid.Nil, errors.New("test")
// }
// func (mock *MockBaseNodeHttp) Join(url *url.URL, n *node.Node) (*models.JoinResponse, error) {
// 	return nil, nil
// }
// func (mock *MockBaseNodeHttp) LinkNode(url *url.URL, request *models.LinkRequest) (bool, error) {
// 	return true, nil
// }
// func (mock *MockBaseNodeHttp) SendToken(from *node.Node, to *node.Node) error {
// 	return nil
// }

// type MockJoiningNodeStartup struct{}

// func (mock *MockJoiningNodeStartup) PingNode(url *url.URL) (uuid.UUID, error) {
// 	return uuid.New(), nil
// }
// func (mock *MockJoiningNodeStartup) Join(url *url.URL, n *node.Node) (*models.JoinResponse, error) {
// 	return nil, nil
// }
// func (mock *MockJoiningNodeStartup) LinkNode(url *url.URL, request *models.LinkRequest) (bool, error) {
// 	return true, nil
// }
// func (mock *MockJoiningNodeStartup) SendToken(from *node.Node, to *node.Node) error {
// 	return nil
// }

// type MockRejectedJoinRequest struct{}

// func (mock *MockRejectedJoinRequest) PingNode(url *url.URL) (uuid.UUID, error) {
// 	return uuid.New(), nil
// }
// func (mock *MockRejectedJoinRequest) Join(url *url.URL, n *node.Node) (*models.JoinResponse, error) {
// 	joinResp := models.JoinResponse{
// 		Ok: false,
// 	}
// 	return &joinResp, nil
// }
// func (mock *MockRejectedJoinRequest) LinkNode(url *url.URL, request *models.LinkRequest) (bool, error) {
// 	return true, nil
// }
// func (mock *MockRejectedJoinRequest) SendToken(from *node.Node, to *node.Node) error {
// 	return nil
// }

// type MockJoiningNodeJoinRequest struct {
// 	LeftUrl, RightUrl string
// }

// func (mock *MockJoiningNodeJoinRequest) PingNode(url *url.URL) (uuid.UUID, error) {
// 	return uuid.New(), nil
// }
// func (mock *MockJoiningNodeJoinRequest) Join(url *url.URL, n *node.Node) (*models.JoinResponse, error) {
// 	joinResp := models.JoinResponse{
// 		Ok:    true,
// 		Left:  mock.LeftUrl,
// 		Right: mock.RightUrl,
// 	}
// 	return &joinResp, nil
// }
// func (mock *MockJoiningNodeJoinRequest) LinkNode(url *url.URL, request *models.LinkRequest) (bool, error) {
// 	return true, nil
// }
// func (mock *MockJoiningNodeJoinRequest) SendToken(from *node.Node, to *node.Node) error {
// 	return nil
// }
