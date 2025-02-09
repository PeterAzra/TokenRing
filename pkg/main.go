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
	join_api "tokenRing/pkg/node-api/join"
	link_api "tokenRing/pkg/node-api/link"
	ping_api "tokenRing/pkg/node-api/ping"
	token_api "tokenRing/pkg/node-api/token"
	node_http "tokenRing/pkg/node-http"
	disconnect_service "tokenRing/pkg/services/disconnect"
	join_service "tokenRing/pkg/services/join"
	link_service "tokenRing/pkg/services/link"
	ping_service "tokenRing/pkg/services/pinger"
	startup_service "tokenRing/pkg/services/startup"
	token_service "tokenRing/pkg/services/token"

	"github.com/gin-gonic/gin"
)

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

	httpClient := node_http.NewHttpClient()

	// Services
	pingSvc := ping_service.NewPingService(httpClient)
	linkSvc := link_service.NewLinkService(httpClient)
	tokenSvc := token_service.NewTokenService(httpClient)
	joinSvc := join_service.NewJoinService(httpClient)
	dcSvc := disconnect_service.NewDisconnectService(tokenSvc, linkSvc)
	startupSvc := startup_service.NewStartupService(pingSvc, joinSvc, linkSvc, tokenSvc)

	// Apis
	linkApi := link_api.NewLinkApi(pingSvc)
	tokenApi := token_api.NewTokenApi(tokenSvc)

	// Startup Server

	baseNode, ok := startupSvc.StartUpBaseNode(baseNodeUrl)

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
			newNode, err := startupSvc.JoinNodeRing(baseNode, newNodeUrl)
			if err != nil {
				logging.Error(err, "%v Unable to join node ring", newNode.Id)
				panic(err)
			}
		}()
	}

	// TODO middleware for unhandled errors?
	router := gin.Default()
	router.GET("/ping", ping_api.Ping)
	router.POST("/joinrequest", join_api.Join)
	router.POST("/left-link", linkApi.LeftLink)
	router.POST("/right-link", linkApi.RightLink)
	// router.GET("/state", nodeApi.PrintState)
	router.POST("/token", tokenApi.Token)

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

	dcSvc.Disconnect(&node.Self)

	log.Println("Server exiting")
}
