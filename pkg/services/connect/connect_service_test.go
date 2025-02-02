package connect

// import (
// 	"net/url"
// 	"testing"
// 	"tokenRing/pkg/node"
// 	node_http_mocks "tokenRing/pkg/test-utils"

// 	"github.com/stretchr/testify/assert"
// )

// func Test_ConnectLeftAdjacentNode_ReturnsTrue_ForSuccessfulRequest(t *testing.T) {
// 	nodeUrl, _ := url.Parse("http://localhost:8080")
// 	adjNodeUrl, _ := url.Parse("http://localhost:8081")
// 	node := node.NewNode(nodeUrl)
// 	result, err := ConnectLeftAdjacentNode(node, adjNodeUrl, &node_http_mocks.MockSuccessfulLinkRequest{})
// 	assert.True(t, result, "ConnectLeftAdjacentNode result was not true for successful request")
// 	assert.Nil(t, err, "ConnectLeftAdjacentNode error was not nil for successful request")
// }

// func Test_ConnectRightAdjacentNode_ReturnsTrue_ForSuccessfulRequest(t *testing.T) {
// 	nodeUrl, _ := url.Parse("http://localhost:8080")
// 	adjNodeUrl, _ := url.Parse("http://localhost:8081")
// 	node := node.NewNode(nodeUrl)
// 	result, err := ConnectRightAdjacentNode(node, adjNodeUrl, &node_http_mocks.MockSuccessfulLinkRequest{})
// 	assert.True(t, result, "ConnectRightAdjacentNode result was not true for successful request")
// 	assert.Nil(t, err, "ConnectRightAdjacentNode error was not nil for successful request")
// }

// func Test_ConnectLeftAdjacentNode_ReturnsFalse_ForUnSuccessfulRequest(t *testing.T) {
// 	nodeUrl, _ := url.Parse("http://localhost:8080")
// 	adjNodeUrl, _ := url.Parse("http://localhost:8081")
// 	node := node.NewNode(nodeUrl)
// 	result, err := ConnectLeftAdjacentNode(node, adjNodeUrl, &node_http_mocks.MockUnSuccessfulLinkRequest{})
// 	assert.False(t, result, "ConnectLeftAdjacentNode result was not false for bad request")
// 	assert.NotNil(t, err, "ConnectLeftAdjacentNode error was nil for bad request")
// }

// func Test_ConnectRightAdjacentNode_ReturnsFalse_ForUnSuccessfulRequest(t *testing.T) {
// 	nodeUrl, _ := url.Parse("http://localhost:8080")
// 	adjNodeUrl, _ := url.Parse("http://localhost:8081")
// 	node := node.NewNode(nodeUrl)
// 	result, err := ConnectRightAdjacentNode(node, adjNodeUrl, &node_http_mocks.MockUnSuccessfulLinkRequest{})
// 	assert.False(t, result, "ConnectRightAdjacentNode result was not false for bad request")
// 	assert.NotNil(t, err, "ConnectRightAdjacentNode error was nil for bad request")
// }
