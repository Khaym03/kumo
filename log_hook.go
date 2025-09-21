package main

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// LogrusHook is a logrus hook that sends logs to the Wails frontend.
type LogrusHook struct {
	ctx context.Context
}

// NewLogrusHook creates a new instance of the hook.
func NewLogrusHook(ctx context.Context) *LogrusHook {
	return &LogrusHook{ctx: ctx}
}

// Fire is the main method that handles the log event.
func (h *LogrusHook) Fire(entry *logrus.Entry) error {
	// Emit a custom event with the log message and level.
	// You can send any data you need.
	runtime.EventsEmit(h.ctx, "log_message", map[string]string{
		"level":   entry.Level.String(),
		"message": entry.Message,
		"time":    entry.Time.Format(time.RFC3339),
	})
	return nil
}

// Levels returns the log levels this hook will be triggered for.
func (h *LogrusHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
