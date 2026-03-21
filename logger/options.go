package logger

import (
	"io"
	"log/slog"
)

// WithWriter sets the output writer for the logger.
func WithWriter(w io.Writer) Option {
	return func(c *config) {
		c.writer = w
	}
}

// WithLevel sets the minimum log level.
func WithLevel(level slog.Leveler) Option {
	return func(c *config) {
		c.level = level
	}
}
