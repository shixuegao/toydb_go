package raft

import (
	"context"
	"errors"
	"fmt"
	"sxg/toydb_go/grpc/proto"
	"sxg/toydb_go/logging"
	"sync"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
)

//状态机
type State interface {
	AppliedIndex() uint64
	Mutate(uint64, []byte) ([]byte, error)
	Query([]byte) ([]byte, error)
}

//命令
type Instruction interface {
	String() string
}

type InstructionAbort struct{}

func (ins *InstructionAbort) String() string {
	return "aborting instruction"
}

type InstructionApply struct {
	Entry *proto.Entry
}

func (ins *InstructionApply) String() string {
	return fmt.Sprintf("applying instruction, Entry: %s", ins.Entry.String())
}

type InstructionNotify struct {
	Id    string
	Addr  *proto.Address
	Index uint64
}

func (ins *InstructionNotify) String() string {
	return fmt.Sprintf("notifying instruction, Id: %s, addr: %s, index: %d", ins.Id, ins.Addr.String(), ins.Index)
}

type InstructionQuery struct {
	Id      string
	Addr    *proto.Address
	Command []byte
	Term    uint64
	Index   uint64
	Quorum  int
}

func (ins *InstructionQuery) String() string {
	return fmt.Sprintf("querying instruction, Id: %s, addr: %s, command: %s, term: %d, index: %d, quorum: %d",
		ins.Id, ins.Addr.String(), string(ins.Command), ins.Term, ins.Index, ins.Quorum)
}

type InstructionStatus struct {
	Id     string
	Addr   *proto.Address
	Status *proto.Status
}

func (ins *InstructionStatus) String() string {
	return fmt.Sprintf("status instruction, Id: %s, addr: %s, status: %s", ins.Id, ins.Addr.String(), ins.Status.String())
}

type InstructionVote struct {
	Term  uint64
	Index uint64
	Addr  *proto.Address
}

func (ins *InstructionVote) String() string {
	return fmt.Sprintf("voting instruction, term: %d, index: %d, addr: %s", ins.Term, ins.Index, ins.Addr.String())
}

//driver
type wrapper struct {
	addr *proto.Address
	id   string
}

type Driver struct {
	logger       logging.Logger
	context      context.Context
	insReceiver  <-chan Instruction
	msgSender    chan<- *proto.Message
	notify       map[string]*wrapper
	queries      *queries
	appliedIndex uint64
}

func NewDriver(context context.Context, insReceiver <-chan Instruction,
	msgSender chan<- *proto.Message, logger logging.Logger) *Driver {
	return &Driver{
		logger:       logger,
		context:      context,
		insReceiver:  insReceiver,
		msgSender:    msgSender,
		notify:       map[string]*wrapper{},
		queries:      newQueries(),
		appliedIndex: 0,
	}
}

func (dri *Driver) Drive(state State, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		dri.logger.Debugf("starting state machine driver")
		for {
			bk := false
			select {
			case <-dri.context.Done():
				bk = true
			case ins := <-dri.insReceiver:
				if err := dri.execute(ins, state); err != nil {
					dri.logger.Errorf("excuting of driver exception: %s", err.Error())
					return
				}
			}
			if bk {
				break
			}
		}
		dri.logger.Debugf("stopping state machine driver")
	}()
}

func (dri *Driver) execute(ins Instruction, state State) error {
	dri.logger.Debugf("executing %s", ins.String())
	switch v := ins.(type) {
	case *InstructionAbort:
		dri.notifyAbort("")
		dri.queryAbort("")
	case *InstructionApply:
		if len(v.Entry.Command) > 0 {
			dri.logger.Debugf("applying state machine command %d: %s", v.Entry.Index, string(v.Entry.Command))
			result, err := state.Mutate(v.Entry.Index, v.Entry.Command)
			if err != nil {
				return err
			}
			dri.notifyApplied(result)
		}
		dri.appliedIndex = v.Entry.Index
		return dri.queryExecute(state)
	case *InstructionNotify:
		if v.Index > state.AppliedIndex() {
			wra := &wrapper{addr: v.Addr, id: v.Id}
			dri.notify[wra.id] = wra
		} else {
			abort, _ := anypb.New(&proto.Abort{Reason: "Index is applied"})
			event := &proto.ClientResponseEvent{Id: v.Id, ResponseType: proto.RespType_RESP_ABORT, Content: abort}
			if err := dri.sendToClient(v.Addr, event); err != nil {
				dri.logger.Errorf("sending abort to client failed: %s", err.Error())
			}
		}
	case *InstructionQuery:
		dri.queries.addQuery(v)
	case *InstructionStatus:
		v.Status.ApplyIndex = state.AppliedIndex()
		status, _ := anypb.New(v.Status)
		event := &proto.ClientResponseEvent{Id: v.Id, ResponseType: proto.RespType_RESP_STATUS, Content: status}
		if err := dri.sendToClient(v.Addr, event); err != nil {
			dri.logger.Errorf("sending status to client failed: %s", err.Error())
		}
	case *InstructionVote:
		dri.queryVote(v.Term, v.Index, v.Addr)
		if err := dri.queryExecute(state); err != nil {
			return err
		}
	default:
		dri.logger.Warnf("unknown type of instruction: %s", v.String())
	}
	return nil
}

func (dri *Driver) queryExecute(state State) error {
	for _, query := range dri.queryReady(dri.appliedIndex) {
		dri.logger.Debugf("Executing query %v", string(query.command))
		result, err := state.Query(query.command)
		if err != nil {
			return err
		}
		content, _ := anypb.New(&proto.State{Result: result})
		response := &proto.ClientResponseEvent{Id: query.id, ResponseType: proto.RespType_RESP_STATE, Content: content}
		if err := dri.sendToClient(query.addr, response); err != nil {
			dri.logger.Errorf("sending state to client failed: %s", err.Error())
		}
	}
	return nil
}

func (dri *Driver) queryReady(appliedIndex uint64) []*query {
	ready := []*query{}
	empty := []uint64{}
	dri.queries.iterRange(0, appliedIndex, func(tree *queryTree) {
		readyIds := []string{}
		tree.iter(func(q *query) {
			if len(q.votes) >= q.quorum {
				readyIds = append(readyIds, q.id)
			}
		})
		for _, id := range readyIds {
			query := tree.del(id)
			ready = append(ready, query)
		}
		if tree.size() == 0 {
			empty = append(empty, tree.index)
		}
	})
	for _, index := range empty {
		dri.queries.del(index)
	}
	return ready
}

func (dri *Driver) queryVote(term uint64, index uint64, addr *proto.Address) {
	dri.queries.iterRange(0, index, func(tree *queryTree) {
		tree.iter(func(q *query) {
			if term >= q.term {
				q.votes[addr] = struct{}{}
			}
		})
	})
}

func (dri *Driver) notifyAbort(reason string) {
	notify := dri.notify
	dri.notify = map[string]*wrapper{}
	abort, _ := anypb.New(&proto.Abort{Reason: reason})
	for _, wrapper := range notify {
		event := &proto.ClientResponseEvent{
			Id:           wrapper.id,
			ResponseType: proto.RespType_RESP_ABORT,
			Content:      abort,
		}
		if err := dri.sendToClient(wrapper.addr, event); err != nil {
			dri.logger.Errorf("sending abort to client failed: %s", err.Error())
		}
	}
}

func (dri *Driver) queryAbort(reason string) {
	queries := dri.queries
	dri.queries = newQueries()
	abort, _ := anypb.New(&proto.Abort{Reason: reason})
	queries.iter(func(tree *queryTree) {
		tree.iter(func(q *query) {
			event := &proto.ClientResponseEvent{
				Id:           q.id,
				ResponseType: proto.RespType_RESP_ABORT,
				Content:      abort,
			}
			if err := dri.sendToClient(q.addr, event); err != nil {
				dri.logger.Errorf("sending abort to client failed: %s", err.Error())
			}
		})
	})
}

func (dri *Driver) notifyApplied(result []byte) {

}

//发送信息到client，因此term可以为0
func (dri *Driver) sendToClient(to *proto.Address, event any) (err error) {
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
		Term:  0,
		From:  from,
		To:    to,
		Event: eve,
	}
	dri.msgSender <- msg
	return nil
}
