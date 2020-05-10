package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/procfs/sysfs"
)

type config struct {
	period     time.Duration
	serverAddr string
}

func parseConfig() config {
	var period = flag.Duration("period", 10*time.Minute, "period for collecting temperature")
	var serverAddr = flag.String("srvAddr", ":5000", "addr to server prometheus metrics")
	flag.Parse()

	return config{period: *period, serverAddr: *serverAddr}
}

func tempLoop(period time.Duration) {
	fs, err := sysfs.NewDefaultFS()
	if err != nil {
		log.Panicf("can't init sysfs: %v", err)
	}
	for {
		zones, err := fs.ClassThermalZoneStats()
		if err != nil {
			log.Printf("error getting timezone: %v", err)
			time.Sleep(period)
			continue
		}

		for _, zone := range zones {
			raspTemp.Set(float64(zone.Temp) / 1000)
			time.Sleep(period)
		}
	}
}

var (
	raspTemp = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "current_temp",
		Help: "Current temperature in Celsius",
	})
)

func newServer(cfg config) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	srv := http.Server{
		Addr:    cfg.serverAddr,
		Handler: mux,
	}

	return &srv
}

func main() {
	cfg := parseConfig()
	srv := newServer(cfg)
	go tempLoop(cfg.period)
	log.Fatal(srv.ListenAndServe())
}
