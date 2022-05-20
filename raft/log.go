package raft

import (
	"errors"
	"sxg/toydb_go/grpc/proto"
	store "sxg/toydb_go/storage/log"
)

type Log struct {
	producer    store.EntryProducer
	store       store.Store
	lastIndex   uint64
	lastTerm    uint64
	commitIndex uint64
	commitTerm  uint64
}

func (log *Log) SaveTerm(term uint64, votedFor string) error {
	return nil
}

func (log *Log) Append(term uint64, command []byte) (index uint64, _ error) {
	index = log.lastIndex + 1
	entry := log.producer(index, term, command)
	if _, err := log.store.Append(entry); err != nil {
		return 0, err
	}
	log.lastIndex = entry.GetIndex()
	log.lastTerm = entry.GetTerm()
	return
}

func (log *Log) CommitIndex() uint64 {
	return log.commitIndex
}

func (log *Log) CommitIndexAndTerm() (uint64, uint64) {
	return log.commitIndex, log.commitTerm
}

func (log *Log) LastIndexAndTerm() (uint64, uint64) {
	return log.lastIndex, log.lastTerm
}

func (log *Log) BatchAppend(entries []*proto.Entry) (index uint64, _ error) {
	for _, entry := range entries {
		if _, err := log.Append(entry.GetTerm(), entry.GetCommand()); err != nil {
			return 0, err
		}
	}
	return log.lastIndex, nil
}

func (log *Log) Commit(index uint64) error {
	if err := log.store.Commit(index); err != nil {
		return err
	}
	entry, _ := log.store.Get(index)
	log.commitIndex = entry.GetIndex()
	log.commitTerm = entry.GetTerm()
	return nil
}

func (log *Log) Truncate(index uint64) error {
	if _, err := log.store.Truncate(index); err != nil {
		return err
	}
	last, _ := log.store.GetLast()
	log.lastIndex = last.GetIndex()
	log.lastTerm = last.GetTerm()
	return nil
}

func (log *Log) Get(index uint64) (proto.Entry, bool) {
	return log.store.Get(index)
}

func (log *Log) Has(index, term uint64) bool {
	entry, ok := log.store.Get(index)
	return ok && entry.GetTerm() == term
}

func (log *Log) Scan(start, end uint64) store.Scan {
	return log.store.Scan(start, end)
}

//日志必须连续，且第一个日志的index最大为lastIndex+1
func (log *Log) Splice(entries []*proto.Entry) (uint64, error) {
	var l int
	if l = len(entries); l == 0 {
		return log.lastIndex, nil
	}
	if entries[0].GetIndex() > log.lastIndex+1 {
		return 0, errors.New("spliced entries cannot begin past last index")
	}
	former := entries[0].GetIndex()
	for i := 1; i < len(entries); i++ {
		if latter := entries[i].GetIndex(); former+1 != latter {
			return 0, errors.New("spliced entries must be contiguous")
		} else {
			former = latter
		}
	}
	//裁剪日志
	first := entries[0]
	if first.GetIndex() == log.lastIndex+1 {
		return log.BatchAppend(entries)
	}
	for _, entry := range entries {
		if current, ok := log.store.Get(entry.GetIndex()); !ok || current.GetTerm() == entry.GetTerm() {
			continue
		}
		if err := log.Truncate(entry.GetIndex() - 1); err != nil {
			return 0, err
		}
		log.Append(entry.GetTerm(), entry.GetCommand())
	}
	return log.lastIndex, nil
}

func NewLog(producer store.EntryProducer, store store.Store) (log *Log, _ error) {
	log = &Log{producer: producer, store: store}
	//committed index, term
	if ci := store.Commited(); ci == 0 {
		log.commitIndex = 0
		log.commitTerm = 0
	} else if entry, ok := store.Get(ci); !ok {
		return nil, errors.New("committed entry not found")
	} else {
		log.commitIndex = entry.GetIndex()
		log.commitTerm = entry.GetTerm()
	}
	//last index, term
	if entry, ok := store.GetLast(); !ok {
		return nil, errors.New("last entry not found")
	} else {
		log.lastIndex = entry.GetIndex()
		log.lastTerm = entry.GetTerm()
	}
	return
}
