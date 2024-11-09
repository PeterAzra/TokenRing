package node_ring

import (
	"log"
	"tokenRing/pkg/node"
)

type NodeRing struct {
	BaseNode *node.Node
	ThisNode *node.Node
	nodes    []*node.Node
}

var Context NodeRing

func NewNodeRing() *NodeRing {
	return &NodeRing{}
}

func InitNodeRing(baseNode *node.Node) {
	if Context.BaseNode == nil {
		log.Printf("Initializing node ring with base node at %v", baseNode.Url.String())
		Context = *NewNodeRing()
		Context.BaseNode = baseNode
	}
}

func (ring *NodeRing) GetBaseNodeUrl() string {
	return ring.BaseNode.Url.String()
}
