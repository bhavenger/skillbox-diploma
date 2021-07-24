package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "skillbox_http_requests_total",
		Help: "The total number of processed requests",
	})
)

func main() {
	logger := log.New(os.Stdout, "[skillbox]: ", log.LstdFlags)
	logger.Println("Server is starting...")

	router := http.NewServeMux()

	router.Handle("/metrics", promhttp.Handler())
	router.Handle("/health", healthCheck())
	router.Handle("/", mainHandler())

	server := &http.Server{
		Addr:     "0.0.0.0:8080",
		Handler:  logRequest(logger)(router),
		ErrorLog: logger,
	}

	logger.Println("Waiting for requests at http://0.0.0.0:8080/")
	server.ListenAndServe()
}

func mainHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		opsProcessed.Inc()
		fmt.Fprintln(w, "Hello, skillbox!")
		fmt.Fprintln(w, "Metod:"+r.Method)
		fmt.Fprintln(w, "Uri:"+r.RequestURI)
		fmt.Fprintln(w, "Host:"+r.Host)
		fmt.Fprintln(w, "Headers are:")
		fmt.Fprintln(w, r.Header)
		fmt.Fprintln(w, "Body are:")
		fmt.Fprintln(w, r.Body)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	})
}

func healthCheck() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
}

func logRequest(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				logger.Println(r.Method, r.RequestURI, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}
