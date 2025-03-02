package node_http

// import (
// 	"bytes"
// 	"encoding/json"
// 	"io"
// 	"net/http"
// 	"net/url"
// 	"strings"
// 	"testing"
// 	"tokenRing/pkg/models"
// 	"tokenRing/pkg/node"
// 	node_token "tokenRing/pkg/node-token"

// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// )

// func Test_PingReturnsGuid_WhenResponseFromNode(t *testing.T) {
// 	sut := NodeHttpClient{
// 		HttpClient: &PingReturnsGuidClient{},
// 	}
// 	testEndpoint, _ := url.Parse("http://localhost:8080/ping")
// 	resp, err := sut.PingNode(testEndpoint)
// 	assert.Nil(t, err)
// 	assert.NotEqual(t, uuid.Nil, resp, "Ping returned zero guid when expected non-zero guid")
// }

// type Ping404Client struct{}

// func (client *Ping404Client) Do(request *http.Request) (*http.Response, error) {
// 	return nil, http.ErrHandlerTimeout
// }

// func Test_PingReturnsZeroGuid_WhenRequestErrors(t *testing.T) {
// 	sut := NodeHttpClient{
// 		HttpClient: &Ping404Client{},
// 	}
// 	testEndpoint, _ := url.Parse("http://localhost:8080/ping")
// 	resp, err := sut.PingNode(testEndpoint)
// 	assert.NotNil(t, err)
// 	assert.Equal(t, uuid.Nil, resp, "Ping returned non-zero guid when expected zero guid")
// }

// func Test_JoinReturnsResponse_OnValidRequest(t *testing.T) {
// 	sut := NodeHttpClient{
// 		HttpClient: &JoinValidResponseClient{},
// 	}
// 	testNodeUrl, _ := url.Parse("http://localhost:8083")
// 	testNode := node.NewNode(testNodeUrl)
// 	baseNodeUrl, _ := url.Parse("http://localhost:8080")
// 	resp, err := sut.Join(baseNodeUrl, testNode)
// 	assert.Equal(t, nil, err, "Expected nil error but found err %w", err)
// 	assert.True(t, resp.Ok)
// 	assert.True(t, len(resp.Left) > 0, "Expected url but found none on left")
// 	assert.True(t, len(resp.Right) > 0, "Expected url but found none on right")
// }

// func Test_SendToken_ReturnsOk_OnSuccessfulRequest(t *testing.T) {
// 	sut := NodeHttpClient{
// 		HttpClient: &TokenRequestSuccessfulClient{},
// 	}
// 	sourceUrl, _ := url.Parse("http://localhost:8080")
// 	destUrl, _ := url.Parse("http://localhost:8081")
// 	sourceNode := node.NewNode(sourceUrl)
// 	sourceNode.Token = node_token.NewToken()
// 	destNode := node.NewNode(destUrl)
// 	err := sut.SendToken(sourceNode, destNode)

// 	assert.True(t, err == nil)
// 	assert.True(t, sourceNode.Token == nil)
// }

// func Test_SendToken_ReturnsError_OnBadRequest(t *testing.T) {
// 	sut := NodeHttpClient{
// 		HttpClient: &TokenRequestBadRequestClient{},
// 	}
// 	sourceUrl, _ := url.Parse("http://localhost:8080")
// 	destUrl, _ := url.Parse("http://localhost:8081")
// 	sourceNode := node.NewNode(sourceUrl)
// 	sourceNode.Token = node_token.NewToken()
// 	destNode := node.NewNode(destUrl)
// 	err := sut.SendToken(sourceNode, destNode)

// 	assert.True(t, err != nil)
// 	assert.True(t, sourceNode.Token != nil)
// }

// type PingReturnsGuidClient struct{}

// func (client *PingReturnsGuidClient) Do(request *http.Request) (*http.Response, error) {
// 	uuid, _ := uuid.NewUUID()
// 	uuidReader := strings.NewReader(uuid.String())
// 	uuidResp := io.NopCloser(uuidReader)

// 	pingResp := http.Response{
// 		StatusCode: http.StatusOK,
// 		Body:       uuidResp,
// 	}
// 	return &pingResp, nil
// }

// type JoinValidResponseClient struct{}

// func (client *JoinValidResponseClient) Do(request *http.Request) (*http.Response, error) {
// 	joinResponse := models.NewJoinResponse("http://localhost:8001", "http://localhost:8002")
// 	bodyData, _ := json.Marshal(joinResponse)
// 	reader := io.NopCloser(bytes.NewBuffer(bodyData))
// 	resp := http.Response{
// 		StatusCode: http.StatusOK,
// 		Body:       reader,
// 	}
// 	return &resp, nil
// }

// type TokenRequestSuccessfulClient struct{}

// func (client *TokenRequestSuccessfulClient) Do(request *http.Request) (*http.Response, error) {
// 	return &http.Response{
// 		StatusCode: http.StatusOK,
// 	}, nil
// }

// type TokenRequestBadRequestClient struct{}

// func (client *TokenRequestBadRequestClient) Do(request *http.Request) (*http.Response, error) {
// 	return &http.Response{
// 		StatusCode: http.StatusBadRequest,
// 	}, nil
// }
