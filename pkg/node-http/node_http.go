package node_http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"tokenRing/pkg/logging"
	"tokenRing/pkg/models"
	"tokenRing/pkg/node"

	"github.com/google/uuid"
)

type NodeClient interface {
	Join(n *node.Node, baseUrl string) (*models.JoinResponse, error)
	PingNode(url *url.URL) (uuid.UUID, error)
	LinkLeftNode(n *node.Node, adjNodeUrl *url.URL) (bool, error)
	LinkRightNode(n *node.Node, adjNodeUrl *url.URL) (bool, error)
}

type NodeHttpClient struct {
	HttpClient HttpSender
}

func NewNodeHttpClient() *NodeHttpClient {
	var httpClient HttpSender = NewHttpClient()
	return &NodeHttpClient{
		HttpClient: httpClient,
	}
}

func (client *NodeHttpClient) PingNode(url *url.URL) (uuid.UUID, error) {
	endpoint := fmt.Sprintf("%v/%v", url, "ping")
	logging.Information("Pinging node %v", endpoint)

	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		logging.Error(err, "create ping request fail")
		return uuid.Nil, err
	}

	resp, err := client.HttpClient.Do(request)
	if err != nil {
		logging.Warning(err.Error())
		return uuid.Nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Warning(err.Error())
		return uuid.Nil, err
	}

	uuidResp, err := uuid.ParseBytes(body)
	if err != nil {
		logging.Warning(err.Error())
		return uuid.Nil, err
	}

	logging.Information("Ping successful %v\n", uuidResp.String())
	return uuidResp, nil
}

func (client *NodeHttpClient) Join(n *node.Node, baseUrl string) (*models.JoinResponse, error) {
	logging.Information("Contacting node to join ring")
	joinRequest := models.NewJoinRequest(n.Id, n.Url.String())
	resp, err := client.sendJoinRequest(joinRequest, baseUrl)
	if err != nil {
		logging.Warning("An error occurred on join request")
		return nil, err
	}

	logging.Information("Join request successful")

	return resp, nil
}

func (client *NodeHttpClient) sendJoinRequest(request *models.JoinRequest, baseUrl string) (*models.JoinResponse, error) {
	endpoint := fmt.Sprintf("%v/%v", baseUrl, "joinrequest")
	logging.Information("Sending join request %v %v", endpoint, request)

	requestData, err := json.Marshal(request)
	if err != nil {
		logging.Error(err, "An error occurred marshalling join request model")
		return nil, err
	}

	//log.Printf("Join request %v", string(requestData))

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(requestData))
	if err != nil {
		logging.Error(err, "link request creation error")
		return nil, err
	}

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		logging.Error(err, "An error occurred on join request")
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Error(err, "An error occurred reading join response")
		return nil, err
	}

	//log.Println(string(body))

	var joinResp models.JoinResponse
	err = json.Unmarshal(body, &joinResp)
	if err != nil {
		logging.Error(err, "An error occurred on join response unmarshal")
		return nil, err
	}

	logging.Information("Join response %v", joinResp)

	return &joinResp, nil
}

// Link the current node's left link with adjacent node's right link
func (client *NodeHttpClient) LinkLeftNode(n *node.Node, adjNodeUrl *url.URL) (bool, error) {
	request := models.NewLinkRequest(n.Url.String())
	endpoint := fmt.Sprintf("%v/%v", adjNodeUrl, "right-link")

	ok, err := client.sendLinkRequest(endpoint, request)
	if !ok || err != nil {
		logging.Warning("An error occurred on right link request")
		return false, err
	}

	return true, nil
}

// Link the current node's right link with adjacent node's left link
func (client *NodeHttpClient) LinkRightNode(n *node.Node, adjNodeUrl *url.URL) (bool, error) {
	request := models.NewLinkRequest(n.Url.String())
	endpoint := fmt.Sprintf("%v/%v", adjNodeUrl, "left-link")

	ok, err := client.sendLinkRequest(endpoint, request)
	if !ok || err != nil {
		logging.Warning("An error occurred on left link request")
		return false, err
	}

	return true, nil
}

func (client *NodeHttpClient) sendLinkRequest(endpoint string, request *models.LinkRequest) (bool, error) {
	logging.Information("Sending link request %v %v", endpoint, request)

	requestData, err := json.Marshal(request)
	if err != nil {
		logging.Error(err, "An error occurred marshalling request")
		return false, err
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(requestData))
	if err != nil {
		logging.Error(err, "link request create error")
		return false, err
	}

	resp, err := client.HttpClient.Do(req)

	if err != nil {
		logging.Error(err, "Error on link post")
		return false, err
	}

	defer resp.Body.Close()
	return true, nil
}
