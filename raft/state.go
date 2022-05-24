package raft

import (
	"context"
	"errors"
	"fmt"
	"sxg/toydb_go/grpc/proto"
	"sxg/toydb_go/logging"
	"sync"

	"github.com/google/btree"
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
	logger      logging.Logger
	context     context.Context
	insReceiver <-chan Instruction
	msgSender   chan<- *proto.Message
	notify      map[string]*wrapper
	queries     *btree.BTree
}

func NewDriver(context context.Context, insReceiver <-chan Instruction,
	msgSender chan<- *proto.Message, logger logging.Logger) *Driver {
	return &Driver{
		logger:      logger,
		context:     context,
		insReceiver: insReceiver,
		msgSender:   msgSender,
		notify:      map[string]*wrapper{},
		queries:     btree.New(2),
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
		if err := dri.notifyAbort(""); err != nil {
			return err
		}
		if err := dri.queryAbort(""); err != nil {
			return err
		}
	case *InstructionApply:
	case *InstructionNotify:
	case *InstructionQuery:
	case *InstructionStatus:
	case *InstructionVote:
	default:
		dri.logger.Warnf("unknown type of instruction: %s", v.String())
	}
	return nil
}

func (dri *Driver) notifyAbort(reason string) error {
	notify := dri.notify
	dri.notify = map[string]*wrapper{}
	abort, _ := anypb.New(&proto.Abort{Reason: reason})
	for _, wrapper := range notify {
		event := &proto.ClientResponseEvent{
			Id:      wrapper.id,
			Content: abort,
		}
		if err := dri.send(wrapper.addr, event); err != nil {
			dri.logger.Errorf("sending abort to client failed: %s", err.Error())
		}
	}
	return nil
}

func (dri *Driver) queryAbort(reason string) error {
	return nil
}

//发送信息到client，因此term可以为0
func (dri *Driver) send(to *proto.Address, event any) (err error) {
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
