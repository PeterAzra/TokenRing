package startup

import (
	"errors"
	"log"
	"net/url"
	"tokenRing/pkg/logging"
	"tokenRing/pkg/node"
	node_http "tokenRing/pkg/node-http"
	node_ring "tokenRing/pkg/node-ring"
	node_token "tokenRing/pkg/node-token"
	node_token_service "tokenRing/pkg/node-token-service"
	"tokenRing/pkg/services/connect"
)

type StartupService struct {
	NodeClient node_http.NodeClient
}

func NewStartupService(nodeClient node_http.NodeClient) *StartupService {
	return &StartupService{NodeClient: nodeClient}
}

func (s *StartupService) StartUpBaseNode(baseNodeUrl *url.URL) (*node.Node, bool) {
	log.Printf("Connecting to node ring %v", baseNodeUrl)
	baseUuid, err := s.NodeClient.PingNode(baseNodeUrl)
	if err != nil {
		log.Println("Node ring not found. Setting node as base node.")
		baseNode := node.InitNode(baseNodeUrl)
		node_ring.InitNodeRing(baseNode)
		baseNode.Left = baseNode
		baseNode.Right = baseNode
		baseNode.Token = node_token.NewToken()
		tknService := node_token_service.NewTokenService(5, baseNode, s.NodeClient)
		go func() {
			tknService.Run()
		}()
		return baseNode, true
	} else {
		log.Printf("Found base node %v", baseUuid)
		baseNode := node.NewNodeWithId(baseNodeUrl, &baseUuid)
		return baseNode, false
	}
}

func (s *StartupService) JoinNodeRing(baseNode *node.Node, thisNodeUrl *url.URL) (*node.Node, error) {
	newNode := node.InitNode(thisNodeUrl)
	node_ring.InitNodeRing(baseNode)

	joinResp, err := s.NodeClient.Join(baseNode.Url, newNode)
	if err != nil {
		log.Printf("unable to join ring")
		return newNode, err
	}

	if joinResp.Ok {
		leftUrl, err := url.Parse(joinResp.Left)
		if err != nil {
			log.Print("unable to parse left node url")
			return newNode, err
		}

		rightUrl, err := url.Parse(joinResp.Right)
		if err != nil {
			log.Print("unable to parse right node url")
			return newNode, err
		}

		leftLinkOk, err := connect.ConnectLeftAdjacentNode(newNode, leftUrl, s.NodeClient)
		if err != nil {
			return newNode, err
		}
		if !leftLinkOk {
			return newNode, errors.New("unable to setup link to left node")
		}

		leftId, err := s.NodeClient.PingNode(leftUrl)
		if err != nil {
			logging.Error(err, "an error occured pinging left node on join")
			return newNode, err
		}
		newNode.Left = node.NewNodeWithId(leftUrl, &leftId)

		rightLinkOk, err := connect.ConnectRightAdjacentNode(newNode, rightUrl, s.NodeClient)
		if err != nil {
			return newNode, err
		}
		if !rightLinkOk {
			return newNode, errors.New("unable to setup link to right node")
		}

		rightId, err := s.NodeClient.PingNode(rightUrl)
		if err != nil {
			logging.Error(err, "an error occurred pining right node on join")
			return newNode, err
		}
		newNode.Right = node.NewNodeWithId(rightUrl, &rightId)
	} else {
		return newNode, errors.New("unable to join node ring")
	}

	logging.Information("Join successful for node %v with left: %v right: %v", newNode.Id, newNode.Left.Url, newNode.Right.Url)

	return newNode, nil
}
