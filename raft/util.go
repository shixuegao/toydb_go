package raft

import (
	"fmt"
	"sxg/toydb_go/grpc/proto"

	"github.com/google/btree"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type InputMsg struct {
	Msg                   *proto.Message
	Event                 protoreflect.ProtoMessage //已经解码好的proto结构体，可以直接用
	ClientResponseContent protoreflect.ProtoMessage //已经解码好的客户端响应proto结构体，可以直接用
}

func VerifyAndWrap(msg *proto.Message) (*InputMsg, error) {
	event, err := msg.Event.UnmarshalNew()
	if err != nil {
		return nil, err
	}
	var matched bool = true
	var content protoreflect.ProtoMessage = nil
	switch msg.EventType {
	case proto.EventType_HEARTBEAT:
		_, matched = event.(*proto.HeartbeatEvent)
	case proto.EventType_CONFIRM_LEADER:
		_, matched = event.(*proto.ConfirmleaderEvent)
	case proto.EventType_SOLICIT_VOTE:
		_, matched = event.(*proto.SolicitVoteEvent)
	case proto.EventType_GRANT_VOTE:
		_, matched = event.(*proto.GrantVoteEvent)
	case proto.EventType_REPLICATE_ENTRIES:
		_, matched = event.(*proto.ReplicateEntriesEvent)
	case proto.EventType_ACCEPT_ENTRIES:
		_, matched = event.(*proto.AcceptEntriesEvent)
	case proto.EventType_REJECT_ENTRIES:
		_, matched = event.(*proto.RejectEntriesEvent)
	case proto.EventType_CLIENT_REQUEST:
		_, matched = event.(*proto.ClientRequestEvent)
	case proto.EventType_CLIENT_RESPONSE:
		if vv, ok := event.(*proto.ClientResponseEvent); ok {
			if content, err = VerifyClientResponse(vv); err != nil {
				return nil, err
			}
			matched = ok
		}
	default:
		return nil, fmt.Errorf("unknown event type of message %v", msg.EventType)
	}
	if !matched {
		return nil, fmt.Errorf("unexpected event, need type %v, the real is %v", msg.EventType, event)
	}
	return &InputMsg{msg, event, content}, nil
}

func VerifyClientResponse(event *proto.ClientResponseEvent) (v protoreflect.ProtoMessage, err error) {
	if v, err = event.Content.UnmarshalNew(); err != nil {
		return
	}
	var matched bool = true
	switch event.ResponseType {
	case proto.RespType_RESP_STATUS:
		if vv, ok := v.(*proto.Status); ok {
			v = vv
			matched = ok
		}
	case proto.RespType_RESP_STATE:
		if vv, ok := v.(*proto.State); ok {
			v = vv
			matched = ok
		}
	case proto.RespType_RESP_ABORT:
		if vv, ok := v.(*proto.Abort); ok {
			v = vv
			matched = ok
		}
	default:
		return nil, fmt.Errorf("unknown client response type %v", event.ResponseType)
	}
	if !matched {
		return nil, fmt.Errorf("unexpected client response content, need type %v, the real is %v", event.ResponseType, v)
	}
	return
}

// query
type query struct {
	id      string
	term    uint64
	addr    *proto.Address
	command []byte
	quorum  int
	votes   map[*proto.Address]struct{}
}

func (q *query) Less(than btree.Item) bool {
	th := than.(*query)
	return q.id < th.id
}

type handleQuery func(q *query)
type handleQueryTree func(tree *queryTree)

// query tree
type queryTree struct {
	index uint64
	tree  *btree.BTree // tree<string, query>
}

func (q *queryTree) Less(than btree.Item) bool {
	th := than.(*queryTree)
	return q.index < th.index
}

func (q *queryTree) iter(handle handleQuery) {
	if q.tree.Len() == 0 {
		return
	}
	q.tree.Ascend(func(i btree.Item) bool {
		q := i.(*query)
		handle(q)
		return true
	})
}

func (tree *queryTree) del(id string) *query {
	if item := tree.tree.Delete(&query{id: id}); item != nil {
		return item.(*query)
	}
	return nil
}

func (tree *queryTree) size() int {
	return tree.tree.Len()
}

// queries
type queries struct {
	tree *btree.BTree //tree<uint64, tree<string, query>>
}

func newQueries() *queries {
	return &queries{tree: btree.New(2)}
}

func (qs *queries) del(index uint64) *queryTree {
	if item := qs.tree.Delete(&queryTree{index: index}); item != nil {
		return item.(*queryTree)
	}
	return nil
}

func (qs *queries) addQuery(ins *InstructionQuery) {
	var qTree *queryTree
	if !qs.tree.Has(&queryTree{index: ins.Index}) {
		qTree = &queryTree{index: ins.Index}
		qs.tree.ReplaceOrInsert(qTree)
	} else {
		qTree = qs.tree.Get(&queryTree{index: ins.Index}).(*queryTree)
	}
	q := &query{
		id:      ins.Id,
		term:    ins.Term,
		addr:    ins.Addr,
		command: ins.Command,
		quorum:  ins.Quorum,
		votes:   make(map[*proto.Address]struct{}),
	}
	qTree.tree.ReplaceOrInsert(q)
}

func (qs *queries) iter(handle handleQueryTree) {
	if qs.tree.Len() == 0 {
		return
	}
	qs.tree.Ascend(func(i btree.Item) bool {
		qTree := i.(*queryTree)
		handle(qTree)
		return true
	})
}

func (qs *queries) iterRange(start, end uint64, handle handleQueryTree) {
	if start > end {
		panic("start index is bigger than end index")
	}
	if qs.tree.Len() == 0 {
		return
	}
	qs.tree.AscendRange(&queryTree{index: start}, &queryTree{index: end}, func(i btree.Item) bool {
		qTree := i.(*queryTree)
		handle(qTree)
		return true
	})
}
