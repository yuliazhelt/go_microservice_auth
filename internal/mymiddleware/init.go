package mymiddleware

import (
	"log"
	"io"
	"bytes"
	"fmt"
	"net/http"
	"go.uber.org/zap"
	"github.com/go-chi/chi/middleware"
	"go.opentelemetry.io/otel/trace"
)

func LoggerMiddleware(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			mw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Fatal(err)
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			defer func() {
				customLogger := l.With(
					zap.String("Body", string(body)),
					zap.String("Headers", fmt.Sprint(r.Header)),
					zap.Int("Status", mw.Status()),
				)
				customLogger.Debug("Done!")
			}()
			next.ServeHTTP(mw, r)
		}
		return http.HandlerFunc(f)
	}
}

func TracerMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			mw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			span := trace.SpanFromContext(r.Context())
			w.Header().Add("TraceID", span.SpanContext().TraceID().String())
			span.End()
			next.ServeHTTP(mw, r)
		}
		return http.HandlerFunc(f)
	}
}