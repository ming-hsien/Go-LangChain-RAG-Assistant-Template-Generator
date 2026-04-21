package tools

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/tmc/langchaingo/llms"
)

// Tool represents a pluggable capability for the AI Agent
type Tool struct {
	Name        string
	Description string
	Schema      llms.Tool
	Execute     func(args string) string
}

var (
	registry = make(map[string]Tool)
	mu       sync.RWMutex
)

// Register adds a new tool to the global registry
func Register(t Tool) {
	mu.Lock()
	defer mu.Unlock()
	registry[t.Name] = t
	log.Printf("Tool registered: %s", t.Name)
}

// GetDefinitions returns the list of tools the AI can use (from the registry)
func GetDefinitions() []llms.Tool {
	mu.RLock()
	defer mu.RUnlock()

	defs := make([]llms.Tool, 0, len(registry))
	for _, t := range registry {
		defs = append(defs, t.Schema)
	}
	return defs
}

// Execute is the dispatcher for all registered tools
func Execute(name, args string) string {
	mu.RLock()
	tool, exists := registry[name]
	mu.RUnlock()

	if !exists {
		return fmt.Sprintf("Tool %s not found in registry", name)
	}

	log.Printf("Executing tool: %s", name)
	return tool.Execute(args)
}

// Helper to wrap common JSON parsing for tools
func ParseArgs(args string, target interface{}) error {
	return json.Unmarshal([]byte(args), target)
}
