package processor

import (
	"context"

	"github.com/ming-hsien/lang-chain-template/internal/config"
	"github.com/ming-hsien/lang-chain-template/internal/rag"
)

type Processor struct {
	ragSvc *rag.RAGService
}

func NewProcessor(ragSvc *rag.RAGService) *Processor {
	return &Processor{
		ragSvc: ragSvc,
	}
}

// Ask handles the core RAG questioning logic
func (p *Processor) Ask(ctx context.Context, question string) (string, error) {
	return p.ragSvc.Query(ctx, question)
}

// Reindex triggers the background document indexing
func (p *Processor) Reindex(ctx context.Context) error {
	// Running in a background context to avoid being killed by a request timeout
	// if called from an API, though the caller is responsible for the exact ctx type.
	return p.ragSvc.IndexDocuments(ctx, "./documents")
}

// GetInfo returns the application configuration info
func (p *Processor) GetInfo() map[string]interface{} {
	return map[string]interface{}{
		"app_name": config.AppConfig.AppName,
	}
}

// GetTools returns the available tools registry
func (p *Processor) GetTools() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":        "ResetNode",
			"description": "Resets a specific core network node.",
			"parameters": []map[string]string{
				{"name": "node_id", "type": "string", "description": "The unique ID of the node"},
			},
		},
		{
			"name":        "GetTopology",
			"description": "Retrieves the current network topology mapping.",
			"parameters":  []map[string]string{},
		},
	}
}
