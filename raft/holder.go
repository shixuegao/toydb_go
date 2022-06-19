package raft

import (
	"context"
	"sxg/toydb_go/grpc/proto"
	"sync"
	"sync/atomic"
	"time"
)

type waiter struct {
	caseId uint64
	c      chan *Case
	e      chan struct{}
	create time.Time
	remove int32
}

type CaseHolder struct {
	id      uint64
	m       *sync.Map
	wg      *sync.WaitGroup
	context context.Context
}

func NewCaseKeeper(context context.Context, wg *sync.WaitGroup) *CaseHolder {
	return &CaseHolder{
		id:      0,
		m:       new(sync.Map),
		wg:      wg,
		context: context}
}

// id不能为0
func (ck *CaseHolder) incrId() uint64 {
	var res uint64
	for {
		if res = atomic.AddUint64(&ck.id, 1); res != 0 {
			break
		}
	}
	return res
}

// 等待信息被处理
func (ck *CaseHolder) Await(msg *proto.Message, f func(cas *Case)) *Case {
	cas := &Case{
		ID:      ck.incrId(),
		Msg:     msg,
		Timeout: false}
	wai := &waiter{
		caseId: cas.ID,
		c:      make(chan *Case, 1),
		e:      make(chan struct{}),
		create: time.Now(),
		remove: 0}
	ck.m.Store(cas.ID, wai)
	ck.wg.Add(1)
	go func() {
		defer ck.wg.Done()
		f(cas)
	}()
	select {
	case r := <-wai.c:
		return r
	case <-wai.e:
		return &Case{Timeout: true}
	}
}

func (ck *CaseHolder) Loop(recv <-chan *Case, timeout time.Duration) {
	ck.wg.Add(2)
	go func() {
		defer ck.wg.Done()
		ck.awake(recv)
	}()
	go func() {
		defer ck.wg.Done()
		ck.clear(timeout)
	}()
}

// 唤醒等待中的任务
func (ck *CaseHolder) awake(recv <-chan *Case) {
	for {
		select {
		case <-ck.context.Done():
			return
		case cas := <-recv:
			if v, ok := ck.m.Load(cas.ID); ok {
				wai := v.(*waiter)
				atomic.StoreInt32(&wai.remove, 1)
				wai.c <- cas
			}
		}
	}
}

// 清理超时任务
func (ck *CaseHolder) clear(timeout time.Duration) {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	ticker := time.NewTicker(timeout / 2)
	for {
		select {
		case <-ck.context.Done():
			return
		case <-ticker.C:
			arrRemove := []*waiter{}
			arrTimeout := []*waiter{}
			ck.m.Range(func(key, value any) bool {
				wai := value.(*waiter)
				if atomic.LoadInt32(&wai.remove) == int32(1) {
					arrRemove = append(arrRemove, wai)
				} else if time.Now().Sub(wai.create) > timeout {
					arrTimeout = append(arrTimeout, wai)
				}
				return true
			})
			for _, wai := range arrRemove {
				ck.m.Delete(wai.caseId)
			}
			for _, wai := range arrTimeout {
				ck.m.Delete(wai.caseId)
				close(wai.e)
			}
		}
	}
}
