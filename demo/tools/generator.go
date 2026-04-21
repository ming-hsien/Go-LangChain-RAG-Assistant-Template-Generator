package demo

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ming-hsien/lang-chain-template/internal/tools"
	"github.com/tmc/langchaingo/llms"
)

func init() {
	tools.Register(tools.Tool{
		Name:        "CreateDocument",
		Description: "Generates a new text document in the local knowledge base.",
		Schema: llms.Tool{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "CreateDocument",
				Description: "Generates a new text document in the local knowledge base.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"filename": map[string]any{
							"type":        "string",
							"description": "The name of the file (e.g., 'report.txt')",
						},
						"content": map[string]any{
							"type":        "string",
							"description": "The text content to write into the file",
						},
					},
					"required": []string{"filename", "content"},
				},
			},
		},
		Execute: func(args string) string {
			var params struct {
				Filename string `json:"filename"`
				Content  string `json:"content"`
			}
			if err := tools.ParseArgs(args, &params); err != nil {
				return fmt.Sprintf("Error parsing arguments for CreateDocument: %v", err)
			}
			return createDocument(params.Filename, params.Content)
		},
	})
}

func createDocument(filename, content string) string {
	if filepath.Ext(filename) != ".txt" {
		filename = filename + ".txt"
	}

	filePath := filepath.Join("documents", filename)

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Sprintf("Error creating document %s: %v", filename, err)
	}

	return fmt.Sprintf("Successfully generated document: %s. Re-index can make it available to the AI.", filename)
}
