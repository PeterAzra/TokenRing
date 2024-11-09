package node_api

import (
	"log"
	"net/http"
	"net/url"
	"tokenRing/pkg/logging"
	"tokenRing/pkg/models"
	"tokenRing/pkg/node"
	node_http "tokenRing/pkg/node-http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type NodeApi struct {
	NodeClient node_http.NodeClient
}

func NewNodeApi(httpClient node_http.NodeClient) *NodeApi {
	return &NodeApi{
		NodeClient: httpClient,
	}
}

// curl http://localhost:8080/joinrequest -X POST -v -H 'content-type: application/json' -d '{"NodeId":"00000001-0001-0001-0001-000000000001","Url":"asdf"}'
func (api *NodeApi) JoinRequest(c *gin.Context) {
	thisNode := &node.Self

	var joinRequestData models.JoinRequest
	if err := c.ShouldBindWith(&joinRequestData, binding.JSON); err != nil {
		log.Printf("Error on JoinRequest content read %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Received Join Request from %v", joinRequestData)

	c.JSON(http.StatusOK, models.JoinResponse{
		Ok:    true,
		Left:  thisNode.Url.String(),
		Right: thisNode.Right.Url.String(),
	})
}

func (api *NodeApi) Ping(c *gin.Context) {
	thisNode := &node.Self
	log.Println(thisNode)
	c.JSON(http.StatusOK, thisNode.Id)
}

// curl http://localhost:8080/right-link -X POST -v -H 'content-type: application/json' -d '{"Url":"http://localhost:42163"}'
func (api *NodeApi) RightLink(c *gin.Context) {
	var linkRequestData models.LinkRequest
	if err := c.ShouldBindWith(&linkRequestData, binding.JSON); err != nil {
		logging.Error(err, "Error on LinkRequest content read %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}

	newNode, err := linkNode(api, linkRequestData.Url)
	if err != nil {
		logging.Error(err, "Error linking node %v", linkRequestData.Url)
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to link node"})
		return
	}

	logging.Information("Updating right node %v %v", newNode.Id, newNode.Url)

	node.Self.Right = newNode
	c.JSON(http.StatusOK, gin.H{})
}

func (api *NodeApi) LeftLink(c *gin.Context) {
	var linkRequestData models.LinkRequest
	if err := c.ShouldBindWith(&linkRequestData, binding.JSON); err != nil {
		logging.Error(err, "Error on LinkRequest content read %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}

	newNode, err := linkNode(api, linkRequestData.Url)
	if err != nil {
		logging.Error(err, "Error linking node %v", linkRequestData.Url)
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to link node"})
		return
	}

	logging.Information("Updating left node %v %v", newNode.Id, newNode.Url)

	node.Self.Left = newNode
	c.JSON(http.StatusOK, gin.H{})
}

func linkNode(api *NodeApi, newLinkUrl string) (*node.Node, error) {
	linkUrl, err := url.Parse(newLinkUrl)
	if err != nil {
		logging.Error(err, "Unable to parse link url %v", newLinkUrl)
		return nil, err
	}

	linkId, err := api.NodeClient.PingNode(linkUrl)
	if err != nil {
		logging.Error(err, "Unable to ping node on link request")
		return nil, err
	}

	return node.NewNodeWithId(linkUrl, &linkId), nil
}

func (n *NodeApi) PrintState(c *gin.Context) {
	msg := node.Self.String()
	logging.Information(msg)
	c.JSON(http.StatusOK, gin.H{"node": msg})
}
