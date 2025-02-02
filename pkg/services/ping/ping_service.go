package ping_service

import (
	"io"
	"net/http"
	"net/url"
	"tokenRing/pkg/logging"
	node_http "tokenRing/pkg/node-http"

	"github.com/google/uuid"
)

type PingService struct {
	client *node_http.HttpClient
}

func NewPingService(client *node_http.HttpClient) *PingService {
	return &PingService{
		client: client,
	}
}

func (s *PingService) Ping(url *url.URL) (uuid.UUID, error) {
	endpoint := url.JoinPath("ping")
	logging.Information("Pinging node %v", endpoint)

	request, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		logging.Error(err, "create ping request fail")
		return uuid.Nil, err
	}

	resp, err := s.client.Do(request)
	if err != nil {
		logging.Warning("%s", err.Error())
		return uuid.Nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Warning("%s", err.Error())
		return uuid.Nil, err
	}

	uuidResp, err := uuid.ParseBytes(body)
	if err != nil {
		logging.Warning("%s", err.Error())
		return uuid.Nil, err
	}

	logging.Information("Ping successful %v\n", uuidResp.String())
	return uuidResp, nil
}
