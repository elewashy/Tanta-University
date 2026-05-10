package allocator_test

import (
	"cloud-alloc/allocator"
	"math/rand"
	"testing"
)

const (
	seed    = 42
	cpuCap  = 16
	ramCap  = 32
)

func setup(nTasks, nServers int) ([]allocator.Task, []allocator.Server) {
	rng := rand.New(rand.NewSource(seed))
	tasks := allocator.GenerateTasks(nTasks, rng)
	servers := allocator.GenerateServers(nServers, cpuCap, ramCap)
	return tasks, servers
}

// ── Correctness tests ──

func TestGreedy_NoOverflow(t *testing.T) {
	tasks, servers := setup(50, 5)
	alloc := allocator.Greedy(tasks, servers)
	assertNoOverflow(t, alloc, tasks, servers, "Greedy")
}

func TestDP_NoOverflow(t *testing.T) {
	tasks, servers := setup(50, 5)
	alloc := allocator.DP(tasks, servers)
	assertNoOverflow(t, alloc, tasks, servers, "DP")
}

func TestHeuristic_NoOverflow(t *testing.T) {
	tasks, servers := setup(50, 5)
	alloc := allocator.Heuristic(tasks, servers, 0.5, 0.5)
	assertNoOverflow(t, alloc, tasks, servers, "Heuristic")
}

func TestDP_BetterOrEqualToGreedy(t *testing.T) {
	tasks, sG := setup(30, 3)
	_, sD := setup(30, 3)

	aG := allocator.Greedy(tasks, sG)
	aD := allocator.DP(tasks, sD)

	scoreG := priorityScore(aG, tasks)
	scoreD := priorityScore(aD, tasks)

	if scoreD < scoreG {
		t.Errorf("DP score %d < Greedy score %d; DP should be >= Greedy", scoreD, scoreG)
	}
}

func TestAllTasks_InAllocation(t *testing.T) {
	tasks, servers := setup(50, 5)
	for name, fn := range map[string]func([]allocator.Task, []allocator.Server) allocator.Allocation{
		"Greedy":    allocator.Greedy,
		"Heuristic": func(tasks []allocator.Task, servers []allocator.Server) allocator.Allocation {
			return allocator.Heuristic(tasks, servers, 0.5, 0.5)
		},
		"DP": allocator.DP,
	} {
		s := allocator.CloneServers(servers)
		alloc := fn(tasks, s)
		if len(alloc) != len(tasks) {
			t.Errorf("%s: allocation has %d entries, expected %d", name, len(alloc), len(tasks))
		}
	}
}

// ── Benchmarks ──

func BenchmarkGreedy_50(b *testing.B) {
	tasks, _ := setup(50, 5)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := allocator.GenerateServers(5, cpuCap, ramCap)
		allocator.Greedy(tasks, s)
	}
}

func BenchmarkDP_50(b *testing.B) {
	tasks, _ := setup(50, 5)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := allocator.GenerateServers(5, cpuCap, ramCap)
		allocator.DP(tasks, s)
	}
}

func BenchmarkHeuristic_50(b *testing.B) {
	tasks, _ := setup(50, 5)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := allocator.GenerateServers(5, cpuCap, ramCap)
		allocator.Heuristic(tasks, s, 0.5, 0.5)
	}
}

func BenchmarkGreedy_500(b *testing.B) {
	rng := rand.New(rand.NewSource(seed))
	tasks := allocator.GenerateTasks(500, rng)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := allocator.GenerateServers(50, cpuCap, ramCap)
		allocator.Greedy(tasks, s)
	}
}

func BenchmarkDP_500(b *testing.B) {
	rng := rand.New(rand.NewSource(seed))
	tasks := allocator.GenerateTasks(500, rng)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := allocator.GenerateServers(50, cpuCap, ramCap)
		allocator.DP(tasks, s)
	}
}

func BenchmarkHeuristic_500(b *testing.B) {
	rng := rand.New(rand.NewSource(seed))
	tasks := allocator.GenerateTasks(500, rng)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := allocator.GenerateServers(50, cpuCap, ramCap)
		allocator.Heuristic(tasks, s, 0.5, 0.5)
	}
}

// ── Helpers ──

func assertNoOverflow(t *testing.T, alloc allocator.Allocation, tasks []allocator.Task, servers []allocator.Server, name string) {
	t.Helper()
	taskMap := make(map[string]allocator.Task, len(tasks))
	for _, tk := range tasks {
		taskMap[tk.ID] = tk
	}

	cpuUsed := make(map[string]int)
	ramUsed := make(map[string]int)
	for tid, sid := range alloc {
		if sid == "" {
			continue
		}
		cpuUsed[sid] += taskMap[tid].CPU
		ramUsed[sid] += taskMap[tid].RAM
	}

	for _, s := range servers {
		if cpuUsed[s.ID] > s.CPUTotal {
			t.Errorf("%s: server %s CPU overflow: used %d / %d", name, s.ID, cpuUsed[s.ID], s.CPUTotal)
		}
		if ramUsed[s.ID] > s.RAMTotal {
			t.Errorf("%s: server %s RAM overflow: used %d / %d", name, s.ID, ramUsed[s.ID], s.RAMTotal)
		}
	}
}

func priorityScore(alloc allocator.Allocation, tasks []allocator.Task) int {
	taskMap := make(map[string]allocator.Task, len(tasks))
	for _, t := range tasks {
		taskMap[t.ID] = t
	}
	score := 0
	for tid, sid := range alloc {
		if sid != "" {
			score += taskMap[tid].Priority
		}
	}
	return score
}
