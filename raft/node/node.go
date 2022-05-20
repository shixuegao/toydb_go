package node

import (
	"sxg/toydb_go/grpc/proto"
	"sxg/toydb_go/raft"

	"google.golang.org/protobuf/types/known/anypb"
)

type wrapper struct {
	addr  *proto.Address
	event *anypb.Any
}

type Node struct {
	id          string
	peers       []string
	term        uint64
	log         *raft.Log
	msgSender   chan<- *proto.Message
	insSender   chan<- raft.Instruction
	queueReqs   []*wrapper
	proxiedReqs map[string]*proto.Address
	comm        *Commission
}

func (n *Node) abortProxied() error {
	return nil
}

func (n *Node) forwardQueued(leader *proto.Address) error {
	return nil
}

func (n *Node) send(to *proto.Address, event any) error {
	return nil
}

func (n *Node) sendIns(ins raft.Instruction) {
	n.insSender <- ins
}

func (n *Node) registerProxiedRequest(id string, addr *proto.Address) {
	n.proxiedReqs[id] = addr
}

func (n *Node) unregisterProxiedRequest(id string) {
	delete(n.proxiedReqs, id)
}

func (n *Node) storeClientRequestsQueue(addr *proto.Address, event *anypb.Any) {
	wra := &wrapper{
		addr:  addr,
		event: event,
	}
	n.queueReqs = append(n.queueReqs, wra)
}

func (n *Node) resetClientRequestsQueue() []*wrapper {
	queue := n.queueReqs
	n.queueReqs = []*wrapper{}
	return queue
}

func (n *Node) quorum() int {
	return (len(n.peers)+1)/2 + 1
}
