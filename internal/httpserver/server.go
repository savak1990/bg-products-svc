package httpserver

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	p "github.com/savak1990/bg-products-svc/internal/products"
)

type Server struct {
	engine *gin.Engine
	repo   p.Repo
}

func NewServer(repo p.Repo) *Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.RecoveryWithWriter(os.Stderr))

	s := &Server{engine: r, repo: repo}

	v1 := r.Group("/v1")
	{
		v1.GET("/products", s.onGetProducts)
		v1.POST("/products", s.onPostProducts)
	}

	r.GET("/healthz/live", s.onGetHealthzLive)
	r.GET("/healthz/ready", s.onGetHealthzReady)
	r.GET("/health", s.onGetHealthzLive)

	return s
}

func (s *Server) onGetProducts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"items": s.repo.List()})
}

func (s *Server) onPostProducts(c *gin.Context) {
	var in p.Product
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p, err := s.repo.Create(in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, p)
}

func (s *Server) onGetHealthzLive(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (s *Server) onGetHealthzReady(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (s *Server) Handler() http.Handler {
	return s.engine
}
