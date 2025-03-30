package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
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
	ping_service "tokenRing/pkg/services/ping"
	startup_service "tokenRing/pkg/services/startup"
	token_service "tokenRing/pkg/services/token"

	"github.com/gin-gonic/gin"
)

const defaultBaseNodeUrl = "http://localhost:8080"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	baseNodeUrl, thisNodePort := readStartUpNodeConfiguration()

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
	baseNode, isStartingAsBaseNode := startupSvc.StartUpBaseNode(baseNodeUrl)

	ln, _ := net.Listen("tcp", fmt.Sprintf(":%v", thisNodePort))
	defer ln.Close()

	_, port, _ := net.SplitHostPort(ln.Addr().String())

	// New node is joining the ring
	if !isStartingAsBaseNode {
		go func() {
			newNodeUrl, err := url.Parse(fmt.Sprintf("http://%v:%v", "localhost", port))
			if err != nil {
				logging.Error(err, "Unable to parse address for new node")
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

func readStartUpNodeConfiguration() (*url.URL, string) {
	for i := 0; i < len(os.Args); i++ {
		logging.Information("Args[%v]: %v", i, os.Args[i])
	}

	argsLen := len(os.Args)

	baseNodeUrl := defaultBaseNodeUrl
	thisNodePort := "0"

	if argsLen > 1 {
		baseNodeUrl = os.Args[1]
	}
	if argsLen > 2 {
		thisNodePort = os.Args[2]
		_, err := strconv.Atoi(thisNodePort)
		if err != nil {
			logging.Error(err, "Error parsing this node port from arguments")
			panic(err)
		}
	}

	baseNode, err := url.Parse(baseNodeUrl)
	if err != nil {
		logging.Error(err, "Error parsing base node url")
		panic(err)
	}

	logging.Information("BaseNodeUrl: %v; ThisNodePort: %v", baseNode, thisNodePort)

	return baseNode, thisNodePort
}
