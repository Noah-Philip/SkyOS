package internal

import "sync"

type Scheduler struct {
	Workload []*Task
	mu       sync.Mutex
}

func (s *Scheduler) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
}

func (s *Scheduler) AddTask(task *Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Workload = append(s.Workload, task)
}

func (s *Scheduler) RemoveTask(taskToRemove *Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, task := range s.Workload {
		if task == taskToRemove {
			s.Workload = append(
				s.Workload[:i], s.Workload[i+1:]...,
			)
			return
		}
	}
}

func (s *Scheduler) GetWorkload() []*Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Workload
}

func (s *Scheduler) SetWorkload(workload []*Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Workload = workload
}

func (s *Scheduler) ClearWorkload() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Workload = nil
}

func (s *Scheduler) IsEmpty() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.Workload) == 0
}

func (s *Scheduler) Size() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.Workload)
}

func (s *Scheduler) scoreNodes(Task *Task, nodes []*Node) map[*Node]float64 {
	scores := make(map[*Node]float64)
	for _, node := range nodes {
		scores[node] = 0
		if node.isEligibleForTask(Task) {
			scores[node] = node.calculateScore(Task)
		}
	}
	return scores
}
