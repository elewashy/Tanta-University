package allocator

import "math/rand"

// GenerateTasks creates n tasks with random CPU (1-8), RAM (1-16), priority (1-10).
func GenerateTasks(n int, rng *rand.Rand) []Task {
	tasks := make([]Task, n)
	for i := range tasks {
		tasks[i] = Task{
			ID:       "T" + itoa(i+1),
			CPU:      rng.Intn(8) + 1,
			RAM:      rng.Intn(16) + 1,
			Priority: rng.Intn(10) + 1,
		}
	}
	return tasks
}

// GenerateServers creates m identical servers.
func GenerateServers(m, cpuCap, ramCap int) []Server {
	servers := make([]Server, m)
	for i := range servers {
		servers[i] = Server{
			ID:           "S" + itoa(i+1),
			CPUTotal:     cpuCap,
			RAMTotal:     ramCap,
			CPURemaining: cpuCap,
			RAMRemaining: ramCap,
		}
	}
	return servers
}

// itoa is a simple int-to-string without importing strconv.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
