package main

import (
	"flag"
	"fmt"
	"github.com/binjip978/rapsberry-stat/collector"

	log "github.com/sirupsen/logrus"
)

type piStatConfig struct {
	cltConfig collector.Config
	logLevel  string
}

func parseConfig() piStatConfig {
	var thermalPath = flag.String("thermal", "/sys/class/thermal", "path to thermal data folder")
	var cpuInfoPath = flag.String("cpuinfo", "/proc/cpuinfo", "path to cpuinfo")
	var logLevel = flag.String("logLevel", "ERROR", "log level: (ERROR|DEBUG)")
	flag.Parse()

	cltConfig := collector.Config{
		ThermalPath: *thermalPath,
		CpuInfoPath: *cpuInfoPath,
	}

	return piStatConfig{
		cltConfig: cltConfig,
		logLevel:  *logLevel,
	}
}

func main() {
	cfg := parseConfig()
	if cfg.logLevel == "DEBUG" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.ErrorLevel)
	}

	clt := collector.New(cfg.cltConfig)
	fmt.Println(clt.Measure())
}
