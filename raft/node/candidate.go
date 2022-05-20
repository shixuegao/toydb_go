package node

import "sxg/toydb_go/grpc/proto"

type candidate struct {
	comm  *Commission
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
	if err := can.comm.node.abortProxied(); err != nil {
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
		EventType:   proto.EventType_HEARTBEAT,
		CommitIndex: commitIndex,
		CommitTerm:  commitTerm,
	}
	if err := lea.comm.node.send(addr, event); err != nil {
		return err
	}
	if err := lea.append([]byte{}); err != nil {
		return err
	}
	if err := lea.comm.node.abortProxied(); err != nil {
		return err
	}
	return nil
}

func (can *candidate) step(msg *proto.Message) {
	if msg.Term > can.comm.term() &&
		msg.From.AddrType == proto.AddrType_PEER {
		if err := can.becomeFollower(msg.From.Peer, msg.Term); err != nil {
			can.comm.logger.Errorf("candidate convert to follower failed: %s", err.Error())
			return
		}
		can.comm.step(msg)
		return
	}
	event, err := msg.Event.UnmarshalNew()
	if err != nil {
		can.comm.logger.Errorf("unknow proto message: %s", err.Error())
		return
	}
	switch v := event.(type) {
	case *proto.HeartbeatEvent:
		if msg.From.AddrType == proto.AddrType_PEER {
			if err := can.becomeFollower(msg.From.Peer, msg.Term); err != nil {
				can.comm.logger.Errorf("candidate convert to follower failed: %s", err.Error())
				return
			}
			can.comm.step(msg)
			return
		}
	case *proto.GrantVoteEvent:
		can.comm.logger.Debugf("received term %d vote from %s", can.comm.term(), msg.From.Peer)
		can.votes++
		if can.votes > can.comm.node.quorum() {
			reqsQueue := can.comm.node.resetClientRequestsQueue()
			can.comm.logger.Infof("won election for term %d, become leader", can.comm.term())
			if err := can.becomeLeader(); err != nil {
				can.comm.logger.Errorf("candidate convert to leader failed: %s", err.Error())
				return
			}
			//将客户端请求转发到本地
			for _, req := range reqsQueue {
				msg := &proto.Message{
					Term:  0,
					From:  req.addr,
					To:    &proto.Address{AddrType: proto.AddrType_LOCAL},
					Event: req.event,
				}
				can.step(msg)
			}
		}
	case *proto.ClientRequestEvent:
		can.comm.node.storeClientRequestsQueue(msg.From, msg.Event)
	case *proto.ClientResponseEvent:
		if v.RequestType == proto.ReqType_STATUS {
			v.Status.Server = can.comm.id()
		}
		can.comm.node.unregisterProxiedRequest(v.Id)
		addr := &proto.Address{AddrType: proto.AddrType_CLIENT}
		if err := can.comm.node.send(addr, v); err != nil {
			can.comm.logger.Errorf("send client response event failed: %s", err.Error())
		}
	case *proto.SolicitVoteEvent: //do nothing
	case *proto.ConfirmleaderEvent, *proto.ReplicateEntriesEvent, *proto.AcceptEntriesEvent, *proto.RejectEntriesEvent:
		can.comm.logger.Warnf("received unexpected message: %s", msg.String())
	}
}
