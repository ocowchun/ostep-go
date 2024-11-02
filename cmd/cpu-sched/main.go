package main

import (
	"fmt"
	cpu_schedule "ostep-go/cpu-schedule"
)

func main() {
	fmt.Println("cpu-sched")

	s := cpu_schedule.NewSimulator()
	jobs := []cpu_schedule.Job{
		{Name: "Job 0", Length: 8000},
		{Name: "Job 1", Length: 7000},
		{Name: "Job 2", Length: 1000},
		{Name: "Job 3", Length: 4000},
		{Name: "Job 4", Length: 5000},
		{Name: "Job 5", Length: 6000},
	}
	//jobs := []cpu_schedule.Job{
	//	{Name: "Job 0", Length: 1000},
	//	{Name: "Job 1", Length: 1000},
	//	{Name: "Job 2", Length: 1000},
	//	{Name: "Job 3", Length: 1000},
	//	{Name: "Job 4", Length: 1000},
	//	{Name: "Job 5", Length: 1000},
	//}

	//fmt.Println("====FIFO====")
	//fifo := cpu_schedule.NewFIFOScheduler()
	//s.Run(fifo, jobs)
	//fmt.Println()

	//fmt.Println("====SJF====")
	//sjf := cpu_schedule.NewSJFScheduler()
	//s.Run(sjf, jobs)
	//fmt.Println()

	fmt.Println("====Round Robin====")
	rr := cpu_schedule.NewRoundRobinScheduler()
	s.Run(rr, jobs)
	fmt.Println("cpu-sched done")
}
