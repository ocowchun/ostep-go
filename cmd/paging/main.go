package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	input, err := parse(os.Stdin)
	if err != nil {
		panic(err)
	}

	table := input.pageTable
	for _, trace := range input.traces {
		va := trace.virtualAddress

		pa, err := table.Translate(va)
		if err != nil {
			fmt.Printf("VA 0x%08x (decimal:    %d) -->  Invalid (%s)\n", va, va, err.Error())
		} else {
			vpn := va >> table.offsetSize
			fmt.Printf("VA 0x%08x (decimal:    %d) --> %08x (decimal    %d) [VPN %d]\n", va, va, pa, pa, vpn)

		}
	}
}

type VirtualAddressTrace struct {
	virtualAddress int64
}

func (t VirtualAddressTrace) String() string {
	return fmt.Sprintf("%08x", t.virtualAddress)
}

type Problem struct {
	pageTable        *PageTable
	traces           []VirtualAddressTrace
	addressSpaceSize int64
}

type PageTable struct {
	offsetSize int
	pfnSize    int
	entries    []int64
}

func (t *PageTable) Translate(virtualAddress int64) (int64, error) {
	vpn := virtualAddress >> t.offsetSize
	if int64(len(t.entries)) <= vpn || t.entries[vpn] == 0 {
		return 0, fmt.Errorf("VPN %d not valid", vpn)
	}

	offsetMask := int64((1 << t.offsetSize) - 1)
	offset := virtualAddress & offsetMask
	pfnMask := 1<<(t.pfnSize) - 1
	pa := t.entries[vpn] & int64(pfnMask)
	pa = (pa << t.offsetSize) + offset

	return pa, nil
}

func (t *PageTable) AddEntry(entry int64) {
	t.entries = append(t.entries, entry)
}

func (t *PageTable) String() string {
	var out bytes.Buffer
	for idx, entry := range t.entries {
		out.WriteString(fmt.Sprintf("%3d : %08x\n", idx, entry))
	}
	return out.String()
}

func parse(reader io.Reader) (Problem, error) {
	scanner := bufio.NewScanner(reader)

	problem := Problem{}
	pageSize := int64(0)
	physicalMemorySize := int64(0)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "Page Table") {
			pfnSize := computeLargestBit(uint64(physicalMemorySize / pageSize))
			offsetSize := computeLargestBit(uint64(pageSize))

			table := &PageTable{offsetSize: offsetSize, pfnSize: pfnSize}
			for scanner.Scan() {
				line := scanner.Text()

				if strings.HasPrefix(line, "  [") {
					tokens := strings.Split(line, "0x")
					i, err := strconv.ParseInt(tokens[1], 16, 64)
					if err != nil {
						return problem, err
					}

					table.AddEntry(i)
				} else {
					problem.pageTable = table
					break
				}
			}

		} else if strings.HasPrefix(line, "Virtual Address Trace") {
			traces := make([]VirtualAddressTrace, 0)
			for scanner.Scan() {
				line := scanner.Text()

				if strings.HasPrefix(line, "  VA") {
					i, err := strconv.ParseInt(line[7:15], 16, 64)
					if err != nil {
						return problem, err
					}
					traces = append(traces, VirtualAddressTrace{virtualAddress: i})
				} else {
					problem.traces = traces
					break
				}
			}

		} else if strings.HasPrefix(line, "ARG page size") {
			tokens := strings.Split(line, "ARG page size ")
			num, err := convertNum(tokens[1])
			if err != nil {
				return problem, err
			}
			pageSize = num
		} else if strings.HasPrefix(line, "ARG address space size") {
			tokens := strings.Split(line, "ARG address space size ")
			num, err := convertNum(tokens[1])
			if err != nil {
				return problem, err
			}
			problem.addressSpaceSize = num
		} else if strings.HasPrefix(line, "ARG phys mem size") {
			tokens := strings.Split(line, "ARG phys mem size ")
			num, err := convertNum(tokens[1])
			if err != nil {
				return problem, err
			}
			physicalMemorySize = num
		}
	}

	return problem, nil
}

func convertNum(str string) (int64, error) {
	unit := str[len(str)-1]
	num, err := strconv.ParseInt(str[:len(str)-1], 10, 64)
	if err != nil {
		return 0, err
	}

	switch unit {
	case 'k':
		return num * 1024, nil
	case 'm':
		return num * 1024 * 1024, nil

	default:
		return 0, fmt.Errorf("invalid str %s", str)

	}
}
func computeLargestBit(num uint64) int {
	i := 63
	for i >= 0 {
		if num >= (1 << i) {
			return i
		}
		i--
	}
	return 0
}
