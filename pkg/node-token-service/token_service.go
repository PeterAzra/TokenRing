package node_token_service

import (
	"time"
	"tokenRing/pkg/node"
	node_http "tokenRing/pkg/node-http"
)

type TokenService struct {
	TokenPassDelaySeconds int
	// Done                  chan bool
	Node     *node.Node
	NodeHttp node_http.NodeClient
}

func NewTokenService(tokenPassDelaySeconds int,
	// done chan bool,
	nodeToWatch *node.Node,
	nodeClient node_http.NodeClient) *TokenService {
	return &TokenService{
		TokenPassDelaySeconds: tokenPassDelaySeconds,
		// Done:                  done,
		Node:     nodeToWatch,
		NodeHttp: nodeClient,
	}
}

func (svc *TokenService) Run() error {
	ticker := time.NewTicker(time.Duration(svc.TokenPassDelaySeconds * 1000 * int(time.Millisecond)))

	for {
		select {
		// case <-svc.Done:
		// 	return nil
		case <-ticker.C:
			if svc.Node.Token != nil {
				_ = svc.NodeHttp.SendToken(svc.Node, svc.Node.Right)
			}
		}
	}
}
