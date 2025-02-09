package joiner

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"tokenRing/pkg/logging"
	"tokenRing/pkg/models"
	"tokenRing/pkg/node"
	node_http "tokenRing/pkg/node-http"
)

type Joiner interface {
	Join(url *url.URL, n *node.Node) (*models.JoinResponse, error)
}

type JoinService struct {
	client node_http.HttpSender
}

func NewJoinService(client node_http.HttpSender) *JoinService {
	return &JoinService{
		client: client,
	}
}

func (s *JoinService) Join(url *url.URL, n *node.Node) (*models.JoinResponse, error) {
	logging.Information("Contacting node to join ring")
	joinRequest := models.NewJoinRequest(n.Id, n.Url.String())
	resp, err := sendJoinRequest(url, joinRequest, s.client)
	if err != nil {
		logging.Warning("An error occurred on join request")
		return nil, err
	}

	logging.Information("Join request successful")

	return resp, nil
}

func sendJoinRequest(url *url.URL, request *models.JoinRequest, client node_http.HttpSender) (*models.JoinResponse, error) {
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

	resp, err := client.Do(req)
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
