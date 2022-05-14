package log

type EntryProducer func(uint64, uint64, []byte) Entry
type Entry interface {
	Index() uint64
	Term() uint64
	Command() []byte
}

type Scan interface {
	Boundary() (uint64, uint64)
	Iterator() (Entry, bool)
	Len() int
}

//注意：保留一定数量的committed的日志，不要全部删除。可以设置一定的阈值
type Store interface {
	Append(entry Entry) (uint64, error)    //新增日志，返回索引
	Commit(index uint64) error             //commit日志
	Commited() uint64                      //返回已经commit的日志索引
	Get(index uint64) (Entry, bool)        //获取指定索引的entry
	GetLast() (Entry, bool)                //获取最后一个日志
	Len() int                              //获取日志长度
	Scan(start, end uint64) Scan           //获取指定索引之间的entry
	Size() int                             //获取日志大小(byte的数量)
	Truncate(index uint64) (uint64, error) //删除指定索引之后的日志，并返回最大的索引
	GetMetaData(key []byte) ([]byte, bool) //获取元数据的值
	SetMetaData(key, value []byte) error   //设置元数据的值
	IsEmpty() bool                         //日志是否有entry
}
