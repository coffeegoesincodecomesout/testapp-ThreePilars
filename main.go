package main

import (
   "fmt"
   "net/http"
   "log/slog"
   "os"
   "time"
   "context" 

   "github.com/prometheus/client_golang/prometheus"
   "github.com/prometheus/client_golang/prometheus/promhttp"

   "go.opentelemetry.io/otel"
   "go.opentelemetry.io/otel/codes"
   "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
   "go.opentelemetry.io/otel/sdk/trace"
)

var pingCounter = prometheus.NewCounter(
   prometheus.CounterOpts{
       Name: "ping_request_count",
       Help: "No of request handled by Ping handler",
   },
)

func main() {
   ctx := context.Background()
   exp, err := otlptracegrpc.New(
           ctx,
	   otlptracegrpc.WithInsecure(),
   )
   if err != nil {
           panic(err)
   }

   tracerProvider := trace.NewTracerProvider(trace.WithBatcher(exp))
   defer func() {
           if err := tracerProvider.Shutdown(ctx); err != nil {
	           panic(err)
	   }
   }()
   otel.SetTracerProvider(tracerProvider)   


   prometheus.MustRegister(pingCounter)

   http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
           trace := otel.Tracer("http-server")
		_, span := trace.Start(r.Context(), "handleRequest")
		defer span.End()

		time.Sleep(1 * time.Second)
		span.SetStatus(codes.Ok, "Status 200")
                
                pingCounter.Inc()
                fmt.Fprintf(w, "pong")
  
                logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
                logger.Info("ping...pong...")

   }) 

   logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))          
   logger.Info("Server is starting... ")              

   http.Handle("/metrics", promhttp.Handler())
   http.ListenAndServe(":8090", nil)
}
