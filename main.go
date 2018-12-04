package main

import (
	// "os"
	"strings"

	"github.com/shelmangroup/monohub/server"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	logJSON  = kingpin.Flag("log-json", "Use structured logging in JSON format").Default("false").Bool()
	logLevel = kingpin.Flag("log-level", "The level of logging").Default("info").Enum("debug", "info", "warn", "error", "panic", "fatal")
)

func main() {
	kingpin.HelpFlag.Short('h')
	kingpin.CommandLine.DefaultEnvars()
	kingpin.Parse()

	switch strings.ToLower(*logLevel) {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	if *logJSON {
		log.SetFormatter(&log.JSONFormatter{})
	}

	kingpin.Parse()

	log.Info("Here we go!")

	srv := server.NewGitServer()
	err := srv.Serve()
	if err != nil {
		log.Fatal(err)
	}

}
