package join

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"
	"tokenRing/pkg/models"
	"tokenRing/pkg/node"
	test_utils "tokenRing/pkg/test-utils"

	"github.com/stretchr/testify/assert"
)

func Test_Join_ReturnsSuccessfulResponse_OnSuccessfulJoin(t *testing.T) {
	baseNodeUrl, _ := url.Parse("http://localhost:8080")
	joiningNodeUrl, _ := url.Parse("http://localhost:8081")
	joiningNode := node.NewNode(joiningNodeUrl)

	leftUrl := "http://localhost:8082"
	rightUrl := "http://localhost:8083"
	responseModel := models.NewJoinResponse(leftUrl, rightUrl)
	responseJson, _ := json.Marshal(responseModel)
	response := http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseJson)),
	}
	httpMock := test_utils.NewHttpClientMock(func(req *http.Request) (*http.Response, error) {
		return &response, nil
	})

	sut := NewJoinService(httpMock)
	resp, err := sut.Join(baseNodeUrl, joiningNode)

	assert.NotNil(t, resp)
	assert.Equal(t, responseModel.Left, resp.Left)
	assert.Equal(t, responseModel.Right, resp.Right)
	assert.Nil(t, err)
}

func Test_Join_ReturnsError_OnHttpError(t *testing.T) {
	baseNodeUrl, _ := url.Parse("http://localhost:8080")
	joiningNodeUrl, _ := url.Parse("http://localhost:8081")
	joiningNode := node.NewNode(joiningNodeUrl)

	httpMock := test_utils.NewHttpClientMock(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("test error")
	})

	sut := NewJoinService(httpMock)
	resp, err := sut.Join(baseNodeUrl, joiningNode)

	assert.Nil(t, resp)
	assert.NotNil(t, err)
}
