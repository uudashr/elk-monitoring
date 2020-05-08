package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.elastic.co/apm/module/apmot"
)

func main() {
	opentracing.SetGlobalTracer(apmot.New())

	r := mux.NewRouter()
	r.Use(opentracingMiddleware)
	r.Use(loggingMiddleware)

	r.HandleFunc("/random/ints/:n", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramN := vars["n"]
		n, err := strconv.Atoi(paramN)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		i := rand.Intn(n)

		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(i); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	// Health check
	r.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		// empty will return 200 OK
	})

	// Metrics
	r.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Handler: r,
		Addr:    ":8080",
	}

	if err := server.ListenAndServe(); err != nil {
		log.Println("Fail to listen and server:", err)
	}
}
