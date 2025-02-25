package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/kuskoman/logstash-exporter/collectors"
	"github.com/kuskoman/logstash-exporter/config"
	"github.com/kuskoman/logstash-exporter/server"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	version := flag.Bool("version", false, "prints the version and exits")

	flag.Parse()
	if *version {
		fmt.Printf("%s\n", config.SemanticVersion)
		return
	}

	warn := godotenv.Load()
	if warn != nil {
		log.Println(warn)
	}

	logger, err := config.SetupSlog()
	if err != nil {
		log.Fatalf("failed to setup slog: %s", err)
	}
	slog.SetDefault(logger)

	port, host := config.Port, config.Host
	logstashUrl := config.LogstashUrl

	slog.Debug("application starting... ")
	versionInfo := config.GetVersionInfo()
	slog.Info(versionInfo.String())

	collectorManager := collectors.NewCollectorManager(logstashUrl)
	appServer := server.NewAppServer(host, port)
	prometheus.MustRegister(collectorManager)

	slog.Info("starting server on", "host", host, "port", port)
	if err := appServer.ListenAndServe(); err != nil {
		slog.Error("failed to listen and serve", "err", err)
		os.Exit(1)
	}
}
