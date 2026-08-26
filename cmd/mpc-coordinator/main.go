package main

import (
	"context"
	"errors"
	"github.com/example/027-mpc-coordinator/internal/application"
	"github.com/example/027-mpc-coordinator/internal/crypto"
	"github.com/example/027-mpc-coordinator/internal/repository"
	"github.com/example/027-mpc-coordinator/internal/transport"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	addr := os.Getenv("MPC_HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	store := repository.NewMemoryStore(os.Getenv("MPC_STATE_FILE"))
	svc := application.NewService(store, crypto.NewKeyProvider())
	srv := &http.Server{Addr: addr, Handler: transport.NewServer(svc, logger).Handler(), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		logger.Info("server_started", "addr", addr)
		if e := srv.ListenAndServe(); e != nil && !errors.Is(e, http.ErrServerClosed) {
			logger.Error("server_failed", "error", e)
			os.Exit(1)
		}
	}()
	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdown)
	logger.Info("server_stopped")
}
