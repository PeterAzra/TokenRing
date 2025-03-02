package link_api

import (
	"net/http"
	"net/url"
	"tokenRing/pkg/logging"
	"tokenRing/pkg/models"
	"tokenRing/pkg/node"
	ping_service "tokenRing/pkg/services/ping"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type LinkApi struct {
	pingSvc *ping_service.PingService
}

func NewLinkApi(pingSvc *ping_service.PingService) *LinkApi {
	return &LinkApi{
		pingSvc: pingSvc,
	}
}

// curl http://localhost:8080/right-link -X POST -v -H 'content-type: application/json' -d '{"Url":"http://localhost:42163"}'
func (api *LinkApi) RightLink(c *gin.Context) {
	var linkRequestData models.LinkRequest
	if err := c.ShouldBindWith(&linkRequestData, binding.JSON); err != nil {
		logging.Error(err, "Error on LinkRequest content read %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}

	newNode, err := linkNode(api.pingSvc, linkRequestData.Url)
	if err != nil {
		logging.Error(err, "Error linking node %v", linkRequestData.Url)
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to link node"})
		return
	}

	logging.Information("Updating right node %v %v", newNode.Id, newNode.Url)

	node.Self.Right = newNode
	c.JSON(http.StatusOK, gin.H{})
}

func (api *LinkApi) LeftLink(c *gin.Context) {
	var linkRequestData models.LinkRequest
	if err := c.ShouldBindWith(&linkRequestData, binding.JSON); err != nil {
		logging.Error(err, "Error on LinkRequest content read %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}

	newNode, err := linkNode(api.pingSvc, linkRequestData.Url)
	if err != nil {
		logging.Error(err, "Error linking node %v", linkRequestData.Url)
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to link node"})
		return
	}

	logging.Information("Updating left node %v %v", newNode.Id, newNode.Url)

	node.Self.Left = newNode
	c.JSON(http.StatusOK, gin.H{})
}

func linkNode(pingSvc *ping_service.PingService, newLinkUrl string) (*node.Node, error) {
	linkUrl, err := url.Parse(newLinkUrl)
	if err != nil {
		logging.Error(err, "Unable to parse link url %v", newLinkUrl)
		return nil, err
	}

	linkId, err := pingSvc.Ping(linkUrl)
	if err != nil {
		logging.Error(err, "Unable to ping node on link request")
		return nil, err
	}

	return node.NewNodeWithId(linkUrl, &linkId), nil
}
