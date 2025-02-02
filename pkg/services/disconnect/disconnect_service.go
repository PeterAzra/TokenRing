package disconnect_service

import (
	"tokenRing/pkg/logging"
	"tokenRing/pkg/models"
	"tokenRing/pkg/node"
	link_service "tokenRing/pkg/services/link"
	token_service "tokenRing/pkg/services/token"
)

type DisconnectService struct {
	tokenSvc *token_service.TokenService
	linkSvc  *link_service.LinkService
}

func NewDisconnectService(tokenSvc *token_service.TokenService, linkSvc *link_service.LinkService) *DisconnectService {
	return &DisconnectService{
		tokenSvc: tokenSvc,
		linkSvc:  linkSvc,
	}
}

func (svc *DisconnectService) Disconnect(disconnectingNode *node.Node) (bool, error) {
	logging.Information("Node %v disconnecting from ring", disconnectingNode.Id)

	if disconnectingNode.Token != nil {
		if err := svc.tokenSvc.SendToken(disconnectingNode, disconnectingNode.Right); err != nil {
			// TODO better handling for sending token on disconnect?
			// TODO Possible race condition with token service
			logging.Error(err, "An error occurred on disconnect when sending token. The token is lost!")
		}
	}

	// Send requests to left and right adj nodes to use new links
	rightLinkRequest := models.NewLinkRequest(disconnectingNode.Right.Url.String())
	leftAdjUrl := disconnectingNode.Left.Url.JoinPath("right-link")
	leftDoneCh := make(chan error, 1)

	leftLinkRequest := models.NewLinkRequest(disconnectingNode.Left.Url.String())
	rightAdjUrl := disconnectingNode.Right.Url.JoinPath("left-link")
	rightDoneCh := make(chan error, 1)

	// Send requests together to try and avoid issues while swapping links
	go func() {
		_, err := svc.linkSvc.LinkNode(leftAdjUrl, rightLinkRequest)
		leftDoneCh <- err
	}()
	go func() {
		_, err := svc.linkSvc.LinkNode(rightAdjUrl, leftLinkRequest)
		rightDoneCh <- err
	}()

	leftErr := <-leftDoneCh
	rightErr := <-rightDoneCh

	// TODO better checking of errors
	if leftErr != nil {
		logging.Error(leftErr, "An error occurred disconnecting from node ring")
		return false, leftErr
	}
	if rightErr != nil {
		logging.Error(rightErr, "An error occurred disconnecting from node ring")
		return false, rightErr
	}

	logging.Information("Node successfully disconnected")
	return true, nil
}
