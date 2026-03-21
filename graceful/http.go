// Package graceful provides utility functions for starting various
// server handlers that gracefully exit when a context is cancelled
package graceful

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// HttpServer start a graceful http server that stops when the context is cancelled
func HttpServer(ctx context.Context, addr string, handler http.Handler) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	done := make(chan struct{})

	go func() {
		<-ctx.Done()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer shutdownCancel()
		srv.Shutdown(shutdownCtx)
		close(done)
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to listen and serve: %w", err)
	}

	<-done

	return nil
}
