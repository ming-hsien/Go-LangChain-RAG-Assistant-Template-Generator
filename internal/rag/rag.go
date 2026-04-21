package rag

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/ming-hsien/lang-chain-template/internal/config"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/qdrant"
)

type RAGService struct {
	LLM         llms.Model
	VectorStore vectorstores.VectorStore
}

func NewRAGService() (*RAGService, error) {
	llm, err := openai.New(
		openai.WithBaseURL(config.AppConfig.LLMBaseURL),
		openai.WithToken(config.AppConfig.GitHubToken),
		openai.WithModel(config.AppConfig.LLMModel),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init llm: %v", err)
	}

	embedClient, err := openai.New(
		openai.WithBaseURL(config.AppConfig.LLMBaseURL),
		openai.WithToken(config.AppConfig.GitHubToken),
		openai.WithEmbeddingModel("text-embedding-3-small"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init embed client: %v", err)
	}

	embedder, err := embeddings.NewEmbedder(embedClient)
	if err != nil {
		return nil, fmt.Errorf("failed to init embedder: %v", err)
	}

	// 3. Parse Qdrant URL
	qURL, err := url.Parse(config.AppConfig.QdrantURL)
	if err != nil {
		return nil, fmt.Errorf("invalid qdrant url: %v", err)
	}

	// 4. Ensure Qdrant Collection exists
	err = ensureCollection(config.AppConfig.QdrantURL, config.AppConfig.CollectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure collection: %v", err)
	}

	// 5. Initialize Qdrant
	store, err := qdrant.New(
		qdrant.WithURL(*qURL),
		qdrant.WithCollectionName(config.AppConfig.CollectionName),
		qdrant.WithEmbedder(embedder),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init qdrant: %v", err)
	}

	return &RAGService{
		LLM:         llm,
		VectorStore: store,
	}, nil
}

func (s *RAGService) IndexDocuments(ctx context.Context, dirPath string) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	splitter := textsplitter.NewRecursiveCharacter()
	splitter.ChunkSize = 1000
	splitter.ChunkOverlap = 200

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(dirPath, entry.Name())
		log.Printf("Indexing document: %s", filePath)

		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Error reading file %s: %v", filePath, err)
			continue
		}

		doc := schema.Document{
			PageContent: string(data),
			Metadata: map[string]interface{}{
				"source": entry.Name(),
			},
		}

		// Split document into chunks
		chunks, err := textsplitter.SplitDocuments(splitter, []schema.Document{doc})
		if err != nil {
			log.Printf("Error splitting file %s: %v", filePath, err)
			continue
		}

		_, err = s.VectorStore.AddDocuments(ctx, chunks)
		if err != nil {
			log.Printf("Error adding documents to vector store: %v", err)
			continue
		}
	}

	return nil
}

func (s *RAGService) Query(ctx context.Context, question string) (string, error) {
	// 1. Retrieve similar documents
	docs, err := s.VectorStore.SimilaritySearch(ctx, question, 3)
	if err != nil {
		return "", fmt.Errorf("failed to search vector store: %v", err)
	}

	// 2. Build context
	contextText := ""
	for _, doc := range docs {
		contextText += doc.PageContent + "\n---\n"
	}

	// 3. Simple QA Prompt
	prompt := fmt.Sprintf(`%s
	
Context:
%s

Question:
%s

Answer:`, config.AppConfig.SystemPrompt, contextText, question)

	// 4. Generate response
	resp, err := llms.GenerateFromSinglePrompt(ctx, s.LLM, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate response: %v", err)
	}

	return resp, nil
}

func ensureCollection(qdrantURL string, collectionName string) error {
	// Check if collection exists
	checkURL := fmt.Sprintf("%s/collections/%s", qdrantURL, collectionName)
	resp, err := http.Get(checkURL)
	if err != nil {
		return fmt.Errorf("failed to check collection: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Printf("Collection '%s' already exists", collectionName)
		return nil
	}

	createURL := fmt.Sprintf("%s/collections/%s", qdrantURL, collectionName)
	config := map[string]interface{}{
		"vectors": map[string]interface{}{
			"size":     1536,
			"distance": "Cosine",
		},
	}
	jsonBody, _ := json.Marshal(config)
	req, _ := http.NewRequest("PUT", createURL, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create collection: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create collection, status: %s", resp.Status)
	}

	log.Printf("Successfully created collection '%s'", collectionName)
	return nil
}
