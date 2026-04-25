package promptmgr

import _ "embed"

//go:embed system.txt
var systemPrompt string

// GetSystemPrompt returns the embedded system prompt content
func GetSystemPrompt() string {
	return systemPrompt
}
