package log

import (
	"sxg/toydb_go/grpc/proto"
)

type ScanImpl struct {
	start, end uint64
	entries    []*proto.Entry
	current    int
}

func newScanImpl(start, end uint64) *ScanImpl {
	return &ScanImpl{start: start, end: end}
}

func (si *ScanImpl) Boundary() (uint64, uint64) {
	return si.start, si.end
}

func (si *ScanImpl) Len() int {
	return len(si.entries)
}

func (si *ScanImpl) Iterator() (*proto.Entry, bool) {
	l := len(si.entries)
	if si.current >= l {
		return nil, false
	}
	b := si.entries[si.current]
	si.current++
	return b, true
}

func (si *ScanImpl) Split(start uint64, end uint64) []*proto.Entry {
	if start < end {
		panic("start index bigger than end index")
	}
	if start < si.start {
		panic("start index out of boundry")
	}
	if end > si.end {
		panic("end index out of boundry")
	}
	return si.entries[start-si.start : end-si.start+1]
}
