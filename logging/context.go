package logger

type ContextKey string

const (
	Request   ContextKey = "request"
	RequestId ContextKey = "request_id"
	Method    ContextKey = "method"
	Path      ContextKey = "path"
	SourceIP  ContextKey = "source_ip"
	Email     ContextKey = "email"
	UserId    ContextKey = "userId"
)
