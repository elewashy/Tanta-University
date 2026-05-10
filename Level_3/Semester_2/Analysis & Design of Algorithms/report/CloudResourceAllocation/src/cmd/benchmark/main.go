// Command benchmark runs the three allocation algorithms on identical workloads
// at multiple scales, prints a formatted table, and writes results to
// results/benchmark.json.
package main

import (
	"cloud-alloc/allocator"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// scenario config
type config struct {
	NTasks   int `json:"n_tasks"`
	NServers int `json:"n_servers"`
	CPUCap   int `json:"cpu_cap"`
	RAMCap   int `json:"ram_cap"`
}

type algoResult struct {
	allocator.Metrics
}

type primaryOutput struct {
	Config  config                `json:"config"`
	Results map[string]algoResult `json:"results"`
}

type scaleEntry struct {
	TimeMs float64 `json:"time_ms"`
	Placed int     `json:"placed"`
	Total  int     `json:"total"`
}

type output struct {
	Primary     primaryOutput                    `json:"primary"`
	Scalability map[string]map[string]scaleEntry `json:"scalability"`
}

func runAlgo(name string, tasks []Task, servers []Server) (allocator.Allocation, time.Duration) {
	s := allocator.CloneServers(servers)
	start := time.Now()

	var alloc allocator.Allocation
	switch name {
	case "Greedy":
		alloc = allocator.Greedy(tasks, s)
	case "DP":
		alloc = allocator.DP(tasks, s)
	case "Heuristic":
		alloc = allocator.Heuristic(tasks, s, 0.5, 0.5)
	}

	elapsed := time.Since(start)
	return alloc, elapsed
}

type Task = allocator.Task
type Server = allocator.Server

func main() {
	const seed = 42
	cfg := config{NTasks: 50, NServers: 5, CPUCap: 16, RAMCap: 32}

	rng := rand.New(rand.NewSource(seed))
	tasks := allocator.GenerateTasks(cfg.NTasks, rng)
	servers := allocator.GenerateServers(cfg.NServers, cfg.CPUCap, cfg.RAMCap)

	algos := []string{"Greedy", "DP", "Heuristic"}
	results := make(map[string]algoResult)

	fmt.Println("======================================================================")
	fmt.Println("  CLOUD RESOURCE ALLOCATION — Go BENCHMARK")
	fmt.Printf("  %d tasks  |  %d servers  |  %d CPU / %d GB RAM each\n",
		cfg.NTasks, cfg.NServers, cfg.CPUCap, cfg.RAMCap)
	fmt.Println("======================================================================")
	fmt.Println()

	// Warm up: run each once to avoid cold-start effects
	for _, name := range algos {
		runAlgo(name, tasks, servers)
	}

	// Actual timed run (average of 100 iterations for greedy/heuristic, 5 for DP)
	for _, name := range algos {
		iters := 100
		if name == "DP" {
			iters = 5
		}
		var totalDur time.Duration
		var lastAlloc allocator.Allocation
		for i := 0; i < iters; i++ {
			a, d := runAlgo(name, tasks, servers)
			totalDur += d
			lastAlloc = a
		}
		avgMs := float64(totalDur.Nanoseconds()) / float64(iters) / 1e6

		freshServers := allocator.GenerateServers(cfg.NServers, cfg.CPUCap, cfg.RAMCap)
		m := allocator.Evaluate(lastAlloc, tasks, freshServers)
		m.TimeMs = avgMs
		results[name] = algoResult{m}
	}

	// Print table
	hdr := fmt.Sprintf("  %-28s %10s %10s %10s", "Metric", "Greedy", "DP", "Heuristic")
	fmt.Println(hdr)
	fmt.Println("  " + repeatDash(len(hdr)-2))

	type row struct{ label, key string }
	rows := []row{
		{"Tasks placed", "placed"},
		{"Priority score", "priority_score"},
		{"Avg CPU utilisation (%)", "avg_cpu_util"},
		{"Avg RAM utilisation (%)", "avg_ram_util"},
		{"CPU load imbalance (%)", "cpu_imbalance"},
		{"Execution time (ms)", "time_ms"},
	}
	for _, r := range rows {
		g := results["Greedy"]
		d := results["DP"]
		h := results["Heuristic"]
		fmt.Printf("  %-28s %10s %10s %10s\n",
			r.label, fmtVal(g, r.key), fmtVal(d, r.key), fmtVal(h, r.key))
	}
	fmt.Println()

	// ── Scalability ──
	scales := []int{20, 50, 100, 200, 500, 1000}
	scaleResults := make(map[string]map[string]scaleEntry)

	for _, n := range scales {
		m := max(3, n/10)
		rng2 := rand.New(rand.NewSource(seed))
		t2 := allocator.GenerateTasks(n, rng2)
		s2 := allocator.GenerateServers(m, cfg.CPUCap, cfg.RAMCap)

		entry := make(map[string]scaleEntry)
		for _, name := range algos {
			iters := 50
			if name == "DP" {
				iters = 3
			}
			if name == "DP" && n > 500 {
				iters = 1
			}
			var totalDur time.Duration
			var lastAlloc allocator.Allocation
			for i := 0; i < iters; i++ {
				a, d := runAlgo(name, t2, s2)
				totalDur += d
				lastAlloc = a
			}
			avgMs := float64(totalDur.Nanoseconds()) / float64(iters) / 1e6

			placed := 0
			for _, sid := range lastAlloc {
				if sid != "" {
					placed++
				}
			}
			entry[name] = scaleEntry{TimeMs: avgMs, Placed: placed, Total: n}
		}
		scaleResults[fmt.Sprintf("%d", n)] = entry
	}

	// ── Write JSON ──
	outDir := filepath.Join("..", "results")
	os.MkdirAll(outDir, 0755)
	outPath := filepath.Join(outDir, "benchmark.json")

	out := output{
		Primary:     primaryOutput{Config: cfg, Results: results},
		Scalability: scaleResults,
	}

	data, _ := json.MarshalIndent(out, "", "  ")
	os.WriteFile(outPath, data, 0644)
	fmt.Printf("Results saved to %s\n", outPath)
}

func fmtVal(r algoResult, key string) string {
	switch key {
	case "placed":
		return fmt.Sprintf("%d", r.Placed)
	case "priority_score":
		return fmt.Sprintf("%d", r.PriorityScore)
	case "avg_cpu_util":
		return fmt.Sprintf("%.1f", r.AvgCPUUtil)
	case "avg_ram_util":
		return fmt.Sprintf("%.1f", r.AvgRAMUtil)
	case "cpu_imbalance":
		return fmt.Sprintf("%.1f", r.CPUImbalance)
	case "time_ms":
		return fmt.Sprintf("%.4f", r.TimeMs)
	}
	return ""
}

func repeatDash(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = '-'
	}
	return string(b)
}
