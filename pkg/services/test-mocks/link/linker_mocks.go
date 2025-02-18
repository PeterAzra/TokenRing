package linker_mocks

import (
	"errors"
	"net/url"
	"tokenRing/pkg/models"
	"tokenRing/pkg/node"
)

type SuccessfulLinkerMock struct {
}

func NewSuccessfulLinkerMock() *SuccessfulLinkerMock {
	return &SuccessfulLinkerMock{}
}

func (svc *SuccessfulLinkerMock) ConnectLeftAdjacentNode(connectingNode *node.Node, leftAdjacentNodeUrl *url.URL) (bool, error) {
	return true, nil
}
func (svc *SuccessfulLinkerMock) ConnectRightAdjacentNode(connectingNode *node.Node, rightAdjacentNodeUrl *url.URL) (bool, error) {
	return true, nil
}
func (svc *SuccessfulLinkerMock) LinkNode(url *url.URL, request *models.LinkRequest) (bool, error) {
	return true, nil
}

type ErrorReturningLinkerMock struct {
}

func NewErrorReturningLinkerMock() *ErrorReturningLinkerMock {
	return &ErrorReturningLinkerMock{}
}

func (svc *ErrorReturningLinkerMock) ConnectLeftAdjacentNode(connectingNode *node.Node, leftAdjacentNodeUrl *url.URL) (bool, error) {
	return false, errors.New("test error")
}
func (svc *ErrorReturningLinkerMock) ConnectRightAdjacentNode(connectingNode *node.Node, rightAdjacentNodeUrl *url.URL) (bool, error) {
	return false, errors.New("test error")
}
func (svc *ErrorReturningLinkerMock) LinkNode(url *url.URL, request *models.LinkRequest) (bool, error) {
	return false, errors.New("test error")
}
