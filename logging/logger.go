package logger

import (
	"context"
	"log/slog"
)

type CustomHandler struct {
	slog.Handler
}

func (l CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx.Value(RequestId) == nil {
		return l.Handler.Handle(ctx, r)
	}

	requestId := ctx.Value(RequestId).(string)
	sourceIp := ctx.Value(SourceIP).(string)
	path := ctx.Value(Path).(string)
	method := ctx.Value(Method).(string)
	email := ctx.Value(Email)
	userId := ctx.Value(UserId)

	requestGroup := slog.Group(
		string(Request),
		slog.String(string(RequestId), requestId),
		slog.String(string(SourceIP), sourceIp),
		slog.String(string(Method), method),
		slog.String(string(Path), path),
		slog.Any(string(Email), email),
		slog.Any(string(UserId), userId),
	)

	r.AddAttrs(requestGroup)

	return l.Handler.Handle(ctx, r)
}
