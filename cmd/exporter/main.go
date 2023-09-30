package main

import (
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/kuskoman/logstash-exporter/collectors"
	"github.com/kuskoman/logstash-exporter/config"
	"github.com/kuskoman/logstash-exporter/server"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	warn := godotenv.Load()
	if warn != nil {
		log.Println(warn)
	}

	logger, err := config.SetupSlog()
	if err != nil {
		log.Fatalf("failed to setup slog: %s", err)
	}
	slog.SetDefault(logger)

	exporterConfig, err := config.GetConfig(config.ExporterConfigLocation)
	if err != nil {
		slog.Error("failed to get exporter config", "err", err)
		os.Exit(1)
	}

	host := exporterConfig.Server.Host
	port := strconv.Itoa(exporterConfig.Server.Port)

	slog.Debug("application starting... ")
	versionInfo := config.GetVersionInfo()
	slog.Info(versionInfo.String())

	collectorManager := collectors.NewCollectorManager(exporterConfig.GetLogstashUrls())
	prometheus.MustRegister(collectorManager)

	appServer := server.NewAppServer(host, port, exporterConfig)

	slog.Info("starting server on", "host", host, "port", port)
	if err := appServer.ListenAndServe(); err != nil {
		slog.Error("failed to listen and serve", "err", err)
		os.Exit(1)
	}
}
