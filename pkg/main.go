package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"tokenRing/pkg/logging"
	"tokenRing/pkg/node"
	node_api "tokenRing/pkg/node-api"
	node_http "tokenRing/pkg/node-http"
	node_token_service "tokenRing/pkg/node-token-service"
	"tokenRing/pkg/services/startup"

	"github.com/gin-gonic/gin"
)

func main() {
	// TODO read from configuration

	baseNodeAddr := "http://localhost"
	baseNodePort := "8080"
	baseNodeUrl, err := url.Parse(fmt.Sprintf("%v:%v", baseNodeAddr, baseNodePort))
	if err != nil {
		log.Println("Error parsing base node url")
		panic(err)
	}

	stopService := make(chan bool)
	tokenServicePassDelaySeconds := 5

	nodeClient := node_http.NewNodeHttpClient()

	startupService := startup.NewStartupService(nodeClient)

	baseNode, ok := startupService.StartUpBaseNode(baseNodeUrl)

	thisNodePort := baseNodeUrl.Port()
	if !ok {
		// set to 0 to get unused port from the system
		thisNodePort = "0"
	}

	ln, _ := net.Listen("tcp", fmt.Sprintf(":%v", thisNodePort))
	defer ln.Close()

	_, port, _ := net.SplitHostPort(ln.Addr().String())

	if !ok {
		go func() {
			newNodeUrl, err := url.Parse(fmt.Sprintf("%v:%v", baseNodeAddr, port))
			if err != nil {
				log.Println("Unable to parse address for new node")
				panic(err)
			}
			newNode, err := startupService.JoinNodeRing(baseNode, newNodeUrl)
			if err != nil {
				logging.Error(err, "%v Unable to join node ring", newNode.Id)
				panic(err)
			}
			StartTokenService(tokenServicePassDelaySeconds, stopService, newNode, nodeClient)
		}()
	} else {
		StartTokenService(tokenServicePassDelaySeconds, stopService, baseNode, nodeClient)
	}

	nodeApi := node_api.NewNodeApi(nodeClient)

	// TODO middleware for unhandled errors
	r := gin.Default()
	r.GET("/ping", nodeApi.Ping)
	r.POST("/joinrequest", nodeApi.JoinRequest)
	r.POST("/left-link", nodeApi.LeftLink)
	r.POST("/right-link", nodeApi.RightLink)
	r.GET("/state", nodeApi.PrintState)
	r.POST("/token", nodeApi.Token)

	http.Serve(ln, r)

	stopService <- true
}

func StartTokenService(timeDelaySeconds int,
	stopService chan bool,
	nodeToWatch *node.Node,
	nodeClient node_http.NodeClient) *node_token_service.TokenService {

	tknService := node_token_service.NewTokenService(timeDelaySeconds, stopService, nodeToWatch, nodeClient)

	go func() {
		tknService.Run()
	}()

	return tknService
}
