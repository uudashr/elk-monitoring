package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/opentracing-contrib/go-gorilla/gorilla"
	"github.com/opentracing/opentracing-go"
)

func opentracingMiddleware(next http.Handler) http.Handler {
	return gorilla.Middleware(opentracing.GlobalTracer(), next)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stderr, next)
}
