package log

import (
	"sxg/toydb_go/grpc/proto"
)

const (
	Endless = 0
)

type EntryProducer func(uint64, uint64, []byte) *proto.Entry

type Scan interface {
	Boundary() (uint64, uint64)
	Iterator() (*proto.Entry, bool)
	Len() int
	Split(uint64, uint64) []*proto.Entry
}

//注意：保留一定数量的committed的日志，不要全部删除。可以设置一定的阈值
type Store interface {
	Append(entry *proto.Entry) (uint64, error) //新增日志，返回索引
	Commit(index uint64) error                 //commit日志
	Commited() uint64                          //返回已经commit的日志索引
	Get(index uint64) (*proto.Entry, bool)     //获取指定索引的entry
	GetLast() (*proto.Entry, bool)             //获取最后一个日志
	Len() int                                  //获取日志长度
	Scan(start, end uint64) Scan               //获取指定索引之间的entry, 若end=Endless(0), 则获取自start开始的所有日志
	Size() int                                 //获取日志大小(byte的数量)
	Truncate(index uint64) (uint64, error)     //删除指定索引之后的日志，并返回最大的索引
	GetMetaData(key []byte) ([]byte, bool)     //获取元数据的值
	SetMetaData(key, value []byte) error       //设置元数据的值
	IsEmpty() bool                             //日志是否有entry
	String() string                            //当前Store的名称
}
