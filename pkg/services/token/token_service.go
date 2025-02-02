package token_service

import (
	"bytes"
	"fmt"
	"net/http"
	"tokenRing/pkg/logging"
	"tokenRing/pkg/node"
	node_http "tokenRing/pkg/node-http"
)

type TokenService struct {
	client *node_http.HttpClient
}

func NewTokenService(client *node_http.HttpClient) *TokenService {
	return &TokenService{
		client: client,
	}
}

func (s *TokenService) SendToken(from *node.Node, to *node.Node) error {
	endpoint := to.Url.JoinPath("token")
	logging.Information("Sending token to node %v %v", endpoint, to.Id)
	request, err := http.NewRequest(http.MethodPost, endpoint.String(), &bytes.Buffer{})

	if err != nil {
		logging.Error(err, "create token request error")
		return err
	}

	resp, err := s.client.Do(request)

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
