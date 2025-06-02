package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/marpaia/graphite-golang"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Общее количество HTTP запросов",
		},
	)
)

func init() {
	prometheus.MustRegister(requestsTotal)
}

func main() {
	graphiteClient, err := graphite.NewGraphite("graphite", 2003)
	if err != nil {
		log.Printf("Graphite not available: %v", err)
	}
	counter := 0

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestsTotal.Inc()
		counter += 1
		_ = graphiteClient.SimpleSend("example_app.requests", fmt.Sprintf("%d", counter))
		fmt.Fprintf(w, "Hello, metrics!\n")
	})

	// http.Handle("/metrics", promhttp.Handler())

	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())

	go func() {
		err := http.ListenAndServe(":8088", metricsMux)
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
