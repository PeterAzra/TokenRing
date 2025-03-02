package node_token

type NodeToken struct {
}

// TODO token id?
func NewToken() *NodeToken {
	return &NodeToken{}
}

func (t *NodeToken) String() string {
	return "I am the token"
}
