package join_api

import (
	"log"
	"net/http"
	"tokenRing/pkg/models"
	"tokenRing/pkg/node"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Join(c *gin.Context) {
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
