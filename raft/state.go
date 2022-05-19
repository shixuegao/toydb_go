package raft

import (
	store "sxg/toydb_go/storage/log"
)

type Instruction any

type InstructionApply struct {
	entry store.Entry
}
