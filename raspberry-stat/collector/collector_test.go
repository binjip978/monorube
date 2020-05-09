package collector

import (
	"github.com/sirupsen/logrus"
	"strings"
	"testing"
)

func TestThermal(t *testing.T) {
	logrus.SetLevel(logrus.PanicLevel)
	clt := New(Config{ThermalPath: "./test_misc/thermal"})
	thermal, err := clt.temp()
	if err != nil {
		t.Fatalf("collector should return without error, but error was: %v", err)
	}

	if len(thermal) != 1 || thermal["cpu-thermal"] != 55.017 {
		t.Fatalf("first measurement: expected %f and size 1, got %v and size %d", 55.017,
			thermal["cpu-thermal"], len(thermal))
	}

	thermal2, err := clt.temp()
	if err != nil {
		t.Fatalf("second measurement should success as well, got error %v", err)
	}
	if len(thermal2) != 1 || thermal2["cpu-thermal"] != 55.017 {
		t.Fatalf("second measurement: expected %f and size 1, got %v and size %d", 55.017,
			thermal2["cpu-thermal"], len(thermal2))
	}
}

func TestPiInfo(t *testing.T) {
	pi := piInfo("./test_misc/cpuinfo")
	if pi.Model != "Raspberry Pi 4 Model B Rev 1.1" {
		t.Errorf("model name was parsed incorrectly: %s", pi.Model)
	}
	if pi.Serial != "10000000bfa8b96a" {
		t.Errorf("serial was parsed incorrectly: %s", pi.Serial)
	}
	if pi.Revision != "b03111" {
		t.Errorf("revison was parsed incorrectly: %s", pi.Revision)
	}
	if pi.Hardware != "BCM2835" {
		t.Errorf("hardware was parsed incorrectly: %s", pi.Hardware)
	}
}

func TestParseVoltOutput(t *testing.T) {
	b := []byte("volt=0.8100V\n")
	res, err := parseVoltOutput(b)
	if err != nil {
		t.Errorf("should parse volt correctly, got error: %v", err)
	}
	if res != 0.81 {
		t.Errorf("should be 0.81 was: %f", res)
	}
}

func TestStatsString(t *testing.T) {
	logrus.SetLevel(logrus.PanicLevel)
	clt := New(Config{ThermalPath: "./test_misc/thermal"})
	stats := clt.Measure()
	if !strings.Contains(stats.String(), "cpu-thermal") {
		t.Error("should contain cpu-thermal")
	}
}
