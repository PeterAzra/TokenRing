package token_sender_mocks

import (
	"errors"
	"tokenRing/pkg/node"
)

type SuccessfulTokenSenderMock struct {
}

func NewSuccessfulTokenSenderMock() *SuccessfulTokenSenderMock {
	return &SuccessfulTokenSenderMock{}
}

func (svc *SuccessfulTokenSenderMock) SendToken(from *node.Node, to *node.Node) error {
	return nil
}

type ErrorReturningTokenSenderMock struct {
}

func NewErrorReturningTokenSenderMock() *ErrorReturningTokenSenderMock {
	return &ErrorReturningTokenSenderMock{}
}

func (svc *ErrorReturningTokenSenderMock) SendToken(from *node.Node, to *node.Node) error {
	return errors.New("test error")
}
