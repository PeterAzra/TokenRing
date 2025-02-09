package pinger

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"
	test_utils "tokenRing/pkg/test-utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Ping_ReturnsNodeId_OnSuccessfulPing(t *testing.T) {
	nodeId, _ := uuid.NewUUID()
	respContent, _ := json.Marshal(nodeId.String())

	httpMock := test_utils.NewHttpClientMock(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(respContent)),
		}, nil
	})

	sut := NewPingService(httpMock)

	url, _ := url.Parse("http://localhost:8080")
	id, err := sut.Ping(url)

	assert.Equal(t, nodeId.String(), id.String())
	assert.Nil(t, err)
}

func Test_Ping_ReturnsError_OnError(t *testing.T) {
	httpMock := test_utils.NewHttpClientMock(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
		}, errors.New("test error")
	})

	sut := NewPingService(httpMock)

	url, _ := url.Parse("http://localhost:8080")
	id, err := sut.Ping(url)

	assert.Nil(t, id)
	assert.NotNil(t, err)
}
