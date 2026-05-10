package allocator

import "sort"

// Heuristic allocates tasks using a weighted least-loaded scoring strategy.
// Tasks are sorted by priority (highest first).  For each task the server
// with the highest remaining-capacity score is selected:
//
//	score = alpha * (cpuRemaining/cpuTotal) + beta * (ramRemaining/ramTotal)
func Heuristic(tasks []Task, servers []Server, alpha, beta float64) Allocation {
	sorted := make([]Task, len(tasks))
	copy(sorted, tasks)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Priority > sorted[j].Priority
	})

	alloc := make(Allocation, len(tasks))

	for i := range sorted {
		t := &sorted[i]
		bestIdx := -1
		bestScore := -1.0

		for j := range servers {
			s := &servers[j]
			if s.CPURemaining < t.CPU || s.RAMRemaining < t.RAM {
				continue
			}
			cpuRatio := float64(s.CPURemaining) / float64(s.CPUTotal)
			ramRatio := float64(s.RAMRemaining) / float64(s.RAMTotal)
			score := alpha*cpuRatio + beta*ramRatio
			if score > bestScore {
				bestScore = score
				bestIdx = j
			}
		}

		if bestIdx >= 0 {
			servers[bestIdx].CPURemaining -= t.CPU
			servers[bestIdx].RAMRemaining -= t.RAM
			alloc[t.ID] = servers[bestIdx].ID
		} else {
			alloc[t.ID] = ""
		}
	}
	return alloc
}
