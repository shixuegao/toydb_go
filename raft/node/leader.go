package node

import (
	"fmt"
	"sxg/toydb_go/grpc/proto"
	"sxg/toydb_go/raft"
	store "sxg/toydb_go/storage/log"
	"sxg/toydb_go/util"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/anypb"
)

type leader struct {
	comm          *commission
	peerNextIndex map[string]uint64
	peerLastIndex map[string]uint64
}

func (lea *leader) becomeFollower(leader string, term uint64) error {
	lea.comm.logger.Infof("discovered new leader %s for term %d", leader, term)
	lea.comm.setTerm(term)
	lea.comm.log().SaveTerm(term, "")
	ins := &raft.InstructionAbort{}
	lea.comm.node.sendIns(ins)
	//转变为follower
	lea.comm.setRole(Follower)
	fol := lea.comm.curRole().(*follower)
	fol.leader = leader
	fol.votedFor = ""
	return nil
}

func (lea *leader) append(p []byte) (index uint64, err error) {
	index, err = lea.comm.log().Append(lea.comm.term(), p)
	if err != nil {
		return
	}
	for _, peer := range lea.comm.node.peers {
		if e := lea.replicate(peer); e != nil {
			lea.comm.logger.Errorf("replicate raft log to peer %s failed: %+v", peer, errors.WithStack(e))
		}
	}
	return
}

func (lea *leader) replicate(peer string) error {
	var ok bool
	var nextIndex, baseIndex, baseTerm uint64
	if nextIndex, ok = lea.peerLastIndex[peer]; !ok {
		return fmt.Errorf("unknown peer %s", peer)
	}
	if nextIndex > 0 {
		baseIndex = nextIndex - 1
	}
	if baseIndex > 0 {
		var base *proto.Entry
		if base, ok = lea.comm.log().Get(baseIndex); !ok {
			return fmt.Errorf("missing base entry at %d", baseIndex)
		}
		baseTerm = base.GetTerm()
	}
	scan := lea.comm.log().Scan(baseIndex, store.Endless)
	start, end := scan.Boundary()
	addr := &proto.Address{
		AddrType: proto.AddrType_PEER,
		Peer:     peer,
	}
	event := &proto.ReplicateEntriesEvent{
		BaseIndex: baseIndex,
		BaseTerm:  baseTerm,
		Entries:   scan.Split(start, end),
	}
	if err := lea.comm.node.send(0, addr, event); err != nil {
		return err
	}
	return nil
}

func (lea *leader) commit() error {
	lastIndex, _ := lea.comm.log().LastIndexAndTerm()
	indexes := []uint64{lastIndex}
	for _, v := range lea.peerLastIndex {
		indexes = append(indexes, v)
	}
	util.Sort(indexes, util.DESC)
	index := indexes[lea.comm.node.quorum()-1]
	if index > lea.comm.log().CommitIndex() {
		//只能commit自己任期内的log
		if entry, ok := lea.comm.log().Get(index); ok && entry.Term == lea.comm.term() {
			oldCommitIndex := lea.comm.log().CommitIndex()
			if err := lea.comm.log().Commit(index); err != nil {
				return err
			}
			scan := lea.comm.log().Scan(oldCommitIndex+1, index)
			for {
				entry, ok := scan.Iterator()
				if !ok {
					break
				}
				ins := &raft.InstructionApply{Entry: entry}
				lea.comm.node.sendIns(ins)
			}
		}
	}
	return nil
}

func (lea *leader) step(cas *raft.Case) {
	msg := cas.Msg
	if msg.Term > lea.comm.term() &&
		msg.From.AddrType == proto.AddrType_PEER {
		if err := lea.becomeFollower(msg.From.Peer, msg.Term); err != nil {
			lea.comm.logger.Errorf("leader convert to follower failed: %+v", errors.WithStack(err))
			return
		}
		fol := lea.comm.curRole().(*follower)
		fol.step(cas)
		return
	}

	event, err := msg.Event.UnmarshalNew()
	if err != nil {
		lea.comm.logger.Errorf("unknown proto message: %+v", errors.WithStack(err))
		return
	}
	switch v := event.(type) {
	case *proto.ConfirmleaderEvent:
		//值得继续分析InstructionVote的目的
		if msg.From.AddrType == proto.AddrType_PEER {
			ins := &raft.InstructionVote{
				Term:  msg.Term,
				Index: v.CommitIndex,
				Addr:  msg.From,
			}
			lea.comm.node.sendIns(ins)
			if !v.HasCommitted {
				if err := lea.replicate(msg.From.Peer); err != nil {
					lea.comm.logger.Errorf("replicating raft log to %s failed: %+v", msg.From.Peer, errors.WithStack(err))
				}
			}
		}
	case *proto.AcceptEntriesEvent:
		if msg.From.AddrType == proto.AddrType_PEER {
			lea.peerLastIndex[msg.From.Peer] = v.LastIndex
			lea.peerNextIndex[msg.From.Peer] = v.LastIndex + 1
		}
		if err := lea.commit(); err != nil {
			lea.comm.logger.Errorf("committing failed: %+v", errors.WithStack(err))
		}
	case *proto.RejectEntriesEvent:
		if msg.From.AddrType == proto.AddrType_PEER {
			if nextIndex := lea.peerNextIndex[msg.From.Peer]; nextIndex > 1 {
				lea.peerNextIndex[msg.From.Peer] = nextIndex - 1
			}
			if err := lea.replicate(msg.From.Peer); err != nil {
				lea.comm.logger.Errorf("replicating raft log to %s failed: %+v", msg.From.Peer, errors.WithStack(err))
			}
		}
	case *proto.SolicitVoteEvent, *proto.GrantVoteEvent: // do nothing
	case *proto.ClientRequestEvent:
		lea.handleClientRequest(cas, v)
	case *proto.ClientResponseEvent:
		if v.ResponseType == proto.RespType_RESP_STATUS {
			status := cas.RespContent.(*proto.Status)
			status.Server = lea.comm.id()
			v.Content, _ = anypb.New(status)
		}
		addr := &proto.Address{AddrType: proto.AddrType_CLIENT}
		if err := lea.comm.node.send(cas.ID, addr, v); err != nil {
			lea.comm.logger.Errorf("sending client response event failed: %+v", errors.WithStack(err))
		}
	case *proto.HeartbeatEvent, *proto.ReplicateEntriesEvent:
		lea.comm.logger.Warnf("received unexpected message: %s", msg.String())
	}
}

func (lea *leader) handleClientRequest(cas *raft.Case, req *proto.ClientRequestEvent) {
	msg := cas.Msg
	if req.RequestType == proto.ReqType_REQ_QUERY {
		ins := &raft.InstructionQuery{
			Id:      req.Id,
			Addr:    msg.From,
			Command: req.Command,
			Term:    lea.comm.term(),
			Index:   lea.comm.log().CommitIndex(),
			Quorum:  lea.comm.node.quorum(),
		}
		lea.comm.node.sendIns(ins)
		if len(lea.comm.node.peers) > 0 {
			addr := &proto.Address{AddrType: proto.AddrType_PEERS}
			index, term := lea.comm.log().CommitIndexAndTerm()
			event := &proto.HeartbeatEvent{
				CommitIndex: index,
				CommitTerm:  term,
			}
			if err := lea.comm.node.send(cas.ID, addr, event); err != nil {
				lea.comm.logger.Errorf("sending heartbeat failed: %+v", errors.WithStack(err))
			}
		}
	} else if req.RequestType == proto.ReqType_REQ_MUTATE {
		index, err := lea.append(req.Command)
		if err != nil {
			lea.comm.logger.Errorf("appending mutate command failed: %+v", errors.WithStack(err))
			return
		}
		ins := &raft.InstructionNotify{
			Id:    req.Id,
			Addr:  msg.From,
			Index: index,
		}
		lea.comm.node.sendIns(ins)
		//单机的情况下，直接commit
		if len(lea.comm.node.peers) == 0 {
			if err := lea.commit(); err != nil {
				lea.comm.logger.Errorf("committing failed: %+v", errors.WithStack(err))
			}
		}
	} else if req.RequestType == proto.ReqType_REQ_STATUS {
		store, size := lea.comm.log().StoreInfo()
		status := &proto.Status{
			Server:      lea.comm.id(),
			Leader:      lea.comm.id(),
			Term:        lea.comm.term(),
			CommitIndex: lea.comm.log().CommitIndex(),
			ApplyIndex:  0,
			Storage:     store,
			StorageSize: uint64(size),
		}
		index, _ := lea.comm.log().LastIndexAndTerm()
		indexes := make(map[string]uint64)
		indexes[lea.comm.id()] = index
		for k, v := range lea.peerLastIndex {
			indexes[k] = v
		}
		status.NodeLastIndex = indexes
		ins := &raft.InstructionStatus{
			Id:     lea.comm.id(),
			Addr:   msg.From,
			Status: status,
		}
		lea.comm.node.sendIns(ins)
	} else {
		lea.comm.logger.Warnf("received unexpected type %d of client's request from %s", req.RequestType, msg.From.Peer)
	}
}

func (lea *leader) reset(peers []string, lastIndex uint64) {
	lea.peerLastIndex = map[string]uint64{}
	lea.peerNextIndex = map[string]uint64{}
	for _, peer := range peers {
		lea.peerNextIndex[peer] = lastIndex + 1
		lea.peerLastIndex[peer] = 0
	}
}
