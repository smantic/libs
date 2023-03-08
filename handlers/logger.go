package handlers

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type ZapMiddleware struct {
	logger *zap.Logger
	otel   otelAdapter
}

func NewZapMiddleware(l *zap.Logger, otel otelAdapter) ZapMiddleware {
	if otel == nil {
		otel = &emptyOtelAdapter{}
	}

	return ZapMiddleware{
		logger: l,
		otel:   otel,
	}
}

type ResponseWriter struct {
	http.ResponseWriter
	StatusCode   int
	BytesWritten int
}

var ctxkey struct{}

func WithLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxkey, l)
}

func (z *ZapMiddleware) LogRequest(h http.Handler) http.Handler {

	f := func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		l := z.logger.With(
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Time("start", start),
			zap.String("trace-id", z.otel.GetTraceID(r)),
			zap.String("request-id", z.otel.GetRequestID(r)),
		)

		r = r.WithContext(WithLogger(r.Context(), l))
		ww := ResponseWriter{ResponseWriter: w}

		h.ServeHTTP(ww, r)
		finish := time.Now()

		dur := time.Since(start)
		l.Info(
			"request-finished",
			zap.Duration("duration", dur),
			zap.Time("finish", finish),
			zap.Int("status", ww.StatusCode),
			zap.Int("bytes-written", ww.BytesWritten),
			zap.String("user-agent", r.UserAgent()),
			zap.Strings("transfer-encodings", r.TransferEncoding),
			zap.String("proto", r.Proto),
			zap.String("referer", r.Referer()),
			zap.Int("content-length", int(r.ContentLength)),
			zap.String("remote-addr", r.RemoteAddr),
		)
	}

	return http.HandlerFunc(f)
}
