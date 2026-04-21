package demo

import (
	"fmt"
	"os"
	"strings"

	"github.com/ming-hsien/lang-chain-template/internal/tools"
	"github.com/tmc/langchaingo/llms"
)

func init() {
	tools.Register(tools.Tool{
		Name:        "GetSystemStats",
		Description: "Retrieves current system memory and CPU load information.",
		Schema: llms.Tool{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "GetSystemStats",
				Description: "Retrieves current system memory and CPU load information.",
				Parameters: map[string]any{
					"type":       "object",
					"properties": map[string]any{},
				},
			},
		},
		Execute: func(args string) string {
			return getSystemStats()
		},
	})
}

func getSystemStats() string {
	memInfo, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return fmt.Sprintf("Error reading memory info: %v", err)
	}

	loadAvg, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return fmt.Sprintf("Error reading load average: %v", err)
	}

	lines := strings.Split(string(memInfo), "\n")
	var memTotal, memAvailable string
	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal:") {
			memTotal = strings.TrimSpace(strings.TrimPrefix(line, "MemTotal:"))
		}
		if strings.HasPrefix(line, "MemAvailable:") {
			memAvailable = strings.TrimSpace(strings.TrimPrefix(line, "MemAvailable:"))
		}
	}

	loadParts := strings.Fields(string(loadAvg))
	var loadStr string
	if len(loadParts) >= 3 {
		loadStr = fmt.Sprintf("1m:%s, 5m:%s, 15m:%s", loadParts[0], loadParts[1], loadParts[2])
	} else {
		loadStr = "Unknown"
	}

	return fmt.Sprintf("System Status:\n- Memory Total: %s\n- Memory Available: %s\n- Load Avg: %s",
		memTotal, memAvailable, loadStr)
}
