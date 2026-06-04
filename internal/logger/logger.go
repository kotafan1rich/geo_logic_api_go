package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)

	With(attrs ...any) Logger
	WithLatency(latency time.Duration) Logger
	WithRequest(status int, method, path, query, ip, userAgent string, latency time.Duration) Logger
	LogError(ctx context.Context, err error, message string)

	GinLoggerMiddleware() gin.HandlerFunc
}

type logger struct {
	*slog.Logger
}

type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

type LogFormat string

const (
	FormatJson LogFormat = "json"
	FormatText LogFormat = "text"
)

func New(level LogLevel, format LogFormat, addSource bool, output io.Writer) Logger {
	var slogLevel slog.Level
	switch level {
	case LevelDebug:
		slogLevel = slog.LevelDebug
	case LevelInfo:
		slogLevel = slog.LevelInfo
	case LevelWarn:
		slogLevel = slog.LevelWarn
	case LevelError:
		slogLevel = slog.LevelError
	default:
		panic("bad log level")
	}

	writer := output
	if writer == nil {
		writer = os.Stdout
	}
	opts := &slog.HandlerOptions{
		Level:     slogLevel,
		AddSource: addSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format(time.RFC3339))
				}
			}
			return a
		},
	}

	var handler slog.Handler
	switch format {
	case FormatJson:
		handler = slog.NewJSONHandler(writer, opts)
	case FormatText:
		handler = slog.NewTextHandler(writer, opts)
	default:
		panic("bad log format")
	}

	return &logger{
		Logger: slog.New(handler),
	}
}

func Default() Logger {
	return New(LevelInfo, FormatText, false, os.Stdout)
}

func (l *logger) Info(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}

func (l *logger) Warn(msg string, args ...any) {
	l.Logger.Warn(msg, args...)
}
func (l *logger) Error(msg string, args ...any) {
	l.Logger.Error(msg, args...)
}

func (l *logger) With(attrs ...any) Logger {
	return &logger{Logger: l.Logger.With(attrs...)}
}

func (l *logger) WithLatency(latency time.Duration) Logger {
	return &logger{
		Logger: l.Logger.With(slog.Duration("latency", latency)),
	}
}

func (l *logger) WithRequest(status int, method, path, query, ip, userAgent string, latency time.Duration) Logger {
	return &logger{
		Logger: l.Logger.With(
			slog.Int("status", status),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("query", query),
			slog.String("ip", ip),
			slog.String("user-agent", userAgent),
			slog.Duration("latency", latency),
		)}
}

func (l *logger) LogError(ctx context.Context, err error, message string) {
	if err != nil {
		pc := make([]uintptr, 10)
		n := runtime.Callers(2, pc)
		frames := runtime.CallersFrames(pc[:n])

		var file string
		var line int

		if frame, more := frames.Next(); more {
			file = frame.File
			line = frame.Line
		}

		l.ErrorContext(ctx, message,
			slog.String("error", err.Error()),
			slog.String("file", file),
			slog.Int("line", line),
		)
	}
}
