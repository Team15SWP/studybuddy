package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"study_buddy/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
)

const (
	appTimeout     = 60 * time.Second
	maxHeaderBytes = 1 << 20
)

type Server struct {
	cfg       *config.Config
	log       *slog.Logger
	router    *gin.Engine
	pgPool    *pgxpool.Pool
	repoLayer *repoLayer
}

func NewServer(cfg *config.Config, log *slog.Logger, pgPool *pgxpool.Pool) *Server {
	router := gin.Default()
	return &Server{
		cfg:       cfg,
		log:       log,
		router:    router,
		pgPool:    pgPool,
		repoLayer: initRepoLayer(pgPool),
	}
}

func (s *Server) Run(ctx context.Context) {
	s.setupMiddlewares()
	s.mapHandlers()
	httpServer := &http.Server{
		Addr:           ":" + s.cfg.RouterConfig.AppPort,
		Handler:        s.router,
		ReadTimeout:    appTimeout,
		WriteTimeout:   appTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		slog.Info(fmt.Sprintf("starting server on port %s...", s.cfg.RouterConfig.AppPort))
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		slog.Info("closing PGConfig pool")
		slog.Info("shutting down the http server...")
		s.pgPool.Close()
		return httpServer.Shutdown(ctx)
	})

	if err := g.Wait(); err != nil {
		slog.Info(fmt.Sprintf("exit reason: %s", err))
	}
}
