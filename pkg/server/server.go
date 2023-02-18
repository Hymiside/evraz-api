package server

import (
	"context"
	"fmt"
	"github.com/Hymiside/evraz-api/pkg/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Server struct{}

func (s *Server) RunServer(ctx context.Context, handler *gin.Engine, c models.ConfigServer) error {
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", c.Host, c.Port),
		Handler: handler,
	}

	go func(ctx context.Context) {
		<-ctx.Done()
		httpServer.Shutdown(ctx)
	}(ctx)

	log.Printf("authentication microservice launched on http://%s:%s/", c.Host, c.Port)
	return httpServer.ListenAndServe()
}
