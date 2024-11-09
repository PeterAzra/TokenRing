package models

import (
	"github.com/google/uuid"
)

type JoinRequest struct {
	NodeId *uuid.UUID
	Url    string
}

func NewJoinRequest(nodeId *uuid.UUID, url string) *JoinRequest {
	return &JoinRequest{
		NodeId: nodeId,
		Url:    url,
	}
}
