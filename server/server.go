package server

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	// "go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/examples/exporter"
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.opencensus.io/zpages"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	httpListen = kingpin.Flag("http-listen", "HTTP API Listen address.").Default(":8822").String()
	traceUrl   = kingpin.Flag("trace-url", "Jaeger Trace URL.").Default("http://localhost:14268").String()
)

func NewGitServer() {
	errCh := make(chan error)
	pe, err := prometheus.NewExporter(prometheus.Options{
		Namespace: "monohub",
	})
	if err != nil {
		log.Fatalf("Failed to create the Prometheus stats exporter: %v", err)
	}
	view.RegisterExporter(pe)

	exporter := &exporter.PrintExporter{}
	// exporter, err := jaeger.NewExporter(jaeger.Options{
	// 	Endpoint:    *traceUrl,
	// 	ServiceName: "gannet",
	// })
	// defer exporter.Flush()
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	view.SetReportingPeriod(1 * time.Second)

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", pe)
		zpages.Handle(mux, "/")
		log.WithField("address", ":8888").Info("Starting metrics server")
		if err := http.ListenAndServe(":8888", mux); err != nil {
			errCh <- err
		}
	}()

	go func() {
		h := NewHttpServer()
		if err := h.Run(); err != nil {
			errCh <- err
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigCh:
		log.Warn("Received SIGTERM, exiting gracefully...")
	case err := <-errCh:
		log.WithError(err).Error("Got an error from errCh, exiting gracefully")
	}
}
