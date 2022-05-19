package raft

import (
	"encoding/json"
	"errors"
	"strconv"
	store "sxg/toydb_go/storage/log"
)

// address type
const (
	AddrTypePeers int = iota
	AddrTypePeer
	AddrTypeLocal
	AddrTypeClient
)

type Addr struct {
	addrType int
	peer     string
}

func NewAddress(addrType int, peer string) Addr {
	switch addrType {
	case AddrTypePeer:
		if peer == "" {
			panic("peer cannot be empty on the type of peer")
		}
		return Addr{addrType, peer}
	case AddrTypePeers, AddrTypeLocal, AddrTypeClient:
		return Addr{addrType: addrType}
	default:
		panic("unknown address type: " + strconv.Itoa(addrType))
	}
}

const (
	EventTypeHeartbeat int = iota
	EventTypeComfirmLeader
	EventTypeSolicitVote
	EventTypeGrantVote
	EventTypeReplicateEntries
	EventTypeAcceptEntries
	EventTypeRejectEntries
	EventTypeClientRequest
	EventTypeClientResponse
)

// event
type Event interface {
	EventType() int
}
type HeartbeatEvent struct {
	commitIndex uint64
	commitTerm  uint64
}

func (h *HeartbeatEvent) EventType() int {
	return EventTypeHeartbeat
}

type ConfirmLeaderEvent struct {
	commitIndex  uint64
	hasCommitted bool
}

func (c *ConfirmLeaderEvent) EventType() int {
	return EventTypeComfirmLeader
}

type SolicitVoteEvent struct {
	lastIndex uint64
	lastTerm  uint64
}

func (s *SolicitVoteEvent) EventType() int {
	return EventTypeSolicitVote
}

type GrantVoteEvent struct{}

func (g *GrantVoteEvent) EventType() int {
	return EventTypeGrantVote
}

type ReplicateEntriesEvent struct {
	baseIndex uint64
	baseTerm  uint64
	entries   []store.Entry
}

func (r *ReplicateEntriesEvent) EventType() int {
	return EventTypeReplicateEntries
}

type RejectEntriesEvent struct{}

func (r *RejectEntriesEvent) EventType() int {
	return EventTypeRejectEntries
}

type AcceptEntriesEvent struct {
	lastIndex uint64
}

func (a *AcceptEntriesEvent) EventType() int {
	return EventTypeAcceptEntries
}

type ClientRequestEvent struct {
	id      string
	request Request
}

func (c *ClientRequestEvent) EventType() int {
	return EventTypeClientRequest
}

type ClientResponseEvent struct {
	id       string
	response Response
}

func (c *ClientResponseEvent) EventType() int {
	return EventTypeClientResponse
}

// message
type Message struct {
	Term  uint64 `json:"term"`
	From  Addr   `json:"from"`
	To    Addr   `json:"to"`
	Event Event  `json:"event"`
}

func (msg *Message) String() string {
	c, _ := json.Marshal(msg)
	return string(c)
}

const (
	ReqTypeQuery int = iota
	ReqTypeMutate
	ReqTypeStatus
)

type Request struct {
	reqType int
	command []byte
}

const (
	RespTypeState int = iota
	RespTypeStatus
)

type Response struct {
	respType int
	content  any
}

type Status struct {
	Server        string
	Leader        string
	Term          uint64
	NodeLastIndex map[string]uint64
	CommitIndex   uint64
	ApplyIndex    uint64
	Storage       string
	StorageSize   uint64
}

func validate(msg *Message) error {
	switch msg.From.addrType {
	case AddrTypePeers:
		return errors.New("message from broadcast address")
	case AddrTypeLocal:
		return errors.New("message from local node")
	case AddrTypeClient:
		if msg.Event.EventType() != EventTypeClientRequest {
			return errors.New("non-request message form client")
		}
	}
	return nil
}
