package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	problem, err := parse(os.Stdin)
	if err != nil {
		panic(err)
	}

	solve(problem)
}

func solve(problem *Problem) {
	pfnMask := (1 << 7) - 1
	ptIndexMask := ((1 << 5) - 1) << 5
	offsetMask := (1 << 5) - 1

	for _, va := range problem.virtualAddress {
		fmt.Printf("Virtual Address 0x%04x:\n", va)

		pdIndex := int(va >> 10)

		pdPage := problem.physicalPages[problem.PageDirectoryPageNum]
		pde := pdPage[pdIndex]
		pdeValidBit := pde >> 7
		pdePFN := int(pde) & pfnMask
		fmt.Printf("  --> pde index:0x%x [decimal %d] pde contents:0x%02x (valid %d, pfn 0x%02x [decimal %d])\n", pdIndex, pdIndex, pde, pdeValidBit, pdePFN, pdePFN)

		ptPage := problem.physicalPages[pdePFN]
		ptIndex := (va & ptIndexMask) >> 5
		pte := ptPage[ptIndex]
		pteValidBit := pte >> 7
		ptePFN := int(pte) & pfnMask
		fmt.Printf("    --> pte index:0x%x [decimal %d] pte contents:0x%x (valid %d, pfn 0x%x [decimal %d])\n", ptIndex, ptIndex, pte, pteValidBit, ptePFN, ptePFN)

		offset := va & offsetMask
		entryPage := problem.physicalPages[ptePFN]
		entry := entryPage[offset]
		physicalAddress := ptePFN<<5 + offset
		if pdeValidBit == 1 && pteValidBit == 1 {
			fmt.Printf("      --> Translates to Physical Address 0x%x --> Value: 0x%02x\n", physicalAddress, entry)
		} else {
			fmt.Printf("      --> Fault (page table entry not valid)\n")
		}

		//Virtual Address 0x611c:
		//  --> pde index:0x18 [decimal 24] pde contents:0xa1 (valid 1, pfn 0x21 [decimal 33])
		//    --> pte index:0x8 [decimal 8] pte contents:0xb5 (valid 1, pfn 0x35 [decimal 53])
		//      --> Translates to Physical Address 0x6bc --> Value: 0x08

	}

}

//ARG seed 0
//ARG allocated 64
//ARG num 10
//
//page   0:1b1d05051d0b19001e00121c1909190c0f0b0a1218151700100a061c06050514
//page   1:0000000000000000000000000000000000000000000000000000000000000000
//page   2:121b0c06001e04130f0b10021e0f000c17091717071e001a0f0408120819060b

type Problem struct {
	PageDirectoryPageNum int
	physicalPages        [128][32]byte
	virtualAddress       []int
}

func parse(reader io.Reader) (*Problem, error) {
	scanner := bufio.NewScanner(reader)

	problem := &Problem{}
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "PDBR: ") {
			tokens := strings.Split(line, " ")
			pageNum, err := strconv.ParseInt(tokens[1], 10, 64)
			if err != nil {
				return nil, err
			}
			problem.PageDirectoryPageNum = int(pageNum)

		} else if strings.HasPrefix(line, "page ") {
			tokens := strings.Split(line, ":")
			pageNum, err := strconv.ParseInt(strings.TrimSpace(tokens[0][4:]), 10, 64)
			if err != nil {
				return nil, err
			}

			i := 0
			content := tokens[1]
			for i < 32 {
				b, _ := strconv.ParseInt(content[i*2:i*2+2], 16, 16)
				problem.physicalPages[pageNum][i] = byte(b)

				i += 1
			}

		} else if strings.HasPrefix(line, "Virtual Address") {
			// Virtual Address 611c: Translates To What Physical Address (And Fetches what Value)? Or Fault?
			tokens := strings.Split(line, " ")
			virtualAddress, err := strconv.ParseInt(tokens[2][:4], 16, 64)
			if err != nil {
				return nil, err
			}
			problem.virtualAddress = append(problem.virtualAddress, int(virtualAddress))
		}
	}
	return problem, nil
}

//Some basic assumptions:
//
//* The page size is an unrealistically-small 32 bytes
//* The virtual address space for the process in question (assume there is only one) is 1024 pages, or 32 KB
//* physical memory consists of 128 pages
//
//Thus, a virtual address needs 15 bits (5 for the offset, 10 for the VPN).
//A physical address requires 12 bits (5 offset, 7 for the PFN).
//
//The system assumes a multi-level page table. Thus, the upper five bits of a virtual
//address are used to index into a page directory; the page directory entry (PDE), if valid,
//points to a page of the page table. Each page table page holds 32 page-table entries
//(PTEs). Each PTE, if valid, holds the desired translation (physical frame number, or PFN)
//of the virtual page in question.
//
//The format of a PTE is thus:
//
//```sh
//  VALID | PFN6 ... PFN0
//```
//
//and is thus 8 bits or 1 byte.
//
//The format of a PDE is essentially identical:
//
//```sh
//  VALID | PT6 ... PT0
//```
//
//You are given two pieces of information to begin with.
//
//First, you are given the value of the page directory base register (PDBR),
//which tells you which page the page directory is located upon.
//

// virtual address -> 611c -> 110000100011100
// 11000 -> index into page directory -> 24
// 01000 -> index into PDE -> 8
// 11100 -> offset

//PDBR 108 -> page directory in page 108
// PTE 0 -> valid bit, 1234567 -> PFN
// PDE 0 -> valid bit, 1234567 -> PT
