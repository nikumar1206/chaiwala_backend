package logger

import (
	"context"
	"log/slog"
)

type FiberHandler struct {
	slog.Handler
}

func (l FiberHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx.Value(RequestId) == nil {
		return l.Handler.Handle(ctx, r)
	}

	requestId := ctx.Value(RequestId).(string)
	sourceIp := ctx.Value(SourceIP).(string)
	path := ctx.Value(Path).(string)
	method := ctx.Value(Method).(string)

	requestGroup := slog.Group(
		string(Request),
		slog.String(string(RequestId), requestId),
		slog.String(string(SourceIP), sourceIp),
		slog.String(string(Method), method),
		slog.String(string(Path), path),
	)

	r.AddAttrs(requestGroup)

	return l.Handler.Handle(ctx, r)
}
