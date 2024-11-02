package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	vm_freespace "ostep-go/vm-freespace"
	"strconv"
	"strings"
)

func foo() {
	best, err := vm_freespace.MakeFreeSpaceStrategy("BEST", 1000, 100)
	if err != nil {
		panic(err)
	}

	res := best.Alloc(vm_freespace.Pointer(0), 3)
	best.Free(vm_freespace.Pointer(0))
	res = best.Alloc(vm_freespace.Pointer(1), 5)
	fmt.Println(res.Visited)
	best.Free(vm_freespace.Pointer(1))
	//err = best.Alloc(vm_freespace.Pointer(2), 8)
	//best.Free(vm_freespace.Pointer(2))
	fmt.Println(best)
}

func main() {
	input, err := parse(os.Stdin)
	if err != nil {
		panic(err)
	}

	strategy, err := vm_freespace.MakeFreeSpaceStrategy(input.StrategyName, input.BaseAddr, input.Space)
	for _, op := range input.Operations {
		switch op := op.(type) {
		case AllocOperation:
			res := strategy.Alloc(vm_freespace.Pointer(op.PointerIndex), op.Size)
			if res.Err != nil {
				panic(res.Err)
			}
			//			ptr[4] = Alloc(2) returned 1000 (searched 4 elements)
			//Free List [ Size 4 ]: [ addr:1002 sz:1 ][ addr:1003 sz:5 ][ addr:1008 sz:8 ][ addr:1016 sz:84 ]
			fmt.Printf("ptr[%d] = Alloc(%d) return %d (searched %d elements)\n", op.PointerIndex, op.Size, res.Addr, res.Visited)
			fmt.Printf("Free List [ Size %d ]: %s\n", strategy.FreeList().Size(), strategy.FreeList().String())
			fmt.Println()
		case FreeOperation:
			strategy.Free(vm_freespace.Pointer(op.PointerIndex))
			fmt.Printf("Free(ptr[%d])\n", op.PointerIndex)
			fmt.Printf("returned %d\n", 0)
			fmt.Printf("Free List [ Size %d ]: %s\n", strategy.FreeList().Size(), strategy.FreeList().String())
			fmt.Println()

		}
	}

	//
	//fmt.Println(ops)

	//reader := os.Stdin
	//lines := make([]string, 0)
	//for scanner.Scan() {
	//	lines = append(lines, scanner.Text())
	//}

	//ptr[0] = Alloc(3) returned ?
	//List?
	//
	//Free(ptr[0])
	//returned ?
	//List?
	//
	//ptr[1] = Alloc(5) returned ?
	//List?
	//

	//ptr[0] = Alloc(3) returned 1000 (searched 1 elements)
	//Free List [ Size 1 ]: [ addr:1003 sz:97 ]
	//
	//Free(ptr[0])
	//returned 0
	//Free List [ Size 2 ]: [ addr:1000 sz:3 ][ addr:1003 sz:97 ]
	//

}

type OperationType uint8

const (
	Alloc OperationType = iota
	Free
)

type Operation interface {
	Type() OperationType
}

type AllocOperation struct {
	PointerIndex int
	Size         int
}

type SimulationInput struct {
	Space        int
	BaseAddr     int
	StrategyName string
	Operations   []Operation
}

func (op AllocOperation) Type() OperationType {
	return Alloc
}

type FreeOperation struct {
	PointerIndex int
}

func (op FreeOperation) Type() OperationType {
	return Free
}

func parse(reader io.Reader) (SimulationInput, error) {
	scanner := bufio.NewScanner(reader)
	ops := make([]Operation, 0)
	sim := SimulationInput{}
	for scanner.Scan() {
		firstLine := scanner.Text()
		if strings.HasPrefix(firstLine, "ptr[") {
			tokens := strings.Split(firstLine, " ")
			ptr, err := strconv.Atoi(tokens[0][4 : len(tokens[0])-1])
			if err != nil {
				return sim, err
			}

			size, err := strconv.Atoi(tokens[2][6 : len(tokens[2])-1])
			if err != nil {
				return sim, err
			}
			ops = append(ops, AllocOperation{PointerIndex: ptr, Size: size})

		} else if strings.HasPrefix(firstLine, "Free(") {
			ptr, err := strconv.Atoi(firstLine[9 : len(firstLine)-2])
			if err != nil {
				return sim, err
			}
			ops = append(ops, FreeOperation{PointerIndex: ptr})
		} else if strings.HasPrefix(firstLine, "size") {
			space, err := strconv.Atoi(firstLine[5:])
			if err != nil {
				return sim, err
			}
			sim.Space = space

		} else if strings.HasPrefix(firstLine, "baseAddr") {
			baseAddr, err := strconv.Atoi(firstLine[9:])
			if err != nil {
				return sim, err
			}
			sim.BaseAddr = baseAddr
		} else if strings.HasPrefix(firstLine, "policy") {
			sim.StrategyName = firstLine[7:]
		}
	}
	sim.Operations = ops

	return sim, nil
}
