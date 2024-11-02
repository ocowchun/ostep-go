package vm_freespace

import (
	"bytes"
	"fmt"
)

type Store struct {
	store map[Pointer]Slot
}

func NewStore() *Store {
	return &Store{make(map[Pointer]Slot)}
}

func (s *Store) Add(pointer Pointer, slot Slot) {
	s.store[pointer] = slot
}

func (s *Store) Remove(pointer Pointer) Slot {
	slot := s.store[pointer]
	delete(s.store, pointer)
	return slot
}

func (s *Store) String() string {
	var out bytes.Buffer
	for pointer, slot := range s.store {
		out.WriteString(fmt.Sprintf("%d->%s\n", pointer, slot))
	}
	return out.String()
}
