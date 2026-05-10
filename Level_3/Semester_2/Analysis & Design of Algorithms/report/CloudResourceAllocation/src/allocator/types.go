// Package allocator implements three cloud resource allocation algorithms:
// Greedy First-Fit, Dynamic Programming (2-D Knapsack), and Heuristic Load
// Balancing.  Each algorithm assigns tasks with CPU, RAM, and priority
// requirements to a pool of servers with finite capacity.
package allocator

// Task represents a workload request with resource requirements.
type Task struct {
	ID       string
	CPU      int
	RAM      int
	Priority int
}

// Server represents a compute node with remaining capacity.
type Server struct {
	ID           string
	CPUTotal     int
	RAMTotal     int
	CPURemaining int
	RAMRemaining int
}

// Allocation maps a Task ID to a Server ID (empty string = unplaced).
type Allocation map[string]string

// CloneServers returns a deep copy so algorithms can mutate freely.
func CloneServers(servers []Server) []Server {
	out := make([]Server, len(servers))
	copy(out, servers)
	return out
}
