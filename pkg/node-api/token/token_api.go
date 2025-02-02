package token_api

import (
	"net/http"
	"time"
	"tokenRing/pkg/logging"
	"tokenRing/pkg/node"
	node_token "tokenRing/pkg/node-token"
	token_service "tokenRing/pkg/services/token"

	"github.com/gin-gonic/gin"
)

type TokenApi struct {
	tokenSvc *token_service.TokenService
}

func NewTokenApi(tokenSvc *token_service.TokenService) *TokenApi {
	return &TokenApi{
		tokenSvc: tokenSvc,
	}
}

// curl http://localhost:8080/token -X POST -v
func (api *TokenApi) Token(c *gin.Context) {
	if node.Self.Token != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node already has the token!"})
	} else {
		logging.Information("Node %v received the token", node.Self.Id)
		node.Self.Token = node_token.NewToken()

		go func() {
			<-time.After(time.Duration(5000 * int(time.Millisecond)))
			if node.Self.Token != nil {
				_ = api.tokenSvc.SendToken(&node.Self, node.Self.Right)
			}
		}()
		c.JSON(http.StatusOK, gin.H{})
	}
}
