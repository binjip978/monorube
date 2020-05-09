package main

import (
	"testing"
)

func TestParseConfig(t *testing.T) {
	cfg := parseConfig()

	if cfg.logLevel != "ERROR" {
		t.Errorf("default log error should be ERROR, was %s", cfg.logLevel)
	}

	if cfg.cltConfig.ThermalPath != "/sys/class/thermal" {
		t.Errorf("default thermal path should be /sys/class/thermal was %s", cfg.cltConfig.ThermalPath)
	}

	if cfg.cltConfig.CpuInfoPath != "/proc/cpuinfo" {
		t.Errorf("default cpuinof path should be /proc/cpuinfo was %s", cfg.cltConfig.CpuInfoPath)
	}
}
