package token_service

import (
	"net/http"
	"net/url"
	"testing"
	"tokenRing/pkg/node"
	test_utils "tokenRing/pkg/test-utils"

	"github.com/stretchr/testify/assert"
)

func Test_SendToken_ReturnsNil_OnSuccessfulRequest(t *testing.T) {
	nodeToUrl, _ := url.Parse("http://localhost:8080")
	nodeFromUrl, _ := url.Parse("http://localhost:8081")
	nodeTo := node.NewNode(nodeToUrl)
	nodeFrom := node.NewNode(nodeFromUrl)

	httpMock := test_utils.NewHttpClientMock(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
		}, nil
	})

	sut := NewTokenService(httpMock)
	result := sut.SendToken(nodeFrom, nodeTo)
	assert.Nil(t, result)
}

func Test_SendToken_ReturnsError_OnBadRequest(t *testing.T) {
	nodeToUrl, _ := url.Parse("http://localhost:8080")
	nodeFromUrl, _ := url.Parse("http://localhost:8081")
	nodeTo := node.NewNode(nodeToUrl)
	nodeFrom := node.NewNode(nodeFromUrl)

	httpMock := test_utils.NewHttpClientMock(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
		}, nil
	})

	sut := NewTokenService(httpMock)
	result := sut.SendToken(nodeFrom, nodeTo)
	assert.NotNil(t, result)
}
