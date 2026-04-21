package processor

import (
	"context"
	"fmt"
	"log"

	"github.com/ming-hsien/lang-chain-template/internal/config"
	"github.com/ming-hsien/lang-chain-template/internal/rag"
	"github.com/ming-hsien/lang-chain-template/internal/tools"
	"github.com/tmc/langchaingo/llms"
)

type Processor struct {
	ragSvc *rag.RAGService
}

func NewProcessor(ragSvc *rag.RAGService) *Processor {
	return &Processor{
		ragSvc: ragSvc,
	}
}

func (p *Processor) Ask(ctx context.Context, question string) (string, error) {
	ragContext, err := p.ragSvc.Search(ctx, question)
	if err != nil {
		log.Printf("RAG search failed: %v", err)
		ragContext = "No relevant documents found."
	}

	// 2. Prepare message history
	systemPrompt := fmt.Sprintf("%s\n\nContext:\n%s", config.AppConfig.SystemPrompt, ragContext)
	messages := []llms.MessageContent{
		{
			Role:  llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{llms.TextPart(systemPrompt)},
		},
		{
			Role:  llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{llms.TextPart(question)},
		},
	}

	toolDefs := tools.GetDefinitions()

	for i := 0; i < 5; i++ {
		resp, err := p.ragSvc.LLM.GenerateContent(ctx, messages, llms.WithTools(toolDefs))
		if err != nil {
			return "", fmt.Errorf("llm generation failed: %v", err)
		}

		choice := resp.Choices[0]
		if len(choice.ToolCalls) == 0 {
			return choice.Content, nil
		}

		assistantMsg := llms.MessageContent{
			Role: llms.ChatMessageTypeAI,
		}
		for _, tc := range choice.ToolCalls {
			assistantMsg.Parts = append(assistantMsg.Parts, llms.ToolCall{
				ID:   tc.ID,
				Type: tc.Type,
				FunctionCall: &llms.FunctionCall{
					Name:      tc.FunctionCall.Name,
					Arguments: tc.FunctionCall.Arguments,
				},
			})
		}
		messages = append(messages, assistantMsg)

		for _, tc := range choice.ToolCalls {
			result := tools.Execute(tc.FunctionCall.Name, tc.FunctionCall.Arguments)

			messages = append(messages, llms.MessageContent{
				Role: llms.ChatMessageTypeTool,
				Parts: []llms.ContentPart{
					llms.ToolCallResponse{
						ToolCallID: tc.ID,
						Name:       tc.FunctionCall.Name,
						Content:    result,
					},
				},
			})
		}
	}

	return "", fmt.Errorf("agent reached max iterations without final answer")
}

func (p *Processor) Reindex(ctx context.Context) error {
	return p.ragSvc.IndexDocuments(ctx, "./documents")
}

func (p *Processor) GetInfo() map[string]interface{} {
	return map[string]interface{}{
		"app_name": config.AppConfig.AppName,
	}
}

func (p *Processor) GetTools() []map[string]interface{} {
	defs := tools.GetDefinitions()
	results := make([]map[string]interface{}, 0)
	for _, d := range defs {
		results = append(results, map[string]interface{}{
			"name":        d.Function.Name,
			"description": d.Function.Description,
		})
	}
	return results
}
