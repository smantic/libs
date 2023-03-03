package handlers

import "net/http"

type otelAdapter interface {
	GetRequestID(r *http.Request) string
	GetTraceID(r *http.Request) string
}

type emptyOtelAdapter struct{}

func (e *emptyOtelAdapter) GetRequestID(r *http.Request) string {
	return ""
}

func (e *emptyOtelAdapter) GetTraceID(r *http.Request) string {
	return ""
}
