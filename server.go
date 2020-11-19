package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"io"

	"github.com/jacob-elektronik/rabbit-amazon-forwarder/config"
	"github.com/jacob-elektronik/rabbit-amazon-forwarder/mapping"
	"github.com/jacob-elektronik/rabbit-amazon-forwarder/supervisor"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

const (
	LogLevel = "LOG_LEVEL"
)

func main() {
	createLogger()

	_, closer := initJaeger("rabbit-amazon-forwarder", false)
	defer closer.Close()

	consumerForwarderMapping, err := mapping.New().Load()
	if err != nil {
		log.WithField("error", err.Error()).Fatalf("Could not load consumer - forwarder pairs")
	}
	supervisor := supervisor.New(consumerForwarderMapping)
	if err := supervisor.Start(); err != nil {
		log.WithField("error", err.Error()).Fatal("Could not start supervisor")
	}

	basePath := os.Getenv("BASE_PATH")

	log.Info(fmt.Sprintf("Starting http server with path %s/restart", basePath))
	http.HandleFunc(fmt.Sprintf("%s/restart", basePath), supervisor.Restart)

	log.Info(fmt.Sprintf("Starting http server with path %s/reload", basePath))
	http.HandleFunc(fmt.Sprintf("%s/reload", basePath), supervisor.Reload)

	log.Info(fmt.Sprintf("Starting http server with path %s/health", basePath))
	http.HandleFunc(fmt.Sprintf("%s/health", basePath), supervisor.Check)

	log.Info("Starting http server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createLogger() {
	if strings.ToLower(os.Getenv(config.LogFormat)) == "text" {
		log.SetFormatter(&log.TextFormatter{})
	} else {
		log.SetFormatter(&log.JSONFormatter{})
	}

	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	if logLevel := os.Getenv(LogLevel); logLevel != "" {
		if level, err := log.ParseLevel(logLevel); err != nil {
			log.Fatal(err)
		} else {
			log.SetLevel(level)
		}
	}
}

func initJaeger(service string, debug bool) (opentracing.Tracer, io.Closer) {
	// see https://github.com/jaegertracing/jaeger-client-go#environment-variables
	cfg, err := jaegerConfig.FromEnv()
	if err != nil {
		log.WithField("error", err).Fatal("cannot init Jaeger")
	}

	tracer, closer, err := cfg.NewTracer(
		jaegerConfig.Logger(jaeger.StdLogger),
	)
	if err != nil {
		log.WithField("error", err).Fatal("cannot init Jaeger")
	}

	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}
