package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/opentracing-contrib/go-gorilla/gorilla"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

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

	log.Println("Listening on port 8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Println("Fail to listen and serve:", err)
	}

	log.Println("Stopped")
}

func opentracingMiddleware(next http.Handler) http.Handler {
	return gorilla.Middleware(opentracing.GlobalTracer(), next)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stderr, next)
}

func fib(n int) (int, error) {
	if n < 0 {
		return 0, errors.New("n cannot be less than 0")
	}

	if n < 2 {
		return n, nil
	}

	a, err := fib(n - 2)
	if err != nil {
		return 0, err
	}

	b, err := fib(n - 1)
	if err != nil {
		return 0, err
	}

	return a + b, nil
}
