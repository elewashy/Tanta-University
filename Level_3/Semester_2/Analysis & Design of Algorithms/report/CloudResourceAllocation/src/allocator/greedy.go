package allocator

import "sort"

// Greedy allocates tasks using a priority-sorted first-fit strategy.
// Tasks are sorted by priority (highest first).  For each task the first
// server with enough remaining CPU and RAM receives the assignment.
func Greedy(tasks []Task, servers []Server) Allocation {
	sorted := make([]Task, len(tasks))
	copy(sorted, tasks)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Priority > sorted[j].Priority
	})

	alloc := make(Allocation, len(tasks))

	for i := range sorted {
		t := &sorted[i]
		placed := false
		for j := range servers {
			s := &servers[j]
			if s.CPURemaining >= t.CPU && s.RAMRemaining >= t.RAM {
				s.CPURemaining -= t.CPU
				s.RAMRemaining -= t.RAM
				alloc[t.ID] = s.ID
				placed = true
				break
			}
		}
		if !placed {
			alloc[t.ID] = ""
		}
	}
	return alloc
}
