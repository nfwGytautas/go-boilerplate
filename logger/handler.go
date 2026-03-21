package logger

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"sync"
)

type handler struct {
	w     io.Writer
	mu    *sync.Mutex
	level slog.Leveler
	attrs []slog.Attr
	group string
}

func newHandler(w io.Writer, level slog.Leveler) *handler {
	return &handler{
		w:     w,
		mu:    &sync.Mutex{},
		level: level,
	}
}

func (h *handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *handler) Handle(_ context.Context, r slog.Record) error {
	fields := make(map[string]any, r.NumAttrs()+3)

	fields["time"] = r.Time
	fields["level"] = r.Level.String()
	fields["msg"] = r.Message

	target := fields
	if h.group != "" {
		g := make(map[string]any)
		fields[h.group] = g
		target = g
	}

	for _, a := range h.attrs {
		target[a.Key] = a.Value.Any()
	}

	r.Attrs(func(a slog.Attr) bool {
		target[a.Key] = a.Value.Any()
		return true
	})

	b, err := json.Marshal(fields)
	if err != nil {
		return err
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	_, err = h.w.Write(append(b, '\n'))
	return err
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs), len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	newAttrs = append(newAttrs, attrs...)

	return &handler{
		w:     h.w,
		mu:    h.mu,
		level: h.level,
		attrs: newAttrs,
		group: h.group,
	}
}

func (h *handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	newGroup := name
	if h.group != "" {
		newGroup = h.group + "." + name
	}

	return &handler{
		w:     h.w,
		mu:    h.mu,
		level: h.level,
		attrs: h.attrs,
		group: newGroup,
	}
}
