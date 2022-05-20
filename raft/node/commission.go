package node

import (
	"errors"
	"fmt"
	"sxg/toydb_go/grpc/proto"
	logging "sxg/toydb_go/logging"
	"sxg/toydb_go/raft"
)

type roleType int32

const (
	Follower roleType = iota
	Candidate
	Leader
)

type Commission struct {
	node     *Node
	logger   logging.Logger
	roleType roleType
	roles    map[roleType]role
}

type role interface {
	step(*proto.Message)
}

func NewCommission(node *Node, logger logging.Logger) *Commission {
	comm := &Commission{
		node:   node,
		logger: logger,
	}
	comm.roles = make(map[roleType]role)
	comm.roleType = Follower
	comm.roles[Follower] = &follower{comm: comm}
	comm.roles[Candidate] = &candidate{comm: comm}
	comm.roles[Leader] = &leader{comm: comm}
	return comm
}

func (comm *Commission) curRole() role {
	return comm.roles[comm.roleType]
}

func (comm *Commission) setRole(roleType roleType) {
	comm.roleType = roleType
}

func (comm *Commission) term() uint64 {
	return comm.node.term
}

func (comm *Commission) setTerm(term uint64) {
	comm.node.term = term
}

func (comm *Commission) id() string {
	return comm.node.id
}

func (comm *Commission) log() *raft.Log {
	return comm.node.log
}

func (comm *Commission) step(msg *proto.Message) {
	if err := comm.validate(msg); err != nil {
		comm.logger.Errorf("ignore invalid message: %s", err.Error())
		return
	}
	comm.roles[comm.roleType].step(msg)
}

func (comm *Commission) validate(msg *proto.Message) error {
	switch msg.From.AddrType {
	case proto.AddrType_PEERS:
		return errors.New("message from broadcast address")
	case proto.AddrType_LOCAL:
		return errors.New("message from local node")
	case proto.AddrType_CLIENT:
		if ok, err := isMsgSpecifiedEvents([]proto.EventType{proto.EventType_CLIENT_REQUEST}, msg); err != nil {
			return err
		} else if !ok {
			return errors.New("non-request message from client")
		}
	}
	//此处将任期比自己小的信息(除客户端的请求或响应外)给屏蔽掉
	if msg.GetTerm() < comm.term() {
		types := []proto.EventType{proto.EventType_CLIENT_REQUEST, proto.EventType_CLIENT_RESPONSE}
		if ok, err := isMsgSpecifiedEvents(types, msg); err != nil {
			return err
		} else if !ok {
			return errors.New("message from past term")
		}
	}
	switch msg.To.AddrType {
	case proto.AddrType_PEER:
		if msg.To.Peer != comm.id() {
			return fmt.Errorf("received message from other node %s", msg.To.Peer)
		}
	case proto.AddrType_CLIENT:
		return errors.New("received message for client")
	}
	return nil
}

func isMsgSpecifiedEvents(eventTypes []proto.EventType, msg *proto.Message) (bool, error) {
	event, err := msg.Event.UnmarshalNew()
	if err != nil {
		return false, err
	}
	var ok bool
	for _, eventType := range eventTypes {
		switch eventType {
		case proto.EventType_HEARTBEAT:
			_, ok = event.(*proto.HeartbeatEvent)
		case proto.EventType_CONFIRM_LEADER:
			_, ok = event.(*proto.ConfirmleaderEvent)
		case proto.EventType_SOLICIT_VOTE:
			_, ok = event.(*proto.SolicitVoteEvent)
		case proto.EventType_GRANT_VOTE:
			_, ok = event.(*proto.GrantVoteEvent)
		case proto.EventType_REPLICATE_ENTRIES:
			_, ok = event.(*proto.ReplicateEntriesEvent)
		case proto.EventType_ACCEPT_ENTRIES:
			_, ok = event.(*proto.AcceptEntriesEvent)
		case proto.EventType_REJECT_ENTRIES:
			_, ok = event.(*proto.RejectEntriesEvent)
		case proto.EventType_CLIENT_REQUEST:
			_, ok = event.(*proto.ClientRequestEvent)
		case proto.EventType_CLIENT_RESPONSE:
			_, ok = event.(*proto.ClientResponseEvent)
		default:
			return false, errors.New("unknown event type of message")
		}
		if ok {
			return true, nil
		}
	}
	return ok, nil
}

func address(addrType proto.AddrType, peer string) *proto.Address {
	return &proto.Address{AddrType: addrType, Peer: peer}
}
