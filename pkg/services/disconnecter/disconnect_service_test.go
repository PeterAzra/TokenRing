package disconnecter

import (
	"net/url"
	"testing"
	"tokenRing/pkg/node"
	linker_mocks "tokenRing/pkg/services/test-mocks/linker"
	token_sender_mocks "tokenRing/pkg/services/test-mocks/token-sender"

	"github.com/stretchr/testify/assert"
)

func Test_Disconnect_ReturnsTrue_OnSuccessfulDisconnect(t *testing.T) {
	linkerMock := linker_mocks.NewSuccessfulLinkerMock()
	tokenSenderMock := token_sender_mocks.NewSuccessfulTokenSenderMock()

	nodeUrl, _ := url.Parse("http://localhost:8080")
	disconnectingNode := node.NewNode(nodeUrl)

	sut := NewDisconnectService(tokenSenderMock, linkerMock)
	success, err := sut.Disconnect(disconnectingNode)

	assert.True(t, success)
	assert.NotNil(t, err)
}

func Test_Disconnect_ReturnsTrue_OnTokenSendFailure(t *testing.T) {
	linkerMock := linker_mocks.NewSuccessfulLinkerMock()
	tokenSenderMock := token_sender_mocks.NewErrorReturningTokenSenderMock()

	nodeUrl, _ := url.Parse("http://localhost:8080")
	disconnectingNode := node.NewNode(nodeUrl)

	sut := NewDisconnectService(tokenSenderMock, linkerMock)
	success, err := sut.Disconnect(disconnectingNode)

	assert.True(t, success)
	assert.NotNil(t, err)
}

func Test_Disconnect_ReturnsError_OnLinkNodeFailure(t *testing.T) {
	linkerMock := linker_mocks.NewErrorReturningLinkerMock()
	tokenSenderMock := token_sender_mocks.NewSuccessfulTokenSenderMock()

	nodeUrl, _ := url.Parse("http://localhost:8080")
	disconnectingNode := node.NewNode(nodeUrl)

	sut := NewDisconnectService(tokenSenderMock, linkerMock)
	success, err := sut.Disconnect(disconnectingNode)

	assert.False(t, success)
	assert.Nil(t, err)
}
