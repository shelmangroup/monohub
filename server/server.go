package server

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	api "github.com/shelmangroup/monohub/api"
	"github.com/shelmangroup/monohub/storage"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/examples/exporter"
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.opencensus.io/zpages"
	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Server struct {
	storage    *storage.Storage
	gitHandler *GitHandler
	httpServer *http.Server
}

var (
	command = kingpin.Command("server", "Server")
	dataDir = command.Flag("data-directory", "Data directory").Short('d').Required().String()

	grpcAddress   = command.Flag("gprc-address", "GRPC address").Default(":8823").String()
	listenAddress = command.Flag("listen-address", "HTTP address").Default(":8822").String()
	traceUrl      = command.Flag("trace-url", "Jaeger Trace URL.").Default("http://localhost:14268").String()
)

func FullCommand() string {
	return command.FullCommand()
}

func NewServer(storage *storage.Storage) *Server {

	server := &Server{
		storage: storage,
	}
	// server.gitHandler = NewGitHandler(storage)
	//
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", server.gitHandler.HandleRequest)
	// server.httpServer = &http.Server{
	// 	Addr:    *listenAddress,
	// 	Handler: mux,
	// }

	return server
}

func (s *Server) Serve() error {
	errCh := make(chan error)

	pe, err := prometheus.NewExporter(prometheus.Options{
		Namespace: "monohub",
	})
	if err != nil {
		return err
	}
	view.RegisterExporter(pe)

	exporter := &exporter.PrintExporter{}
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
		h := NewHttpServer(s.storage)
		if err := h.Run(); err != nil {
			errCh <- err
		}
	}()

	go func() {
		log.WithField("address", *grpcAddress).Info("Starting GRPC server")
		listen, err := net.Listen("tcp", *grpcAddress)
		if err != nil {
			errCh <- err
			return
		}
		gs := grpc.NewServer()
		api.RegisterMonohubServer(gs, s)
		err = gs.Serve(listen)
		if err != nil {
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
		return err
	}

	return nil
}

func RunServer() {
	storage := storage.NewStorage(*dataDir)
	srv := NewServer(storage)
	err := srv.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
