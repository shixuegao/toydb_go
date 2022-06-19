package node

import (
	"sxg/toydb_go/grpc/proto"
	"sxg/toydb_go/raft"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/anypb"
)

type candidate struct {
	comm  *commission
	votes int
}

func (can *candidate) becomeFollower(leader string, term uint64) error {
	can.comm.setTerm(term)
	can.comm.log().SaveTerm(term, "")
	//转变为follower，但不会投票
	can.comm.setRole(Follower)
	fol := can.comm.curRole().(*follower)
	fol.leader = leader
	fol.votedFor = ""
	addr := &proto.Address{
		AddrType: proto.AddrType_PEER,
		Peer:     leader,
	}
	if err := can.comm.node.abortProxied("candidate converts to follower"); err != nil {
		return err
	}
	if err := can.comm.node.forwardQueued(addr); err != nil {
		return err
	}
	return nil
}

func (can *candidate) becomeLeader() error {
	can.comm.setRole(Leader)
	lea := can.comm.curRole().(*leader)
	lastIndex, _ := can.comm.log().LastIndexAndTerm()
	lea.reset(can.comm.node.peers, lastIndex)

	commitIndex, commitTerm := can.comm.log().CommitIndexAndTerm()
	addr := &proto.Address{AddrType: proto.AddrType_PEERS}
	event := &proto.HeartbeatEvent{
		CommitIndex: commitIndex,
		CommitTerm:  commitTerm,
	}
	if err := lea.comm.node.send(0, addr, event); err != nil {
		return err
	}
	if _, err := lea.append([]byte{}); err != nil {
		return err
	}
	if err := lea.comm.node.abortProxied("candidate converts to leader"); err != nil {
		return err
	}
	return nil
}

func (can *candidate) step(cas *raft.Case) {
	msg := cas.Msg
	if msg.Term > can.comm.term() &&
		msg.From.AddrType == proto.AddrType_PEER {
		if err := can.becomeFollower(msg.From.Peer, msg.Term); err != nil {
			can.comm.logger.Errorf("candidate convert to follower failed: %+v", errors.WithStack(err))
			return
		}
		fol := can.comm.curRole().(*follower)
		fol.step(cas)
		return
	}
	switch v := cas.Event.(type) {
	case *proto.HeartbeatEvent:
		if msg.From.AddrType == proto.AddrType_PEER {
			if err := can.becomeFollower(msg.From.Peer, msg.Term); err != nil {
				can.comm.logger.Errorf("candidate convert to follower failed: %+v", errors.WithStack(err))
				return
			}
			fol := can.comm.curRole().(*follower)
			fol.step(cas)
			return
		}
	case *proto.GrantVoteEvent:
		can.comm.logger.Debugf("received term %d vote from %s", can.comm.term(), msg.From.Peer)
		can.votes++
		if can.votes > can.comm.node.quorum() {
			requests := can.comm.node.resetClientRequests()
			can.comm.logger.Infof("won election for term %d, become leader", can.comm.term())
			if err := can.becomeLeader(); err != nil {
				can.comm.logger.Errorf("candidate convert to leader failed: %+v", errors.WithStack(err))
				return
			}
			//当前节点升级为leader，将客户端的请求转发至本地
			for _, request := range requests {
				event, _ := anypb.New(request.event)
				msg := &proto.Message{
					Term:  0,
					From:  request.addr,
					To:    &proto.Address{AddrType: proto.AddrType_LOCAL},
					Event: event,
				}
				lea := can.comm.curRole().(*leader)
				lea.step(&raft.Case{ID: cas.ID, Msg: msg, Event: event, RespContent: nil})
			}
		}
	case *proto.ClientRequestEvent:
		can.comm.node.storeClientRequests(msg.From, v)
	case *proto.ClientResponseEvent:
		if v.ResponseType == proto.RespType_RESP_STATUS {
			status := cas.RespContent.(*proto.Status)
			status.Server = can.comm.id()
			v.Content, _ = anypb.New(status)
		}
		can.comm.node.unregisterProxiedRequest(v.Id)
		addr := &proto.Address{AddrType: proto.AddrType_CLIENT}
		if err := can.comm.node.send(cas.ID, addr, v); err != nil {
			can.comm.logger.Errorf("send client response event failed: %+v", errors.WithStack(err))
		}
	case *proto.SolicitVoteEvent: //do nothing
	case *proto.ConfirmleaderEvent, *proto.ReplicateEntriesEvent, *proto.AcceptEntriesEvent, *proto.RejectEntriesEvent:
		can.comm.logger.Warnf("received unexpected message: %s", msg.String())
	}
}
