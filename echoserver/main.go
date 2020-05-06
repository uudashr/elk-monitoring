package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/opentracing/opentracing-go"
	"go.elastic.co/apm/module/apmot"
)

func main() {
	opentracing.SetGlobalTracer(apmot.New())

	r := mux.NewRouter()
	r.Use(opentracingMiddleware)
	r.Use(loggingMiddleware)

	// Echo
	r.HandleFunc("/echo", func(w http.ResponseWriter, req *http.Request) {
		contentType := req.Header.Get("Content-Type")
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println("Fail to read body:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", contentType)
		w.Write(b)
	}).Methods("POST")

	r.HandleFunc("/fib/{n}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		n, err := strconv.Atoi(vars["n"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		v, err := fib(n)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(strconv.Itoa(v)))
	}).Methods("GET")

	// Health check
	r.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		// empty will return 200 OK
	})

	// Metrics
	r.Handle("/metrics", promhttp.Handler())

	log.Println("Listening on port 8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Println("Fail to listen and serve:", err)
	}

	log.Println("Stopped")
}
