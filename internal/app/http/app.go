package httpapp

import (
	"context"
	"fmt"
	"github.com/ARUMANDESU/uniclubs-filestorage-service/internal/config"
	"log/slog"
	"net/http"
)

type App struct {
	log        *slog.Logger
	httpServer *http.Server
}

func New(cfg *config.Config, log *slog.Logger, handler http.Handler) *App {
	httpServer := &http.Server{
		Addr:           cfg.HTTP.Address,
		Handler:        handler,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		ReadTimeout:    cfg.HTTP.Timeout,
		WriteTimeout:   cfg.HTTP.Timeout,
		IdleTimeout:    cfg.HTTP.IdleTimeout,
	}

	return &App{httpServer: httpServer, log: log}
}

// MustRun starts the HTTP server and panics if an error occurs.
// This method is a convenience wrapper around the Run method,
// ensuring that if the server fails to start, the application will
// terminate immediately with a panic.
//
// Usage:
//
//	MustRun should only be used if you want the application to exit
//	in case the server fails to start. For more controlled error handling,
//	consider using the Run method directly.
func (a App) MustRun() {
	err := a.Run()
	if err != nil {
		panic(err)
	}
}

// Run starts http server.
//
// Returns:
//   - An error if the starting process encounters any issues; otherwise, nil.
func (a App) Run() error {
	const op = "app.Run"
	log := a.log.With(
		slog.String("op", op),
		slog.String("address", a.httpServer.Addr),
	)

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("http server is running")

	return nil
}

// Stop gracefully shuts down the server without interrupting any active connections.
// It waits for all the active requests to complete and then shuts down the server.
// This method is typically used for gracefully shutting down the server,
// for instance, when the application is receiving a termination signal.
//
// Parameters:
//   - ctx: A context.Context used to provide a deadline for the shutdown process.
//     The server will wait for active requests to finish until the context deadline.
//
// Returns:
//   - An error if the shutdown process encounters any issues; otherwise, nil.
func (a App) Stop(ctx context.Context) error {
	const op = "app.Stop"
	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	a.log.With(slog.String("op", op)).Info("stopping HTTP Server")
	return nil
}
