package cpu_schedule

import (
	"fmt"
	"sort"
	"strconv"
)

type Scheduler interface {
	Schedule(jobs []Job) (Job, error)
}

type FIFOScheduler struct {
}

func NewFIFOScheduler() *FIFOScheduler {
	return &FIFOScheduler{}
}

func (f FIFOScheduler) Schedule(jobs []Job) (Job, error) {
	if len(jobs) == 0 {
		return Job{}, fmt.Errorf("no jobs provided")
	}

	return jobs[0], nil
}

type SJFScheduler struct{}

func NewSJFScheduler() *SJFScheduler {
	return &SJFScheduler{}
}
func (s SJFScheduler) Schedule(jobs []Job) (Job, error) {
	if len(jobs) == 0 {
		return Job{}, fmt.Errorf("no jobs provided")
	}

	var best Job
	for _, job := range jobs {
		if best.Length == 0 || job.Length < best.Length {
			best = job
		}
	}

	return best, nil
}

type RoundRobinScheduler struct {
	runs map[string]int64
}

func NewRoundRobinScheduler() *RoundRobinScheduler {
	return &RoundRobinScheduler{
		runs: make(map[string]int64),
	}
}

func (r *RoundRobinScheduler) Schedule(jobs []Job) (Job, error) {
	if len(jobs) == 0 {
		return Job{}, fmt.Errorf("no jobs provided")
	}

	var best Job
	for _, job := range jobs {
		if best.Length == 0 {
			best = job
			continue
		}
		if r.runs[job.Name] < r.runs[best.Name] {
			best = job
		}
	}

	r.runs[best.Name] += 1
	return best, nil
}

type Simulator struct {
}

func NewSimulator() *Simulator {
	return &Simulator{}
}

type Job struct {
	Name   string
	Length uint64
}

type Jobs struct {
	jobs []Job
}

func newJobs(jobs []Job) *Jobs {
	copied := make([]Job, len(jobs))
	copy(copied, jobs)
	return &Jobs{jobs: copied}
}

func (j *Jobs) remove(name string) {
	idx := -1
	for i, job := range j.jobs {
		if job.Name == name {
			idx = i
		}
	}
	if idx != -1 {
		j.jobs = append(j.jobs[:idx], j.jobs[idx+1:]...)
	}
}

func (j *Jobs) value() []Job {
	return j.jobs
}
func (j *Jobs) isEmpty() bool {
	return len(j.jobs) == 0
}

func (s *Simulator) Run(scheduler Scheduler, jobs []Job) {
	js := newJobs(jobs)

	currentTimeInMs := int64(0)
	step := int64(1000)
	arrivals := make(map[string]int64)
	executeds := make(map[string]int64)
	firstSchedules := make(map[string]int64)
	completedAts := make(map[string]int64)

	for _, job := range js.value() {
		arrivals[job.Name] = currentTimeInMs
		executeds[job.Name] = 0
		firstSchedules[job.Name] = -1
		completedAts[job.Name] = -1
	}

	for !js.isEmpty() {
		job, err := scheduler.Schedule(js.value())
		if err != nil {
			panic(err)
		}

		if firstSchedules[job.Name] == -1 {
			firstSchedules[job.Name] = currentTimeInMs
		}
		executeds[job.Name] += step

		currentTimeInMs = currentTimeInMs + step
		if executeds[job.Name] >= int64(job.Length) {
			completedAts[job.Name] = currentTimeInMs
			js.remove(job.Name)
		}
	}

	// print stats
	totalResponse := float64(0)
	totalTurnaround := float64(0)
	totalWait := float64(0)

	sort.Slice(jobs, func(i, j int) bool {
		return firstSchedules[jobs[i].Name] < firstSchedules[jobs[j].Name]
	})
	for _, job := range jobs {
		response := float64(firstSchedules[job.Name]) / 1000
		turnaround := float64(completedAts[job.Name]) / 1000
		wait := turnaround - (float64(arrivals[job.Name]) / 1000) - (float64(job.Length / 1000))

		fmt.Printf("%s -- Response: %s Turnaround %s Wait %s\n", job.Name,
			strconv.FormatFloat(response, 'f', 2, 64),
			strconv.FormatFloat(turnaround, 'f', 2, 64),
			strconv.FormatFloat(wait, 'f', 2, 64),
		)
		totalResponse += response
		totalTurnaround += turnaround
		totalWait += wait
	}

	fmt.Println()
	jobLen := float64(len(jobs))
	fmt.Printf("Average -- Response: %s Turnaround %s Wait %s\n",
		strconv.FormatFloat(totalResponse/jobLen, 'f', 2, 64),
		strconv.FormatFloat(totalTurnaround/jobLen, 'f', 2, 64),
		strconv.FormatFloat(totalWait/jobLen, 'f', 2, 64),
	)

}

type Simulation struct {
	currentTimeInMs uint64
	jobs            []Job
}

// create a scheduler (either FIFO, SJF or RR), a simulator and share the final result
