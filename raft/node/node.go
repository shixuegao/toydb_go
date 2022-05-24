package node

import (
	"context"
	"errors"
	"fmt"
	"sxg/toydb_go/grpc/proto"
	logging "sxg/toydb_go/logging"
	"sxg/toydb_go/raft"
	"sync"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
)

type clientRequestWrapper struct {
	addr  *proto.Address
	event *proto.ClientRequestEvent
}

type Node struct {
	context         context.Context
	id              string
	peers           []string
	term            uint64
	log             *raft.Log
	msgSender       chan<- *proto.Message
	insSender       chan<- raft.Instruction
	clientRequests  []*clientRequestWrapper
	proxiedRequests map[string]*proto.Address
	logger          logging.Logger
	comm            *Commission
	wg              *sync.WaitGroup
	dri             *raft.Driver
}

func (n *Node) abortProxied(reason string) error {
	requests := n.proxiedRequests
	n.proxiedRequests = map[string]*proto.Address{}
	abort, _ := anypb.New(&proto.Abort{Reason: reason})
	for _, addr := range requests {
		event := &proto.ClientResponseEvent{
			Id:      n.id,
			Content: abort,
		}
		if err := n.send(addr, event); err != nil {
			n.logger.Errorf("sending abort to client failed: %s", err.Error())
		}
	}
	return nil
}

func (n *Node) forwardQueued(leader *proto.Address) error {
	wrappers := n.resetClientRequests()
	for _, wrapper := range wrappers {
		v := wrapper.event
		n.registerProxiedRequest(v.Id, wrapper.addr)
		addr := wrapper.addr
		if addr.AddrType == proto.AddrType_CLIENT {
			addr = &proto.Address{AddrType: proto.AddrType_LOCAL}
		}
		eve, _ := anypb.New(v)
		msg := &proto.Message{
			Term:  0,
			From:  addr,
			To:    leader,
			Event: eve,
		}
		n.msgSender <- msg
	}

	return nil
}

//发送信息到peer节点
func (n *Node) send(to *proto.Address, event any) (err error) {
	from := &proto.Address{AddrType: proto.AddrType_LOCAL}
	var v protoreflect.ProtoMessage
	var eve *anypb.Any
	if vv, ok := event.(protoreflect.ProtoMessage); !ok {
		return errors.New("invalid event type")
	} else {
		v = vv
	}
	if eve, err = anypb.New(v); err != nil {
		return
	}
	msg := &proto.Message{
		Term:  n.term,
		From:  from,
		To:    to,
		Event: eve,
	}
	n.msgSender <- msg
	return nil
}

func (n *Node) sendIns(ins raft.Instruction) {
	n.insSender <- ins
}

func (n *Node) registerProxiedRequest(id string, addr *proto.Address) {
	n.proxiedRequests[id] = addr
}

func (n *Node) unregisterProxiedRequest(id string) {
	delete(n.proxiedRequests, id)
}

func (n *Node) storeClientRequests(addr *proto.Address, event *proto.ClientRequestEvent) {
	wrapper := &clientRequestWrapper{addr, event}
	n.clientRequests = append(n.clientRequests, wrapper)
}

func (n *Node) resetClientRequests() []*clientRequestWrapper {
	requests := n.clientRequests
	n.clientRequests = []*clientRequestWrapper{}
	return requests
}

func (n *Node) quorum() int {
	return (len(n.peers)+1)/2 + 1
}

func NewNode(context context.Context, id string, peers []string, msgSender chan<- *proto.Message,
	log *raft.Log, state raft.State, logger logging.Logger) (*Node, error) {
	appliedIndex := state.AppliedIndex()
	commitIndex, _ := log.CommitIndexAndTerm()
	if appliedIndex > commitIndex {
		return nil, fmt.Errorf("state machine applied index %d greater than log commit %d", appliedIndex, commitIndex)
	}
	insSender := make(chan raft.Instruction)
	term, votedFor := log.LoadTerm()
	node := &Node{
		context:         context,
		id:              id,
		peers:           peers,
		term:            term,
		log:             log,
		logger:          logger,
		msgSender:       msgSender,
		insSender:       insSender,
		clientRequests:  []*clientRequestWrapper{},
		proxiedRequests: map[string]*proto.Address{},
		wg:              new(sync.WaitGroup),
	}
	//commission
	comm := newCommission(node, logger)
	(comm.curRole().(*follower)).votedFor = votedFor
	if len(peers) == 0 {
		logger.Info("no peers specified, starting as leader")
		comm.setRole(Leader)
		lastIndex, _ := log.LastIndexAndTerm()
		(comm.curRole().(*leader)).reset(peers, lastIndex)
	}
	//driver
	node.dri = raft.NewDriver(node.context, insSender, node.msgSender, logger)
	node.dri.Drive(state, node.wg)
	return node, nil
}
