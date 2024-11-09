package node_http

import "net/http"

type HttpSender interface {
	Do(req *http.Request) (*http.Response, error)
}

type HttpClient struct{}

func (client *HttpClient) Do(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}

func NewHttpClient() *HttpClient {
	return &HttpClient{}
}
