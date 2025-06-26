package main

import (
   "fmt"
   "net/http"
   "log/slog"
   "os" 

   "github.com/prometheus/client_golang/prometheus"
   "github.com/prometheus/client_golang/prometheus/promhttp"
)

var pingCounter = prometheus.NewCounter(
   prometheus.CounterOpts{
       Name: "ping_request_count",
       Help: "No of request handled by Ping handler",
   },
)

func ping(w http.ResponseWriter, req *http.Request) {
   pingCounter.Inc()
   fmt.Fprintf(w, "pong")
   
   logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
   logger.Info("ping...pong...")

}

func main() {
   prometheus.MustRegister(pingCounter)

   logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
   logger.Info("Server is starting... ")

   http.HandleFunc("/ping", ping)
   http.Handle("/metrics", promhttp.Handler())
   http.ListenAndServe(":8090", nil)
}
