package models

type LinkRequest struct {
	Url string
}

func NewLinkRequest(url string) *LinkRequest {
	return &LinkRequest{Url: url}
}
