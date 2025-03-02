package disconnect

import (
	"net/url"
	"testing"
	"tokenRing/pkg/node"
	linker_mocks "tokenRing/pkg/services/test-mocks/link"
	token_sender_mocks "tokenRing/pkg/services/test-mocks/token-sender"

	"github.com/stretchr/testify/assert"
)

func Test_Disconnect_ReturnsTrue_OnSuccessfulDisconnect(t *testing.T) {
	linkerMock := linker_mocks.NewSuccessfulLinkerMock()
	tokenSenderMock := token_sender_mocks.NewSuccessfulTokenSenderMock()

	nodeUrl, _ := url.Parse("http://localhost:8080")
	rightUrl, _ := url.Parse("http://localhost:8081")
	leftUrl, _ := url.Parse("http://localhost:8082")
	disconnectingNode := node.NewNode(nodeUrl)
	disconnectingNode.Right = node.NewNode(rightUrl)
	disconnectingNode.Left = node.NewNode(leftUrl)

	sut := NewDisconnectService(tokenSenderMock, linkerMock)
	success, err := sut.Disconnect(disconnectingNode)

	assert.True(t, success)
	assert.Nil(t, err)
}

func Test_Disconnect_ReturnsTrue_OnTokenSendFailure(t *testing.T) {
	linkerMock := linker_mocks.NewSuccessfulLinkerMock()
	tokenSenderMock := token_sender_mocks.NewErrorReturningTokenSenderMock()

	nodeUrl, _ := url.Parse("http://localhost:8080")
	rightUrl, _ := url.Parse("http://localhost:8081")
	leftUrl, _ := url.Parse("http://localhost:8082")
	disconnectingNode := node.NewNode(nodeUrl)
	disconnectingNode.Right = node.NewNode(rightUrl)
	disconnectingNode.Left = node.NewNode(leftUrl)

	sut := NewDisconnectService(tokenSenderMock, linkerMock)
	success, err := sut.Disconnect(disconnectingNode)

	assert.True(t, success)
	assert.Nil(t, err)
}

func Test_Disconnect_ReturnsError_OnLinkNodeFailure(t *testing.T) {
	linkerMock := linker_mocks.NewErrorReturningLinkerMock()
	tokenSenderMock := token_sender_mocks.NewSuccessfulTokenSenderMock()

	nodeUrl, _ := url.Parse("http://localhost:8080")
	rightUrl, _ := url.Parse("http://localhost:8081")
	leftUrl, _ := url.Parse("http://localhost:8082")
	disconnectingNode := node.NewNode(nodeUrl)
	disconnectingNode.Right = node.NewNode(rightUrl)
	disconnectingNode.Left = node.NewNode(leftUrl)

	sut := NewDisconnectService(tokenSenderMock, linkerMock)
	success, err := sut.Disconnect(disconnectingNode)

	assert.False(t, success)
	assert.NotNil(t, err)
}
