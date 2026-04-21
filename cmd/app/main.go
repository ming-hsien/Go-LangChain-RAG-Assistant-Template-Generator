package main

import (
	"log"

	"github.com/ming-hsien/lang-chain-template/internal/api"
	"github.com/ming-hsien/lang-chain-template/internal/config"
	"github.com/ming-hsien/lang-chain-template/internal/rag"
)

func main() {
	// 1. Load Config from .env or Env Vars
	config.LoadConfig()

	// 2. Initialize RAG Service (GitHub Models + Qdrant)
	ragSvc, err := rag.NewRAGService()
	if err != nil {
		log.Fatalf("Critical: Failed to initialize RAG Service: %v", err)
	}

	// 3. Initialize API Server with the RAG Service
	server := api.NewServer(ragSvc)

	// 4. Start the Server
	log.Printf("Core Network CMS AI Agent starting on port %s...", config.AppConfig.ServerPort)
	if err := server.Run(); err != nil {
		log.Fatalf("Critical: Server failed to run: %v", err)
	}
}
