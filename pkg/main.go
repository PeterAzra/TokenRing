package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os/signal"
	"syscall"
	"time"
	"tokenRing/pkg/logging"
	"tokenRing/pkg/node"
	node_api "tokenRing/pkg/node-api"
	node_http "tokenRing/pkg/node-http"
	"tokenRing/pkg/services/disconnect"
	"tokenRing/pkg/services/startup"

	"github.com/gin-gonic/gin"
)

var nodeClient node_http.NodeClient

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// TODO read from configuration

	baseNodeAddr := "http://localhost"
	baseNodePort := "8080"
	baseNodeUrl, err := url.Parse(fmt.Sprintf("%v:%v", baseNodeAddr, baseNodePort))
	if err != nil {
		log.Println("Error parsing base node url")
		panic(err)
	}

	nodeClient = node_http.NewNodeHttpClient()

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

	// New node is joining the ring
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
		}()
	}

	nodeApi := node_api.NewNodeApi(nodeClient)

	// TODO middleware for unhandled errors?
	router := gin.Default()
	router.GET("/ping", nodeApi.Ping)
	router.POST("/joinrequest", nodeApi.JoinRequest)
	router.POST("/left-link", nodeApi.LeftLink)
	router.POST("/right-link", nodeApi.RightLink)
	router.GET("/state", nodeApi.PrintState)
	router.POST("/token", nodeApi.Token)

	server := &http.Server{
		Addr:    ln.Addr().String(),
		Handler: router,
	}

	go func() {
		if err := server.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	disconnect.DisconnectNode(&node.Self, nodeClient)

	log.Println("Server exiting")
}
