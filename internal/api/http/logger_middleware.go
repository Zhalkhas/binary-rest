package http

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"
)

type SlogFormatter struct{ slog.Logger }

func (l SlogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &slogEntry{
		args: make([]slog.Attr, 0, 8),
	}

	entry.args = append(entry.args, slog.String("from", r.RemoteAddr))

	reqID := middleware.GetReqID(r.Context())
	if reqID != "" {
		entry.args = append(entry.args, slog.String("req_id", reqID))
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	entry.msg = fmt.Sprintf("%s://%s%s %s", scheme, r.Host, r.RequestURI, r.Proto)

	return entry
}

type slogEntry struct {
	msg  string
	args []slog.Attr
}

func (l *slogEntry) Write(status, bytes int, _ http.Header, elapsed time.Duration, _ interface{}) {
	l.args = append(l.args, slog.Int("status", status))
	l.args = append(l.args, slog.Int("bytes", bytes))
	l.args = append(l.args, slog.Duration("elapsed", elapsed))
	slog.LogAttrs(context.Background(), slog.LevelInfo, l.msg, l.args...)
}

func (l *slogEntry) Panic(v interface{}, stack []byte) {
	slog.Log(context.Background(), slog.LevelError, fmt.Sprint(v), "stack", string(stack))
}
