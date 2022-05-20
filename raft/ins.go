package raft

import (
	"sxg/toydb_go/grpc/proto"
)

type Instruction any

type InstructionApply struct {
	Entry *proto.Entry
}
