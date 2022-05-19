package raft

import (
	"errors"
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

func (log *Log) saveTerm(term uint64, votedFor string) error {
	return nil
}

func (log *Log) append(term uint64, command []byte) (index uint64, _ error) {
	index = log.lastIndex + 1
	entry := log.producer(index, term, command)
	if _, err := log.store.Append(entry); err != nil {
		return 0, err
	}
	log.lastIndex = entry.Index()
	log.lastTerm = entry.Term()
	return
}

func (log *Log) batchAppend(entries []store.Entry) (index uint64, _ error) {
	for _, entry := range entries {
		if _, err := log.append(entry.Term(), entry.Command()); err != nil {
			return 0, err
		}
	}
	return log.lastIndex, nil
}

func (log *Log) commit(index uint64) error {
	if err := log.store.Commit(index); err != nil {
		return err
	}
	entry, _ := log.store.Get(index)
	log.commitIndex = entry.Index()
	log.commitTerm = entry.Term()
	return nil
}

func (log *Log) truncate(index uint64) error {
	if _, err := log.store.Truncate(index); err != nil {
		return err
	}
	last, _ := log.store.GetLast()
	log.lastIndex = last.Index()
	log.lastTerm = last.Term()
	return nil
}

func (log *Log) get(index uint64) (store.Entry, bool) {
	return log.store.Get(index)
}

func (log *Log) has(index, term uint64) bool {
	entry, ok := log.store.Get(index)
	return ok && entry.Term() == term
}

func (log *Log) scan(start, end uint64) store.Scan {
	return log.store.Scan(start, end)
}

//日志必须连续，且第一个日志的index最大为lastIndex+1
func (log *Log) splice(entries []store.Entry) (uint64, error) {
	var l int
	if l = len(entries); l == 0 {
		return log.lastIndex, nil
	}
	if entries[0].Index() > log.lastIndex+1 {
		return 0, errors.New("spliced entries cannot begin past last index")
	}
	former := entries[0].Index()
	for i := 1; i < len(entries); i++ {
		if latter := entries[i].Index(); former+1 != latter {
			return 0, errors.New("spliced entries must be contiguous")
		} else {
			former = latter
		}
	}
	//裁剪日志
	first := entries[0]
	if first.Index() == log.lastIndex+1 {
		return log.batchAppend(entries)
	}
	for _, entry := range entries {
		if current, ok := log.store.Get(entry.Index()); !ok || current.Term() == entry.Term() {
			continue
		}
		if err := log.truncate(entry.Index() - 1); err != nil {
			return 0, err
		}
		log.append(entry.Term(), entry.Command())
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
		log.commitIndex = entry.Index()
		log.commitTerm = entry.Index()
	}
	//last index, term
	if entry, ok := store.GetLast(); !ok {
		return nil, errors.New("last entry not found")
	} else {
		log.lastIndex = entry.Index()
		log.lastTerm = entry.Term()
	}
	return
}
