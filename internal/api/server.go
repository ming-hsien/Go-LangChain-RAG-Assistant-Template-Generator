package api

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ming-hsien/lang-chain-template/internal/config"
	"github.com/ming-hsien/lang-chain-template/internal/rag"
)

//go:embed frontend/*
var staticFiles embed.FS

type Server struct {
	router *gin.Engine
	ragSvc *rag.RAGService
}

func NewServer(ragSvc *rag.RAGService) *Server {
	r := gin.Default()

	s := &Server{
		router: r,
		ragSvc: ragSvc,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// API Routes
	v1 := s.router.Group("/v1")
	{
		v1.POST("/query", s.handleQuery)
		v1.GET("/tools", s.handleGetTools)
		v1.POST("/index", s.handleReindex)
		v1.GET("/config", s.handleGetConfig)
	}

	// Static UI
	sub, _ := fs.Sub(staticFiles, "frontend")
	s.router.StaticFS("/ui", http.FS(sub))
	s.router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/ui/")
	})
}

func (s *Server) handleQuery(c *gin.Context) {
	var req struct {
		Question string `json:"question" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := s.ragSvc.Query(c.Request.Context(), req.Question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"answer": resp})
}

func (s *Server) handleGetTools(c *gin.Context) {
	// Mock tool registry for future AI Agent
	tools := []map[string]interface{}{
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

	c.JSON(http.StatusOK, gin.H{"tools": tools})
}

func (s *Server) handleReindex(c *gin.Context) {
	go func() {
		err := s.ragSvc.IndexDocuments(context.Background(), "./documents")
		if err != nil {
			log.Printf("Reindex failed: %v", err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Reindexing started in background"})
}

func (s *Server) handleGetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"app_name": config.AppConfig.AppName,
	})
}

func (s *Server) Run() error {
	addr := ":" + config.AppConfig.ServerPort
	return s.router.Run(addr)
}
