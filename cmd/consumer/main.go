package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func main() {
	http.HandleFunc("healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		if err := http.ListenAndServe(":81", nil); err != nil {
			panic(err)
		}
	}()
}
