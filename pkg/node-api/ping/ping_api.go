package ping_api

import (
	"log"
	"net/http"
	"tokenRing/pkg/node"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	thisNode := &node.Self
	log.Println(thisNode)
	c.JSON(http.StatusOK, thisNode.Id)
}
