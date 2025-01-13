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
	Join(url *url.URL, n *node.Node) (*models.JoinResponse, error)
	PingNode(url *url.URL) (uuid.UUID, error)
	LinkNode(url *url.URL, request *models.LinkRequest) (bool, error)
	SendToken(from *node.Node, to *node.Node) error
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
	endpoint := url.JoinPath("ping")
	logging.Information("Pinging node %v", endpoint)

	request, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
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

func (client *NodeHttpClient) Join(url *url.URL, n *node.Node) (*models.JoinResponse, error) {
	logging.Information("Contacting node to join ring")
	joinRequest := models.NewJoinRequest(n.Id, n.Url.String())
	resp, err := client.sendJoinRequest(url, joinRequest)
	if err != nil {
		logging.Warning("An error occurred on join request")
		return nil, err
	}

	logging.Information("Join request successful")

	return resp, nil
}

func (client *NodeHttpClient) sendJoinRequest(url *url.URL, request *models.JoinRequest) (*models.JoinResponse, error) {
	endpoint := url.JoinPath("joinrequest")
	logging.Information("Sending join request %v %v", endpoint, request)

	requestData, err := json.Marshal(request)
	if err != nil {
		logging.Error(err, "An error occurred marshalling join request model")
		return nil, err
	}

	//log.Printf("Join request %v", string(requestData))

	req, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewBuffer(requestData))
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

func (client *NodeHttpClient) LinkNode(url *url.URL, request *models.LinkRequest) (bool, error) {
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

	resp, err := client.HttpClient.Do(req)

	if err != nil {
		logging.Error(err, "Error on link post")
		return false, err
	}

	defer resp.Body.Close()
	return true, nil
}

func (client *NodeHttpClient) SendToken(from *node.Node, to *node.Node) error {
	endpoint := to.Url.JoinPath("token")
	logging.Information("Sending token to node %v %v", endpoint, to.Id)
	request, err := http.NewRequest(http.MethodPost, endpoint.String(), &bytes.Buffer{})

	if err != nil {
		logging.Error(err, "create token request error")
		return err
	}

	resp, err := client.HttpClient.Do(request)

	if err != nil {
		logging.Error(err, "token response error")
		return err
	} else if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("token response was not successful with HttpStatusCode: %v", resp.StatusCode)
		logging.Error(err, "unsuccessful token response")
		return err
	} else {
		from.Token = nil
		logging.Information("Token pass successful")
		return nil
	}
}
