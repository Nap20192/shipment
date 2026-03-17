package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Nap20192/shipment/internal/pkg/kernel"
)

type LogSubscriber struct {
	log *slog.Logger
}

func NewLogSubscriber() *LogSubscriber {
	return &LogSubscriber{
		log: slog.Default(),
	}
}

func (l *LogSubscriber) Handle(ctx context.Context, event kernel.DomainEvent) error {
	l.log.InfoContext(ctx, fmt.Sprintf("Event received: %s, payload: %+v", event.Name(), event.Payload()))
	return nil
}


type HandlerMiddleware struct {
	log slog.Handler
}

func NewHandlerMiddleware(log slog.Handler) *HandlerMiddleware {
	return &HandlerMiddleware{
		log: log,
	}
}

func (h *HandlerMiddleware) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.log.Enabled(ctx, rec)
}

func (h *HandlerMiddleware) Handle(ctx context.Context, r slog.Record) error {
	if c, ok := ctx.Value(ctxLoggerKey).(logCtx); ok {
		if c.requestID != "" {
			r.Add(slog.String("request_id", c.requestID))
		}
		if c.action != "" {
			r.Add(slog.String("action", c.action))
		}
		if c.service != "" {
			r.Add(slog.String("service", c.service))
		}
		if c.hostname != "" {
			r.Add(slog.String("hostname", c.hostname))
		}
		if c.userID != "" {
			r.Add(slog.String("user_id", c.userID))
		}
	}

	return h.log.Handle(ctx, r)
}

func (h *HandlerMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.log.WithAttrs(attrs)
}

func (h *HandlerMiddleware) WithGroup(name string) slog.Handler {
	return h.log.WithGroup(name)
}

type contextKey struct{}

var ctxLoggerKey = contextKey{}

type logCtx struct {
	requestID string
	action    string
	service   string
	hostname  string
	userID    string
}

func WithUserID(ctx context.Context, userID string) context.Context {
	c := getOrCreateLogCtx(ctx)
	c.userID = userID
	return context.WithValue(ctx, ctxLoggerKey, c)
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	c := getOrCreateLogCtx(ctx)
	c.requestID = requestID
	return context.WithValue(ctx, ctxLoggerKey, c)
}

func WithAction(ctx context.Context, action string) context.Context {
	c := getOrCreateLogCtx(ctx)
	c.action = action
	return context.WithValue(ctx, ctxLoggerKey, c)
}

func WithService(ctx context.Context, service string) context.Context {
	c := getOrCreateLogCtx(ctx)
	c.service = service
	return context.WithValue(ctx, ctxLoggerKey, c)
}

func WithHostname(ctx context.Context, hostname string) context.Context {
	c := getOrCreateLogCtx(ctx)
	c.hostname = hostname
	return context.WithValue(ctx, ctxLoggerKey, c)
}

func getOrCreateLogCtx(ctx context.Context) logCtx {
	if c, ok := ctx.Value(ctxLoggerKey).(logCtx); ok {
		return c
	}
	return logCtx{}
}

func WithSessionId(ctx context.Context, sessionId string) context.Context {
	return WithRequestID(ctx, sessionId)
}

func WithName(ctx context.Context, name string) context.Context {
	return WithAction(ctx, name)
}
