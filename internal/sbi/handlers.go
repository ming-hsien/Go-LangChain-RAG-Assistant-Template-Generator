package sbi

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) handleQuery(c *gin.Context) {
	var req struct {
		Question string `json:"question" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := s.processor.Ask(c.Request.Context(), req.Question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"answer": resp})
}

func (s *Server) handleGetTools(c *gin.Context) {
	tools := s.processor.GetTools()
	c.JSON(http.StatusOK, gin.H{"tools": tools})
}

func (s *Server) handleReindex(c *gin.Context) {
	go func() {
		err := s.processor.Reindex(context.Background())
		if err != nil {
			log.Printf("Reindex failed: %v", err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Reindexing started in background"})
}

func (s *Server) handleGetConfig(c *gin.Context) {
	info := s.processor.GetInfo()
	c.JSON(http.StatusOK, info)
}
