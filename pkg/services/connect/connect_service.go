package connect

import (
	"net/url"
	"tokenRing/pkg/logging"
	"tokenRing/pkg/models"
	"tokenRing/pkg/node"
	node_http "tokenRing/pkg/node-http"
)

// Setup link to left adjacent node
func ConnectLeftAdjacentNode(connectingNode *node.Node, leftAdjacentNodeUrl *url.URL, client node_http.NodeClient) (bool, error) {
	request := models.NewLinkRequest(connectingNode.Url.String())
	// Tell the adjacent node that connecting node is it's right link
	leftAdjacentNodeUrl = leftAdjacentNodeUrl.JoinPath("right-link")
	ok, err := client.LinkNode(leftAdjacentNodeUrl, request)
	if !ok || err != nil {
		logging.Warning("An error occurred on right link request")
		return false, err
	}

	return true, nil
}

// Setup link to right adjacent node
func ConnectRightAdjacentNode(connectingNode *node.Node, rightAdjacentNodeUrl *url.URL, client node_http.NodeClient) (bool, error) {
	request := models.NewLinkRequest(connectingNode.Url.String())
	// Tell the adjacent node that connecting node is it's left link
	rightAdjacentNodeUrl = rightAdjacentNodeUrl.JoinPath("left-link")
	ok, err := client.LinkNode(rightAdjacentNodeUrl, request)
	if !ok || err != nil {
		logging.Warning("An error occurred on left link request")
		return false, err
	}

	return true, nil
}
