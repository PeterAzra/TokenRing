package test_utils

import (
	"net/http"
)

type HttpClientMock struct {
	DoFunc func(*http.Request) (*http.Response, error)
}

func NewHttpClientMock(doFunc func(*http.Request) (*http.Response, error)) *HttpClientMock {
	return &HttpClientMock{
		DoFunc: doFunc,
	}
}

func (client *HttpClientMock) Do(req *http.Request) (*http.Response, error) {
	return client.DoFunc(req)
}
