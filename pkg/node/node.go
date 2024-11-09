package node

import (
	"fmt"
	"net/url"
	"tokenRing/pkg/logging"
	node_token "tokenRing/pkg/node-token"

	"github.com/google/uuid"
)

type Node struct {
	Left, Right *Node
	Id          *uuid.UUID
	Url         *url.URL
	Token       *node_token.NodeToken
}

var Self Node

func InitNode(url *url.URL) *Node {
	if Self.Url == nil {
		logging.Information("Initializing node")
		Self = *NewNode(url)
	}
	return &Self
}

func NewNode(url *url.URL) *Node {
	newid := uuid.New()
	return NewNodeWithId(url, &newid)
}

func NewNodeWithId(url *url.URL, id *uuid.UUID) *Node {
	logging.Information("Creating node %v %v", id, url)
	return &Node{
		Id:    id,
		Left:  nil,
		Right: nil,
		Url:   url,
	}
}

func (n *Node) String() string {
	l, r, tkn := "nil", "nil", "DNE"
	if n.Left != nil {
		l = fmt.Sprintf("{ID: %v, URL: %v}", n.Left.Id.String(), n.Left.Url)
	}
	if n.Right != nil {
		r = fmt.Sprintf("{ID: %v, URL: %v}", n.Right.Id, n.Right.Url)
	}
	if n.Token != nil {
		tkn = "Holder of the token"
	}
	return fmt.Sprintf("{ID: %v, Url: %v, Left: %v, Right: %v, Token: %v}", n.Id, n.Url.String(), l, r, tkn)
}
