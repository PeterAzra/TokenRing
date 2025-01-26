package snapshot

import (
	"tokenRing/pkg/logging"
	"tokenRing/pkg/node"
	node_http "tokenRing/pkg/node-http"
)

type SnapshotService struct {
	SnapshotHistory map[int]int
}

func NewSnapshotService() *SnapshotService {
	history := make(map[int]int)
	return &SnapshotService{
		SnapshotHistory: history,
	}
}

func (svc *SnapshotService) Snapshot(node *node.Node, forwardTo *node.Node, nodeClient node_http.NodeClient, snapshotId int) {
	logging.Information("%s", node.String())
	_, ok := svc.SnapshotHistory[snapshotId]
	if !ok {
		nodeClient.SendSnapshot(forwardTo, snapshotId)
		svc.SnapshotHistory[snapshotId] = snapshotId
	}
}
