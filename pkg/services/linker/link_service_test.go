package linker

import (
	"errors"
	"net/http"
	"net/url"
	"testing"
	"tokenRing/pkg/models"
	test_utils "tokenRing/pkg/test-utils"

	"github.com/stretchr/testify/assert"
)

func Test_LinkNode_ReturnsTrue_OnSuccessfulLink(t *testing.T) {
	nodeUrl, _ := url.Parse("http://localhost:8080")
	adjUrl, _ := url.Parse("http://localhost:8081")

	linkRequest := models.NewLinkRequest(nodeUrl.String())

	mockHttp := test_utils.NewHttpClientMock(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
		}, nil
	})

	sut := NewLinkService(mockHttp)
	result, err := sut.LinkNode(adjUrl, linkRequest)

	assert.True(t, result)
	assert.Nil(t, err)
}

func Test_ConnectLeftAdjacentNode_ReturnsFalse_OnInvalidLinkRequest(t *testing.T) {
	nodeUrl, _ := url.Parse("http://localhost:8080")
	adjUrl, _ := url.Parse("http://localhost:8081")
	linkRequest := models.NewLinkRequest(nodeUrl.String())

	mockHttp := test_utils.NewHttpClientMock(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
		}, nil
	})

	sut := NewLinkService(mockHttp)
	result, err := sut.LinkNode(adjUrl, linkRequest)

	assert.False(t, result)
	assert.Nil(t, err)
}

func Test_ConnectLeftAdjacentNode_ReturnsError_OnHttpError(t *testing.T) {
	nodeUrl, _ := url.Parse("http://localhost:8080")
	adjUrl, _ := url.Parse("http://localhost:8081")
	linkRequest := models.NewLinkRequest(nodeUrl.String())

	mockHttp := test_utils.NewHttpClientMock(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("test error")
	})

	sut := NewLinkService(mockHttp)
	result, err := sut.LinkNode(adjUrl, linkRequest)

	assert.False(t, result)
	assert.NotNil(t, err)
}
