package disconnect

import (
	"net/url"
	"testing"
	"tokenRing/pkg/node"
	node_http_mocks "tokenRing/pkg/test-utils"

	"github.com/stretchr/testify/assert"
)

func Test_Disconnect_ReturnsTrue_OnSuccessfulDisconnect(t *testing.T) {
	nodeUrl, _ := url.Parse("http://localhost:8080")
	leftUrl, _ := url.Parse("http://localhost:8081")
	rightUrl, _ := url.Parse("http://localhost:8082")
	disconnectingNode := node.NewNode(nodeUrl)
	disconnectingNode.Left = node.NewNode(leftUrl)
	disconnectingNode.Right = node.NewNode(rightUrl)

	result, err := DisconnectNode(disconnectingNode, &node_http_mocks.MockSuccessfulLinkRequest{})

	assert.True(t, result, "Expected true, but found false on disconnect")
	assert.Nil(t, err, "Expected nil error, but found error on disconnect")
}

func Test_Disconnect_ReturnsError_OnUnsuccessfulDisconnect(t *testing.T) {
	nodeUrl, _ := url.Parse("http://localhost:8080")
	leftUrl, _ := url.Parse("http://localhost:8081")
	rightUrl, _ := url.Parse("http://localhost:8082")
	disconnectingNode := node.NewNode(nodeUrl)
	disconnectingNode.Left = node.NewNode(leftUrl)
	disconnectingNode.Right = node.NewNode(rightUrl)

	result, err := DisconnectNode(disconnectingNode, &node_http_mocks.MockUnSuccessfulLinkRequest{})

	assert.False(t, result, "Expected false, but found true on disconnect")
	assert.NotNil(t, err, "Expected error, but found nil on disconnect")
}
