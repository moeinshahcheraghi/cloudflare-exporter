package main

import (
	"log"
	"net/http"
	"time"

	"cloudflare-exporter/internal/collector"
	"cloudflare-exporter/internal/config"
	"cloudflare-exporter/internal/metrics"
	"cloudflare-exporter/pkg/cloudflare"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Load configuration
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf(" Configuration error: %v", err)
	}

	log.Println(" Cloudflare Prometheus Exporter")
	log.Printf(" Zone: %s | Port: %s | Interval: %v", cfg.ZoneID, cfg.Port, cfg.ScrapeInterval)

	cfClient := cloudflare.NewClient(cfg.APIToken)

	metricsRegistry := metrics.NewMetrics()
	metricsRegistry.Register()

	col := collector.NewCollector(cfClient, metricsRegistry, cfg.ZoneID)

	startPeriodicCollection(col, cfg.ScrapeInterval)

	setupHTTPServer(cfg.Port)
}

func startPeriodicCollection(col *collector.Collector, interval time.Duration) {
	log.Println(" Performing initial metrics collection...")
	if err := col.CollectAll(); err != nil {
		log.Printf("  Initial collection had errors: %v", err)
	}

	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			if err := col.CollectAll(); err != nil {
				log.Printf("  Collection had errors: %v", err)
			}
		}
	}()
}

func setupHTTPServer(port string) {
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<html>
			<head><title>Cloudflare Exporter</title></head>
			<body>
				<h1>Cloudflare Prometheus Exporter</h1>
				<p><a href="/metrics">Metrics</a></p>
				<p><a href="/health">Health Check</a></p>
			</body>
			</html>
		`))
	})

	log.Printf(" Server listening on http://localhost:%s", port)
	log.Printf(" Metrics available at http://localhost:%s/metrics", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf(" Failed to start server: %v", err)
	}
}