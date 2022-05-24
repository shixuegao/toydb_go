package raft

import (
	"fmt"
	"sxg/toydb_go/grpc/proto"

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
