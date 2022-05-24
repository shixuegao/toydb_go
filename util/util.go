package util

import "sort"

const (
	ASC = iota
	DESC
)

type Number interface {
	~byte | int16 | int32 | int64 | uint16 | uint32 | uint64
}

type sorting[T Number] struct {
	arr       []T
	direction int
}

func (s *sorting[T]) Len() int {
	return len(s.arr)
}

func (s *sorting[T]) Less(i, j int) bool {
	if s.direction == ASC {
		return s.arr[i] <= s.arr[j]
	} else {
		return s.arr[i] > s.arr[j]
	}
}

func (s *sorting[T]) Swap(i, j int) {
	s.arr[i] = s.arr[j]
}

func Sort[T Number](arr []T, direction int) {
	sort.Sort(&sorting[T]{arr, direction})
}
