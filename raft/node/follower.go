package node

import (
	"fmt"
	"sxg/toydb_go/grpc/proto"
	"sxg/toydb_go/raft"

	"google.golang.org/protobuf/types/known/anypb"
)

type follower struct {
	comm     *Commission
	leader   string
	votedFor string
}

func (fol *follower) becomeFollower(leader string, term uint64) error {
	var votedFor string
	if term > fol.comm.term() {
		fol.comm.logger.Infof("discovered new term %d, following leader %s", term, leader)
		//设置新值
		fol.comm.node.term = term
		if err := fol.comm.log().SaveTerm(term, ""); err != nil {
			return err
		}
	} else {
		//这里只会存在term相等的情况, 因为之前的validate已经拦截掉了
		fol.comm.logger.Infof("discovered leader %s, following", leader)
		votedFor = fol.votedFor
	}
	fol.leader = leader
	fol.votedFor = votedFor
	if err := fol.comm.node.abortProxied("follower converts to follower"); err != nil {
		return fmt.Errorf("abort proxied failed: %s", err.Error())
	}
	if err := fol.comm.node.forwardQueued(address(proto.AddrType_PEER, leader)); err != nil {
		return fmt.Errorf("forward queued failed: %s", err.Error())
	}
	return nil
}

func (fol *follower) step(input *raft.InputMsg) {
	//如果信息是其他节点发过来的，且对方的Term比本地的大或者本地leader为空，则考虑重置当前节点，并继续处理msg
	//这里保证follower能够跟随当前最高的Term
	msg := input.Msg
	if msg.From.AddrType == proto.AddrType_PEER &&
		(msg.Term > fol.comm.term() || fol.leader == "") {
		if err := fol.becomeFollower(msg.From.Peer, msg.Term); err != nil {
			fol.comm.logger.Error(err.Error())
		}
		fol.step(input)
		return
	}
	if fol.leader == msg.From.Peer {
		//ticks
	}
	switch v := input.Event.(type) {
	case *proto.HeartbeatEvent:
		if fol.leader != msg.From.Peer {
			return
		}
		has := fol.comm.log().Has(v.CommitIndex, v.CommitTerm)
		curCommitIndex := fol.comm.log().CommitIndex()
		if has && v.CommitIndex > curCommitIndex {
			if err := fol.comm.log().Commit(v.CommitIndex); err != nil {
				fol.comm.logger.Errorf("commit raft log failed: %s", err.Error())
				return
			}
			scan := fol.comm.log().Scan(curCommitIndex+1, v.CommitIndex)
			for {
				entry, ok := scan.Iterator()
				if !ok {
					break
				}
				ins := &raft.InstructionApply{Entry: entry}
				fol.comm.node.sendIns(ins)
			}
		}
		event := &proto.ConfirmleaderEvent{
			CommitIndex:  v.CommitIndex,
			HasCommitted: has,
		}
		if err := fol.comm.node.send(msg.From, event); err != nil {
			fol.comm.logger.Errorf("send confirm leader event failed: %s", err.Error())
		}
	case *proto.SolicitVoteEvent:
		if fol.votedFor != "" && fol.votedFor != msg.From.Peer {
			return
		}
		lastIndex, lastTerm := fol.comm.log().LastIndexAndTerm()
		if v.LastTerm < lastTerm ||
			(v.LastTerm == lastTerm && v.LastIndex < lastIndex) {
			return
		}
		if msg.From.AddrType == proto.AddrType_PEER {
			fol.comm.logger.Infof("voting for %s in term %d election", msg.From.Peer, fol.comm.term())
			event := &proto.GrantVoteEvent{}
			if err := fol.comm.node.send(msg.From, event); err != nil {
				fol.comm.logger.Errorf("send grant vote event failed: %s", err.Error())
				return
			}
			fol.comm.log().SaveTerm(fol.comm.term(), msg.From.Peer)
			fol.votedFor = msg.From.Peer
		}
	case *proto.ReplicateEntriesEvent:
		if fol.leader != msg.From.Peer {
			return
		}
		if v.BaseIndex > 0 && !fol.comm.log().Has(v.BaseIndex, v.BaseTerm) {
			fol.comm.logger.Debugf("reject log entries at base %d", v.BaseIndex)
		} else {
			lastIndex, err := fol.comm.log().Splice(v.Entries)
			if err != nil {
				fol.comm.logger.Errorf("splice entries failed: %s", err.Error())
				return
			}
			event := &proto.AcceptEntriesEvent{
				LastIndex: lastIndex,
			}
			fol.comm.node.send(msg.From, event)
		}
	case *proto.ClientRequestEvent:
		if fol.leader != "" {
			fol.comm.node.registerProxiedRequest(v.Id, msg.From)
			addr := &proto.Address{
				AddrType: proto.AddrType_PEER,
				Peer:     fol.leader,
			}
			if err := fol.comm.node.send(addr, msg.Event); err != nil {
				fol.comm.logger.Errorf("send client request event failed: %s", err.Error())
			}
		} else {
			//将客户端的请求先保存下来，之后再转发出去
			fol.comm.node.storeClientRequests(msg.From, v)
		}
	case *proto.ClientResponseEvent:
		if v.ResponseType == proto.RespType_RESP_STATUS {
			status := input.ClientResponseContent.(*proto.Status)
			status.Server = fol.comm.id()
			v.Content, _ = anypb.New(status)
		}
		fol.comm.node.unregisterProxiedRequest(v.Id)
		addr := &proto.Address{AddrType: proto.AddrType_CLIENT}
		if err := fol.comm.node.send(addr, v); err != nil {
			fol.comm.logger.Errorf("send client response event failed: %s", err.Error())
		}
	case *proto.GrantVoteEvent: // do nothing
	case *proto.ConfirmleaderEvent, *proto.AcceptEntriesEvent, *proto.RejectEntriesEvent:
		fol.comm.logger.Warnf("received unexpected message: %s", msg.String())
	}
}
