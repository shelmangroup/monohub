package main

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/src-d/go-git.v4"
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

	url := "https://github.com/shelmangroup/oidc-agent.git"
	log.Infof("git clone %s", url)

	s, err := NewStorage()
	if err != nil {
		log.Fatal(err)
	}

	_, err = git.Clone(s, nil, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})

	if err != nil {
		log.Fatal(err)
	}

}
