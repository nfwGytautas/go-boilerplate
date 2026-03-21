package graceful

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

// GrpcServer starts a gRPC server and stops it gracefully when the context is cancelled
func GrpcServer(ctx context.Context, addr string, srv *grpc.Server) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on: %s (%w)", addr, err)
	}
	defer listener.Close()

	done := make(chan struct{})

	go func() {
		<-ctx.Done()
		srv.GracefulStop()
		close(done)
	}()

	if err := srv.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC server %w", err)
	}

	<-done

	return nil
}
