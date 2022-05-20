package proto_test

import (
	"fmt"
	"sxg/toydb_go/grpc/proto"
	"testing"

	"google.golang.org/protobuf/types/known/anypb"
)

func Test(t *testing.T) {
	event := &proto.GrantVoteEvent{EventType: proto.EventType_GRANT_VOTE}
	anyyy, _ := anypb.New(event)
	msg := &proto.Message{
		Term:  0,
		From:  nil,
		To:    nil,
		Event: anyyy,
	}

	eve, _ := msg.Event.UnmarshalNew()
	_, ok := eve.(*proto.GrantVoteEvent)
	_, ok2 := eve.(*proto.AcceptEntriesEvent)
	fmt.Println(ok, ok2)
}
