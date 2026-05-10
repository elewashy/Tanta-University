package allocator

import "math"

// Metrics holds evaluation results for one algorithm run.
type Metrics struct {
	Placed        int                `json:"placed"`
	Total         int                `json:"total"`
	PriorityScore int                `json:"priority_score"`
	AvgCPUUtil    float64            `json:"avg_cpu_util"`
	AvgRAMUtil    float64            `json:"avg_ram_util"`
	CPUImbalance  float64            `json:"cpu_imbalance"`
	PerServerCPU  map[string]float64 `json:"per_server_cpu"`
	PerServerRAM  map[string]float64 `json:"per_server_ram"`
	TimeMs        float64            `json:"time_ms"`
}

// Evaluate computes metrics for an allocation against the original
// task list and server definitions (capacities only; remaining is ignored).
func Evaluate(alloc Allocation, tasks []Task, servers []Server) Metrics {
	taskMap := make(map[string]*Task, len(tasks))
	for i := range tasks {
		taskMap[tasks[i].ID] = &tasks[i]
	}

	placed := 0
	priorityScore := 0
	cpuUsed := make(map[string]int, len(servers))
	ramUsed := make(map[string]int, len(servers))

	for _, s := range servers {
		cpuUsed[s.ID] = 0
		ramUsed[s.ID] = 0
	}

	for tid, sid := range alloc {
		if sid == "" {
			continue
		}
		placed++
		t := taskMap[tid]
		priorityScore += t.Priority
		cpuUsed[sid] += t.CPU
		ramUsed[sid] += t.RAM
	}

	cpuUtil := make(map[string]float64, len(servers))
	ramUtil := make(map[string]float64, len(servers))
	var sumCPU, sumRAM float64
	minCPU, maxCPU := 100.0, 0.0

	for _, s := range servers {
		cu := float64(cpuUsed[s.ID]) / float64(s.CPUTotal) * 100
		ru := float64(ramUsed[s.ID]) / float64(s.RAMTotal) * 100
		cpuUtil[s.ID] = math.Round(cu*10) / 10
		ramUtil[s.ID] = math.Round(ru*10) / 10
		sumCPU += cu
		sumRAM += ru
		if cu < minCPU {
			minCPU = cu
		}
		if cu > maxCPU {
			maxCPU = cu
		}
	}

	n := float64(len(servers))
	return Metrics{
		Placed:        placed,
		Total:         len(tasks),
		PriorityScore: priorityScore,
		AvgCPUUtil:    math.Round(sumCPU/n*10) / 10,
		AvgRAMUtil:    math.Round(sumRAM/n*10) / 10,
		CPUImbalance:  math.Round((maxCPU-minCPU)*10) / 10,
		PerServerCPU:  cpuUtil,
		PerServerRAM:  ramUtil,
	}
}
