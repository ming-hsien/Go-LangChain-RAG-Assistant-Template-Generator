package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ming-hsien/lang-chain-template/internal/config"
	"github.com/ming-hsien/lang-chain-template/internal/sbi/processor"
	"github.com/ming-hsien/lang-chain-template/web"
)

type Server struct {
	router    *gin.Engine
	processor *processor.Processor
}

func NewServer(proc *processor.Processor) *Server {
	r := gin.Default()

	s := &Server{
		router:    r,
		processor: proc,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	v1 := s.router.Group("/v1")
	{
		v1.POST("/query", s.handleQuery)
		v1.DELETE("/history", s.handleClearHistory)
		v1.GET("/tools", s.handleGetTools)
		v1.POST("/index", s.handleReindex)
		v1.GET("/config", s.handleGetConfig)
	}

	// Static UI serving from the 'web' package
	s.router.StaticFS("/ui", http.FS(web.Assets))
	s.router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/ui/")
	})
}

func (s *Server) Run() error {
	addr := ":" + config.AppConfig.ServerPort
	return s.router.Run(addr)
}
