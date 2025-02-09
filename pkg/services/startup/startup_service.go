package startup_service

import (
	"errors"
	"log"
	"net/url"
	"time"
	"tokenRing/pkg/logging"
	"tokenRing/pkg/node"
	node_ring "tokenRing/pkg/node-ring"
	node_token "tokenRing/pkg/node-token"
	join_service "tokenRing/pkg/services/joiner"
	link_service "tokenRing/pkg/services/linker"
	ping_service "tokenRing/pkg/services/pinger"
	token_service "tokenRing/pkg/services/token_sender"
)

type StartupService struct {
	pingSvc  ping_service.Pinger
	joinSvc  join_service.Joiner
	linkSvc  link_service.Linker
	tokenSvc token_service.TokenSender
}

func NewStartupService(pingSvc ping_service.Pinger,
	joinSvc join_service.Joiner,
	linkSvc link_service.Linker,
	tokenSvc token_service.TokenSender) *StartupService {
	return &StartupService{
		pingSvc:  pingSvc,
		joinSvc:  joinSvc,
		linkSvc:  linkSvc,
		tokenSvc: tokenSvc,
	}
}

func (svc *StartupService) StartUpBaseNode(baseNodeUrl *url.URL) (*node.Node, bool) {
	log.Printf("Connecting to node ring %v", baseNodeUrl)
	baseUuid, err := svc.pingSvc.Ping(baseNodeUrl)
	if err != nil {
		log.Println("Node ring not found. Setting node as base node.")
		baseNode := node.InitNode(baseNodeUrl)
		node_ring.InitNodeRing(baseNode)
		baseNode.Left = baseNode
		baseNode.Right = baseNode
		baseNode.Token = node_token.NewToken()

		// TODO copied from token api
		go func() {
			<-time.After(time.Duration(5000 * int(time.Millisecond)))
			if node.Self.Token != nil {
				_ = svc.tokenSvc.SendToken(&node.Self, node.Self.Right)
			}
		}()

		return baseNode, true
	} else {
		log.Printf("Found base node %v", baseUuid)
		baseNode := node.NewNodeWithId(baseNodeUrl, &baseUuid)
		return baseNode, false
	}
}

func (svc *StartupService) JoinNodeRing(baseNode *node.Node, thisNodeUrl *url.URL) (*node.Node, error) {
	newNode := node.InitNode(thisNodeUrl)
	node_ring.InitNodeRing(baseNode)

	joinResp, err := svc.joinSvc.Join(baseNode.Url, newNode)
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

		leftLinkOk, err := svc.linkSvc.ConnectLeftAdjacentNode(newNode, leftUrl)
		if err != nil {
			return newNode, err
		}
		if !leftLinkOk {
			return newNode, errors.New("unable to setup link to left node")
		}

		leftId, err := svc.pingSvc.Ping(leftUrl)
		if err != nil {
			logging.Error(err, "an error occured pinging left node on join")
			return newNode, err
		}
		newNode.Left = node.NewNodeWithId(leftUrl, &leftId)

		rightLinkOk, err := svc.linkSvc.ConnectRightAdjacentNode(newNode, rightUrl)
		if err != nil {
			return newNode, err
		}
		if !rightLinkOk {
			return newNode, errors.New("unable to setup link to right node")
		}

		rightId, err := svc.pingSvc.Ping(rightUrl)
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
