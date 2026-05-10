package allocator

// dpSingleServer solves the 2-D 0/1 knapsack for one server.
// It returns the selected task indices and the total priority achieved.
func dpSingleServer(tasks []Task, cpuCap, ramCap int) ([]int, int) {
	n := len(tasks)

	// dp[i][c][r] = best priority using tasks 0..i-1 with c CPU, r RAM
	// Allocate as a flat slice for cache-friendliness.
	size := (n + 1) * (cpuCap + 1) * (ramCap + 1)
	flat := make([]int, size)

	idx := func(i, c, r int) int {
		return i*(cpuCap+1)*(ramCap+1) + c*(ramCap+1) + r
	}

	for i := 1; i <= n; i++ {
		ci := tasks[i-1].CPU
		ri := tasks[i-1].RAM
		pi := tasks[i-1].Priority
		for c := 0; c <= cpuCap; c++ {
			for r := 0; r <= ramCap; r++ {
				skip := flat[idx(i-1, c, r)]
				flat[idx(i, c, r)] = skip
				if c >= ci && r >= ri {
					take := flat[idx(i-1, c-ci, r-ri)] + pi
					if take > skip {
						flat[idx(i, c, r)] = take
					}
				}
			}
		}
	}

	best := flat[idx(n, cpuCap, ramCap)]

	// Backtrack
	selected := make([]int, 0)
	c, r := cpuCap, ramCap
	for i := n; i >= 1; i-- {
		if flat[idx(i, c, r)] != flat[idx(i-1, c, r)] {
			selected = append(selected, i-1)
			c -= tasks[i-1].CPU
			r -= tasks[i-1].RAM
		}
	}
	return selected, best
}

// DP allocates tasks across servers using dynamic programming.
// Each server is filled optimally via the 2-D knapsack; tasks placed on
// one server are removed from the pool before solving the next.
func DP(tasks []Task, servers []Server) Allocation {
	alloc := make(Allocation, len(tasks))
	remaining := make([]Task, len(tasks))
	copy(remaining, tasks)

	for j := range servers {
		if len(remaining) == 0 {
			break
		}
		sel, _ := dpSingleServer(remaining, servers[j].CPUTotal, servers[j].RAMTotal)
		placed := make(map[int]bool, len(sel))
		for _, idx := range sel {
			t := remaining[idx]
			alloc[t.ID] = servers[j].ID
			servers[j].CPURemaining -= t.CPU
			servers[j].RAMRemaining -= t.RAM
			placed[idx] = true
		}
		// Remove placed tasks
		next := make([]Task, 0, len(remaining)-len(sel))
		for k, t := range remaining {
			if !placed[k] {
				next = append(next, t)
			}
		}
		remaining = next
	}

	for _, t := range remaining {
		alloc[t.ID] = ""
	}
	return alloc
}
