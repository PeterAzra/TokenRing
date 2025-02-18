package startup_service

import (
	"net/url"
	"testing"
	"tokenRing/pkg/node"
	joiner_mock "tokenRing/pkg/services/test-mocks/join"
	linker_mocks "tokenRing/pkg/services/test-mocks/link"
	pinger_mocks "tokenRing/pkg/services/test-mocks/ping"
	token_sender_mocks "tokenRing/pkg/services/test-mocks/token-sender"

	"github.com/stretchr/testify/assert"
)

func Test_StartupBaseNode_ReturnsTrue_WhenBaseNodeNotExist(t *testing.T) {
	pingSvc := pinger_mocks.NewErrorReturningPingMock()
	tokenSvc := token_sender_mocks.NewSuccessfulTokenSenderMock()
	joinSvc := joiner_mock.NewSuccessfulJoinMock("http://localhost:8081", "http://localhost:8082")
	linkSvc := linker_mocks.NewSuccessfulLinkerMock()

	sut := NewStartupService(pingSvc, joinSvc, linkSvc, tokenSvc)

	baseNodeUrl, _ := url.Parse("http://localhost:8080")
	baseNode := node.NewNode(baseNodeUrl)

	baseNodeTest, err := sut.JoinNodeRing(baseNode, baseNodeUrl)
	assert.Nil(t, err)
	assert.NotNil(t, baseNodeTest)
}
