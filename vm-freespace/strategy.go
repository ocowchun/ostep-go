package vm_freespace

import (
	"bytes"
	"fmt"
)

type Pointer uint32

type FreeSpaceStrategy interface {
	Alloc(pointer Pointer, size int) AllocResponse
	Free(pointer Pointer)
	FreeList() *FreeList
}

type Slot struct {
	Addr int
	Size int
}

func (s Slot) String() string {
	return fmt.Sprintf("[ addr:%d sz:%d ]", s.Addr, s.Size)
}

type BestStrategy struct {
	freeList *FreeList
	store    *Store
}

type AllocResponse struct {
	Err     error
	Visited int
	Addr    int
}

func (s *BestStrategy) Alloc(pointer Pointer, size int) AllocResponse {
	//ptr[0] = Alloc(3) returned 1000 (searched 1 elements)
	//Free List [ Size 1 ]: [ addr:1003 sz:97 ]

	var candidate Slot
	visited := 0
	for _, slot := range s.freeList.Slots() {
		if slot.Size > size && (candidate.Size == 0 || candidate.Size < slot.Size) {
			candidate = slot
		}
		visited++
	}
	if candidate.Size == 0 {
		return AllocResponse{Err: fmt.Errorf("no available slot")}
	}

	s.freeList.Remove(candidate)
	allocatedSlot := Slot{Addr: candidate.Addr, Size: size}
	s.store.Add(pointer, allocatedSlot)
	if candidate.Size > size {
		remainingSlot := Slot{Addr: candidate.Addr + size, Size: candidate.Size - size}
		s.freeList.Add(remainingSlot)
	}

	return AllocResponse{
		Err:     nil,
		Visited: visited,
		Addr:    allocatedSlot.Addr,
	}
}

func (s *BestStrategy) Free(pointer Pointer) {
	slot := s.store.Remove(pointer)
	s.freeList.Add(slot)
}

func (s *BestStrategy) FreeList() *FreeList {
	return s.freeList
}

func (s *BestStrategy) String() string {
	var output bytes.Buffer
	output.WriteString("freelist:")
	output.WriteString(s.freeList.String())
	output.WriteString(",")
	output.WriteString("store:")
	output.WriteString(s.store.String())
	return output.String()
}

func MakeFreeSpaceStrategy(strategyName string, baseAddr int, size int) (FreeSpaceStrategy, error) {
	switch strategyName {
	case "BEST":
		return &BestStrategy{
			freeList: NewFreeList(baseAddr, size),
			store:    NewStore(),
		}, nil

	default:
		return nil, fmt.Errorf("unknown strategy %s", strategyName)
	}
}
