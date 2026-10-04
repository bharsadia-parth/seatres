package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bharsadia-parth/seatres/internal/api"
	"github.com/bharsadia-parth/seatres/internal/store"
	"github.com/bharsadia-parth/seatres/migrations"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		slog.Error("DATABASE_URL is required")
		os.Exit(1)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// ctx is cancelled when we receive Ctrl+C or SIGTERM (what Docker/Render sends).
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	st, err := store.New(ctx, dsn)
	if err != nil {
		slog.Error("store init failed", "err", err)
		os.Exit(1)
	}

	defer st.Close()

	if st.Migrate(ctx, migrations.FS); err != nil {
		slog.Error("migrate failed", "err", err)
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           api.NewRouter(st),
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Run the server in its own goroutine so main can wait for the shutdown signal.
	go func() {
		slog.Info("listening", "port", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "err", err)
			stop()
		}
	}()

	<-ctx.Done() // block until signal
	slog.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx) // finish in-flight requests, then exit

}
