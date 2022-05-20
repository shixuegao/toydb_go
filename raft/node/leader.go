package node

import "sxg/toydb_go/grpc/proto"

type leader struct {
	comm          *Commission
	peerNextIndex map[string]uint64
	peerLastIndex map[string]uint64
}

func (l *leader) step(msg *proto.Message) {

}

func (le *leader) append(p []byte) error {
	return nil
}

func (le *leader) reset(peers []string, lastIndex uint64) {
	le.peerLastIndex = map[string]uint64{}
	le.peerNextIndex = map[string]uint64{}
	for _, peer := range peers {
		le.peerNextIndex[peer] = lastIndex
		le.peerLastIndex[peer] = 0
	}
}

func newLeader(peers []string, lastIndex uint64) *leader {
	le := &leader{}
	le.reset(peers, lastIndex)
	return le
}
