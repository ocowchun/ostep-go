package vm_freespace

import (
	"bytes"
	"sort"
)

type FreeList struct {
	slots []Slot
}

func (l *FreeList) String() string {
	var out bytes.Buffer
	for _, slot := range l.slots {
		out.WriteString(slot.String())
	}
	return out.String()
}

func NewFreeList(bassAddr int, size int) *FreeList {
	return &FreeList{
		slots: []Slot{
			Slot{bassAddr, size},
		},
	}
}

func (l *FreeList) Size() int {
	return len(l.slots)
}

func (l *FreeList) Add(slot Slot) {
	l.slots = append(l.slots, slot)
	sort.Slice(l.slots, func(i, j int) bool {
		return l.slots[i].Addr < l.slots[j].Addr
	})
}
func (l *FreeList) Remove(victim Slot) {
	for idx, slot := range l.slots {
		if slot.Addr == victim.Addr {
			l.slots = append(l.slots[:idx], l.slots[idx+1:]...)
			return
		}
	}
}

func (l *FreeList) Slots() []Slot {
	return l.slots
}
