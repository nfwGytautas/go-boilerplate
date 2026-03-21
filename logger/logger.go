// Package logger provides utility functions for setting up the go
// standard `slog` package with some "niceities" or for testing
package logger

import (
	"io"
	"log/slog"
	"os"
)

type config struct {
	writer io.Writer
	level  slog.Leveler
}

// Option various options for the logger
type Option func(*config)

// Initialize sets up the default slog logger with the given options
func Initialize(opts ...Option) {
	cfg := &config{
		writer: os.Stdout,
		level:  slog.LevelInfo,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	h := newHandler(cfg.writer, cfg.level)
	slog.SetDefault(slog.New(h))
}
