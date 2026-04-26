package processor

import (
	"context"
	"fmt"
	"log"

	"github.com/ming-hsien/lang-chain-template/internal/config"
	"github.com/ming-hsien/lang-chain-template/internal/history"
	"github.com/ming-hsien/lang-chain-template/internal/promptmgr"
	"github.com/ming-hsien/lang-chain-template/internal/rag"
	"github.com/ming-hsien/lang-chain-template/internal/tools"
	"github.com/tmc/langchaingo/llms"
)

type Processor struct {
	ragSvc  *rag.RAGService
	history *history.HistoryManager
}

func NewProcessor(ragSvc *rag.RAGService) *Processor {
	return &Processor{
		ragSvc:  ragSvc,
		history: history.NewHistoryManager(),
	}
}

func (p *Processor) Ask(ctx context.Context, sessionID, question string) (string, error) {
	ragContext, err := p.ragSvc.Search(ctx, question)
	if err != nil {
		log.Printf("RAG search failed: %v", err)
		ragContext = "No relevant documents found."
	}

	var messages []llms.MessageContent
	if sessionID != "" {
		messages = p.history.Get(sessionID)
	}

	systemPrompt := fmt.Sprintf("%s\n\nContext:\n%s", promptmgr.GetSystemPrompt(), ragContext)

	// We construct a full list for the LLM: [System] + [History] + [Current Question]
	fullMessages := []llms.MessageContent{
		{
			Role:  llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{llms.TextPart(systemPrompt)},
		},
	}
	fullMessages = append(fullMessages, messages...)
	fullMessages = append(fullMessages, llms.MessageContent{
		Role:  llms.ChatMessageTypeHuman,
		Parts: []llms.ContentPart{llms.TextPart(question)},
	})

	toolDefs := tools.GetDefinitions()

	for i := 0; i < 5; i++ {
		resp, err := p.ragSvc.LLM.GenerateContent(ctx, fullMessages, llms.WithTools(toolDefs))
		if err != nil {
			return "", fmt.Errorf("llm generation failed: %v", err)
		}

		choice := resp.Choices[0]
		if len(choice.ToolCalls) == 0 {
			finalAnswer := choice.Content

			if sessionID != "" {
				messages = append(messages, llms.MessageContent{
					Role:  llms.ChatMessageTypeHuman,
					Parts: []llms.ContentPart{llms.TextPart(question)},
				})
				messages = append(messages, llms.MessageContent{
					Role:  llms.ChatMessageTypeAI,
					Parts: []llms.ContentPart{llms.TextPart(finalAnswer)},
				})

				if len(messages) > 10 {
					messages = messages[len(messages)-10:]
				}

				p.history.Set(sessionID, messages)
			}

			return finalAnswer, nil
		}

		// Tool calls needed
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
		fullMessages = append(fullMessages, assistantMsg)

		for _, tc := range choice.ToolCalls {
			result := tools.Execute(tc.FunctionCall.Name, tc.FunctionCall.Arguments)

			fullMessages = append(fullMessages, llms.MessageContent{
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

func (p *Processor) ClearHistory(sessionID string) {
	p.history.Clear(sessionID)
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
