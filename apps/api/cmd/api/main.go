package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/application"
	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/domain"
	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/infrastructure/memory"
	mysqlrepo "github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/infrastructure/mysql"
	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/dataset/transport/httpapi"
	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/platform/config"
	"github.com/dantetest/GPT_Blctek_IP_Chain_715/apps/api/internal/platform/httpx"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.LoadAPI()
	if err != nil {
		logger.Error("configuration rejected", "error", err)
		os.Exit(1)
	}

	repository, closeRepository, err := buildDatasetRepository(cfg)
	if err != nil {
		logger.Error("dataset repository initialization failed", "backend", cfg.DatasetRepository, "error", err)
		os.Exit(1)
	}
	defer closeRepository()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{Code: "HEALTHY", Message: "service is healthy", Data: map[string]string{"service": "api"}})
	})
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, _ *http.Request) {
		httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{Code: "READY", Message: "service is ready", Data: map[string]string{"dataset_repository": cfg.DatasetRepository}})
	})

	service := application.NewService(repository, application.RandomIDGenerator{}, nil)
	httpapi.New(service).Register(mux)

	server := &http.Server{
		Addr:              cfg.Address,
		Handler:           httpx.RequestID(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		logger.Info("api listening", "address", cfg.Address, "environment", cfg.Environment, "dataset_repository", cfg.DatasetRepository)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("api terminated unexpectedly", "error", err)
			stop()
		}
	}()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}
	logger.Info("api stopped")
}

func buildDatasetRepository(cfg config.API) (domain.Repository, func(), error) {
	if cfg.DatasetRepository == config.DatasetRepositoryMemory {
		return memory.NewRepository(), func() {}, nil
	}
	db, err := sql.Open("mysql", cfg.DatabaseDSN)
	if err != nil {
		return nil, func() {}, err
	}
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, func() {}, err
	}
	return mysqlrepo.New(db), func() { _ = db.Close() }, nil
}
