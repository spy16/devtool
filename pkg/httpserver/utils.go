package httpserver

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/spy16/devtool/pkg/log"
)

// Serve runs an HTTP server with given handler. When the provided context is
// cancelled, server shuts down gracefully.
func Serve(ctx context.Context, addr string, h http.Handler) error {
	httpServer := &http.Server{
		Addr:    addr,
		Handler: h,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("failed to shutdown server gracefully: %v", err)
		}
	}()

	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
