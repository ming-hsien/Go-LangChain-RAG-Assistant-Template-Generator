package main

import (
	"log"

	"github.com/ming-hsien/lang-chain-template/internal/config"
	"github.com/ming-hsien/lang-chain-template/internal/rag"
	"github.com/ming-hsien/lang-chain-template/internal/sbi"
	"github.com/ming-hsien/lang-chain-template/internal/sbi/processor"
)

func main() {
	config.LoadConfig()

	ragSvc, err := rag.NewRAGService()
	if err != nil {
		log.Fatalf("Critical: Failed to initialize RAG Service: %v", err)
	}

	proc := processor.NewProcessor(ragSvc)
	server := sbi.NewServer(proc)

	log.Printf("Core Network CMS AI Agent starting on port %s...", config.AppConfig.ServerPort)
	if err := server.Run(); err != nil {
		log.Fatalf("Critical: Server failed to run: %v", err)
	}
}
