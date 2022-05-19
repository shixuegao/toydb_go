package raft

import (
	logging "sxg/toydb_go/logging"
	"time"
)

const (
	Follower int = iota
	Candidate
	Leader
)

type Role struct {
	rType  int
	roles  map[int]any
	node   *Node
	logger logging.Logger
}

type follower struct {
	leader   string
	votedFor string
	ticker   *time.Ticker
	timeout  time.Duration
}

type candidate struct {
	votes int
}

type leader struct {
	peerNextIndex map[string]uint64
	peerLastIndex map[string]uint64
}

func (le *leader) append(p []byte) error {
	return nil
}

func newLeader(peers []string, lastIndex uint64) *leader {
	le := &leader{
		peerNextIndex: map[string]uint64{},
		peerLastIndex: map[string]uint64{},
	}
	for _, peer := range peers {
		le.peerNextIndex[peer] = lastIndex
		le.peerLastIndex[peer] = 0
	}
	return le
}

func newRole(node *Node, logger logging.Logger) *Role {
	roles := make(map[int]any)
	roles[Follower] = &follower{}
	roles[Candidate] = &candidate{}
	roles[Leader] = &leader{}
	r := Role{
		rType:  Follower,
		roles:  roles,
		node:   node,
		logger: logger,
	}
	return &r
}

func (r *Role) step(msg *Message) {
	if err := validate(msg); err != nil {
		r.logger.Errorf("ignore invalid message: %s", err.Error())
		return
	}
	switch r.rType {
	case Follower:
		r.stepOnFollower(msg)
	case Candidate:
		r.stepOnCandidate(msg)
	case Leader:
		r.stepOnLeader(msg)
	}
}

func (r *Role) stepOnFollower(msg *Message) {
	follower := r.roles[r.rType].(*follower)
	//如果信息是其他节点发过来的，且对方的Term比本地的大或者本地leader为空，则考虑重置当前节点，并继续处理msg
	if msg.From.addrType == AddrTypePeer &&
		(msg.Term > r.node.term || follower.leader == "") {
		var votedFor string
		if msg.Term > r.node.term {
			r.logger.Infof("discovered new term %d, following leader %s", msg.Term, msg.From.peer)
			r.node.term = msg.Term
			r.node.log.saveTerm(msg.Term, "")
		} else {
			r.logger.Infof("discovered leader %s, following", msg.From.peer)
			votedFor = follower.votedFor
		}
		follower.leader = msg.From.peer
		follower.votedFor = votedFor
		if err := r.node.abortProxied(); err != nil {
			r.logger.Errorf("abort proxied failed: %s", err.Error())
			return
		}
		if err := r.node.forwardQueued(msg.From); err != nil {
			r.logger.Errorf("forward queued failed: %s", err.Error())
			return
		}
		r.step(msg)
		return
	}
	if follower.leader == msg.From.peer {
		follower.ticker.Reset(follower.timeout)
	}
	switch msg.Event.EventType() {
	case EventTypeHeartbeat:
		heartbeat := msg.Event.(*HeartbeatEvent)
		if follower.leader == msg.From.peer {
			has := r.node.log.has(heartbeat.commitIndex, heartbeat.commitTerm)
			if has && heartbeat.commitIndex > r.node.log.commitIndex {
				oldCommittedIndex := r.node.log.commitIndex
				if err := r.node.log.commit(heartbeat.commitIndex); err != nil {
					r.logger.Errorf("Committing raft logger failed: %s", err.Error())
					return
				}
				scan := r.node.log.scan(oldCommittedIndex+1, heartbeat.commitIndex)
				for {
					entry, ok := scan.Iterator()
					if !ok {
						break
					}
					r.node.stateSender <- &InstructionApply{entry: entry}
				}
			}
			r.node.send(msg.From, &ConfirmLeaderEvent{heartbeat.commitIndex, has})
		}
	case EventTypeSolicitVote:
		solicitVote := msg.Event.(*SolicitVoteEvent)
		if follower.votedFor != "" && follower.votedFor != msg.From.peer {
			return
		}
		if solicitVote.lastTerm < r.node.log.lastTerm {
			return
		}
		if solicitVote.lastTerm == r.node.log.lastTerm &&
			solicitVote.lastIndex < r.node.log.lastIndex {
			return
		}
		if msg.From.addrType == AddrTypePeer {
			r.logger.Infof("voting for %s in term %d election", msg.From.peer, r.node.term)
			if err := r.node.send(msg.From, &GrantVoteEvent{}); err != nil {
				r.logger.Errorf("sending grant vote event failed: %s", err.Error())
				return
			}
			r.node.log.saveTerm(r.node.term, msg.From.peer)
			follower.votedFor = msg.From.peer
		}
	case EventTypeReplicateEntries:
		replicateEntries := msg.Event.(*ReplicateEntriesEvent)
		if follower.leader == msg.From.peer {
			if replicateEntries.baseIndex > 0 &&
				!r.node.log.has(replicateEntries.baseIndex, replicateEntries.baseTerm) {
				r.logger.Debugf("rejecting log entries at base %d", replicateEntries.baseIndex)
				r.node.send(msg.From, &RejectEntriesEvent{})
			} else {
				lastIndex, err := r.node.log.splice(replicateEntries.entries)
				if err != nil {
					r.logger.Errorf("splice entries failed: %s", err.Error())
					return
				}
				r.node.send(msg.From, &AcceptEntriesEvent{lastIndex: lastIndex})
			}
		}
	case EventTypeClientRequest:
		clientRequest := msg.Event.(*ClientRequestEvent)
		if follower.leader != "" {
			r.node.proxiedReqs[clientRequest.id] = msg.From
			r.node.send(Addr{AddrTypePeer, follower.leader}, msg.Event)
		} else {
			r.node.queuedReqs = append(r.node.queuedReqs,
				struct {
					event Event
					addr  Addr
				}{msg.Event, msg.From})
		}
	case EventTypeClientResponse:
		clientResponse := msg.Event.(*ClientResponseEvent)
		if clientResponse.response.respType == RespTypeStatus {
			status := clientResponse.response.content.(*Status)
			status.Server = r.node.id
			clientResponse.response.content = status
		}
		delete(r.node.proxiedReqs, clientResponse.id)
		r.node.send(Addr{addrType: AddrTypeClient}, &ClientResponseEvent{clientResponse.id, clientResponse.response})
	case EventTypeGrantVote: //do nothing
	case EventTypeComfirmLeader, EventTypeAcceptEntries, EventTypeRejectEntries:
		r.logger.Warnf("received unexpected message: %s", msg.String())
	}
}

func (r *Role) stepOnCandidate(msg *Message) {
	candidate := r.roles[Candidate].(*candidate)
	if msg.Term > r.node.term &&
		msg.From.addrType == AddrTypePeer {
		r.node.term = msg.Term
		r.node.log.saveTerm(msg.Term, "")
		//转变为follower, 只是转变，暂时不会投票
		r.rType = Follower
		follower := r.roles[Follower].(*follower)
		follower.leader = msg.From.peer
		follower.votedFor = ""
		r.node.abortProxied()
		r.node.forwardQueued(msg.From)
		r.step(msg)
		return
	}
	switch msg.Event.EventType() {
	case EventTypeHeartbeat:
		if msg.Term > r.node.term {
			r.node.term = msg.Term
			r.node.log.saveTerm(msg.Term, "")
			//转变为follower, 只是转变，暂时不会投票
			r.rType = Follower
			follower := r.roles[Follower].(*follower)
			follower.leader = msg.From.peer
			follower.votedFor = ""
			r.node.abortProxied()
			r.node.forwardQueued(msg.From)
			r.step(msg)
			return
		}
	case EventTypeGrantVote:
		r.logger.Debugf("received term %d vote from %s", r.node.term, msg.From.peer)
		candidate.votes += 1
		if candidate.votes >= r.node.quorum() {
			queued := r.node.queuedReqs
			r.node.queuedReqs = []struct {
				event Event
				addr  Addr
			}{}
			//
			r.logger.Infof("won election for term %d, becoming leader", r.node.term)
			r.rType = Leader
			r.roles[Leader] = newLeader(r.node.peers, r.node.log.lastIndex)
			r.node.send(Addr{addrType: AddrTypePeers}, &HeartbeatEvent{
				commitIndex: r.node.log.commitIndex,
				commitTerm:  r.node.log.commitTerm,
			})
			//
			if err := (r.roles[Leader].(*leader)).append([]byte{}); err != nil {
				r.logger.Errorf("append info to peers failed: %s", err.Error())
				return
			}
			if err := r.node.abortProxied(); err != nil {
				r.logger.Errorf("abort proxied failed: %s", err.Error())
				return
			}
			//将客户端请求发送到本地
			for _, wrapper := range queued {
				r.step(&Message{Term: 0, From: wrapper.addr, To: Addr{addrType: AddrTypeLocal}, Event: wrapper.event})
			}
		}
	case EventTypeClientRequest:
		r.node.queuedReqs = append(r.node.queuedReqs, struct {
			event Event
			addr  Addr
		}{msg.Event, msg.From})
	case EventTypeClientResponse:
		clientResponse := msg.Event.(*ClientResponseEvent)
		if clientResponse.response.respType == RespTypeStatus {
			status := clientResponse.response.content.(*Status)
			status.Server = r.node.id
			clientResponse.response.content = status
		}
		delete(r.node.proxiedReqs, clientResponse.id)
		r.node.send(Addr{addrType: AddrTypeClient}, &ClientResponseEvent{clientResponse.id, clientResponse.response})
	case EventTypeSolicitVote:
		{
			//ignore other candidates when we're also campaigning, why?
		}
	case EventTypeComfirmLeader, EventTypeReplicateEntries, EventTypeAcceptEntries, EventTypeRejectEntries:
		r.logger.Warnf("received unexpected message: %s", msg.String())
	}
}

func (r *Role) stepOnLeader(msg *Message) {}
