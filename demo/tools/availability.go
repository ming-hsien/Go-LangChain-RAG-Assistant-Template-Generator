package demo

import (
	"fmt"

	"github.com/ming-hsien/lang-chain-template/internal/service"
	"github.com/ming-hsien/lang-chain-template/internal/tools"
	"github.com/tmc/langchaingo/llms"
)

var employeeSvc = service.NewEmployeeService()

func init() {
	tools.Register(tools.Tool{
		Name:        "CheckAvailability",
		Description: "Checks the real-time availability and calendar status of a specific employee.",
		Schema: llms.Tool{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "CheckAvailability",
				Description: "Checks the real-time availability and calendar status of a specific employee.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"name": map[string]any{
							"type":        "string",
							"description": "The name of the employee to check (e.g., 'Carol', 'Alonza')",
						},
					},
					"required": []string{"name"},
				},
			},
		},
		Execute: checkAvailability,
	})
}

func checkAvailability(args string) string {
	var params struct {
		Name string `json:"name"`
	}
	if err := tools.ParseArgs(args, &params); err != nil {
		return fmt.Sprintf("Error parsing arguments for CheckAvailability: %v", err)
	}
	return employeeSvc.CheckAvailability(params.Name)
}
