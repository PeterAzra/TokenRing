package link

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"tokenRing/pkg/logging"
	"tokenRing/pkg/models"
	"tokenRing/pkg/node"
	node_http "tokenRing/pkg/node-http"
)

type Linker interface {
	// TODO these could be helper methods and only need LinkNode on the interface
	ConnectLeftAdjacentNode(connectingNode *node.Node, leftAdjacentNodeUrl *url.URL) (bool, error)
	ConnectRightAdjacentNode(connectingNode *node.Node, rightAdjacentNodeUrl *url.URL) (bool, error)
	LinkNode(url *url.URL, request *models.LinkRequest) (bool, error)
}

type LinkService struct {
	client node_http.HttpSender
}

func NewLinkService(client node_http.HttpSender) *LinkService {
	return &LinkService{
		client: client,
	}
}

// Setup link to left adjacent node
func (svc *LinkService) ConnectLeftAdjacentNode(connectingNode *node.Node, leftAdjacentNodeUrl *url.URL) (bool, error) {
	request := models.NewLinkRequest(connectingNode.Url.String())
	// Tell the adjacent node that connecting node is it's right link
	leftAdjacentNodeUrl = leftAdjacentNodeUrl.JoinPath("right-link")
	ok, err := svc.LinkNode(leftAdjacentNodeUrl, request)
	if !ok || err != nil {
		logging.Warning("An error occurred on right link request")
		return false, err
	}

	return true, nil
}

// Setup link to right adjacent node
func (svc *LinkService) ConnectRightAdjacentNode(connectingNode *node.Node, rightAdjacentNodeUrl *url.URL) (bool, error) {
	request := models.NewLinkRequest(connectingNode.Url.String())
	// Tell the adjacent node that connecting node is it's left link
	rightAdjacentNodeUrl = rightAdjacentNodeUrl.JoinPath("left-link")
	ok, err := svc.LinkNode(rightAdjacentNodeUrl, request)
	if !ok || err != nil {
		logging.Warning("An error occurred on left link request")
		return false, err
	}

	return true, nil
}

func (svc *LinkService) LinkNode(url *url.URL, request *models.LinkRequest) (bool, error) {
	logging.Information("Sending link request %v %v", url, request)

	requestData, err := json.Marshal(request)
	if err != nil {
		logging.Error(err, "An error occurred marshalling request")
		return false, err
	}

	req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewBuffer(requestData))
	if err != nil {
		logging.Error(err, "link request create error")
		return false, err
	}

	resp, err := svc.client.Do(req)

	if err != nil {
		logging.Error(err, "Error on link post")
		return false, err
	}

	if resp.StatusCode == http.StatusBadRequest {
		logging.Warning("Invalid link response")
		return false, nil
	}

	defer resp.Body.Close()
	return true, nil
}
