package disconnect

import (
	"tokenRing/pkg/logging"
	"tokenRing/pkg/models"
	"tokenRing/pkg/node"
	node_http "tokenRing/pkg/node-http"
)

func DisconnectNode(disconnectingNode *node.Node, client node_http.NodeClient) (bool, error) {
	// Send requests to left and right adj nodes to use new links
	rightLinkRequest := models.NewLinkRequest(disconnectingNode.Right.Url.String())
	leftAdjUrl := disconnectingNode.Left.Url.JoinPath("right-link")
	leftDoneCh := make(chan error, 1)

	leftLinkRequest := models.NewLinkRequest(disconnectingNode.Left.Url.String())
	rightAdjUrl := disconnectingNode.Right.Url.JoinPath("left-link")
	rightDoneCh := make(chan error, 1)

	// Send requests together to try and avoid issues while swapping links
	go func() {
		_, err := client.LinkNode(leftAdjUrl, rightLinkRequest)
		leftDoneCh <- err
	}()
	go func() {
		_, err := client.LinkNode(rightAdjUrl, leftLinkRequest)
		rightDoneCh <- err
	}()

	leftErr := <-leftDoneCh
	rightErr := <-rightDoneCh

	// TODO better checking of errors
	if leftErr != nil {
		logging.Error(leftErr, "An error occurred disconnecting from node ring")
		return false, leftErr
	}
	if rightErr != nil {
		logging.Error(rightErr, "An error occurred disconnecting from node ring")
		return false, rightErr
	}

	return true, nil
}
