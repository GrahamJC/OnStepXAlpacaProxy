package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
)

const (
	reset = "\033[0m"

	black        = 30
	red          = 31
	green        = 32
	yellow       = 33
	blue         = 34
	magenta      = 35
	cyan         = 36
	lightGray    = 37
	darkGray     = 90
	lightRed     = 91
	lightGreen   = 92
	lightYellow  = 93
	lightBlue    = 94
	lightMagenta = 95
	lightCyan    = 96
	white        = 97
)

func colorize(s string, colorCode int) string {
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), s, reset)
}

type Handler struct {
	h slog.Handler
	b *bytes.Buffer
	m *sync.Mutex
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{h: h.h.WithAttrs(attrs), b: h.b, m: h.m}
}

func (h *Handler) WithGroup(groupName string) slog.Handler {
	return &Handler{h: h.h.WithGroup(groupName), b: h.b, m: h.m}
}

const (
	timeFormat = "[15:04:05.000]"
)

func (h *Handler) ComputeAttrs(ctx context.Context, r slog.Record) (map[string]any, error) {

	// Call inner handler to write formatted output to buffer
	h.m.Lock()
	defer func() {
		h.b.Reset()
		h.m.Unlock()
	}()
	if err := h.h.Handle(ctx, r); err != nil {
		return nil, fmt.Errorf("error when calling inner handler's Handle(): %w", err)
	}

	// Unmarshal JSON into map
	var attrs map[string]any
	err := json.Unmarshal(h.b.Bytes(), &attrs)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling inner handler's Handle result: %w", err)
	}
	return attrs, nil
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {

	// Get level and colorize it
	level := r.Level.String() + ":"
	switch r.Level {
	case slog.LevelDebug:
		level = colorize(level, darkGray)
	case slog.LevelInfo:
		level = colorize(level, cyan)
	case slog.LevelWarn:
		level = colorize(level, yellow)
	case slog.LevelError:
		level = colorize(level, red)
	}

	// Format name/value pairs using JSON handler
	attrs, err := h.ComputeAttrs(ctx, r)
	if err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(attrs, "", "  ")
	if err != nil {
		return fmt.Errorf("error unmarshaling attrs: %w", err)
	}

	// Log message to console
	fmt.Println(
		colorize(r.Time.Format(timeFormat), lightGray),
		level,
		colorize(r.Message, white),
		colorize(string(bytes), darkGray),
	)

	// no error
	return nil
}

func suppressDefaults(next func([]string, slog.Attr) slog.Attr) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey || a.Key == slog.LevelKey || a.Key == slog.MessageKey {
			return slog.Attr{}
		}
		if next == nil {
			return a
		}
		return next(groups, a)
	}
}

func NewHandler(opts *slog.HandlerOptions) *Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	b := &bytes.Buffer{}
	return &Handler{
		h: slog.NewJSONHandler(b, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: suppressDefaults(opts.ReplaceAttr),
		}),
		b: b,
		m: &sync.Mutex{},
	}
}
