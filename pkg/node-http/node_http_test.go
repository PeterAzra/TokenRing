package node_http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"tokenRing/pkg/models"
	"tokenRing/pkg/node"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type PingReturnsGuidClient struct{}

func (client *PingReturnsGuidClient) Do(request *http.Request) (*http.Response, error) {
	uuid, _ := uuid.NewUUID()
	uuidReader := strings.NewReader(uuid.String())
	uuidResp := io.NopCloser(uuidReader)

	pingResp := http.Response{
		StatusCode: http.StatusOK,
		Body:       uuidResp,
	}
	return &pingResp, nil
}

func Test_PingReturnsGuid_WhenResponseFromNode(t *testing.T) {
	sut := NodeHttpClient{
		HttpClient: &PingReturnsGuidClient{},
	}
	testEndpoint, _ := url.Parse("http://localhost:8080/ping")
	resp, err := sut.PingNode(testEndpoint)
	assert.Nil(t, err)
	assert.NotEqual(t, uuid.Nil, resp, "Ping returned zero guid when expected non-zero guid")
}

type Ping404Client struct{}

func (client *Ping404Client) Do(request *http.Request) (*http.Response, error) {
	return nil, http.ErrHandlerTimeout
}

func Test_PingReturnsZeroGuid_WhenRequestErrors(t *testing.T) {
	sut := NodeHttpClient{
		HttpClient: &Ping404Client{},
	}
	testEndpoint, _ := url.Parse("http://localhost:8080/ping")
	resp, err := sut.PingNode(testEndpoint)
	assert.NotNil(t, err)
	assert.Equal(t, uuid.Nil, resp, "Ping returned non-zero guid when expected zero guid")
}

type JoinValidResponseClient struct{}

func (client *JoinValidResponseClient) Do(request *http.Request) (*http.Response, error) {
	joinResponse := models.NewJoinResponse("http://localhost:8001", "http://localhost:8002")
	bodyData, _ := json.Marshal(joinResponse)
	reader := io.NopCloser(bytes.NewBuffer(bodyData))
	resp := http.Response{
		StatusCode: http.StatusOK,
		Body:       reader,
	}
	return &resp, nil
}

func Test_JoinReturnsResponse_OnValidRequest(t *testing.T) {
	sut := NodeHttpClient{
		HttpClient: &JoinValidResponseClient{},
	}
	testNodeUrl, _ := url.Parse("http://localhost:8083")
	testNode := node.NewNode(testNodeUrl)
	resp, err := sut.Join(testNode, "http://localhost:8080")
	assert.Equal(t, nil, err, "Expected nil error but found err %w", err)
	assert.True(t, resp.Ok)
	assert.True(t, len(resp.Left) > 0, "Expected url but found none on left")
	assert.True(t, len(resp.Right) > 0, "Expected url but found none on right")
}
