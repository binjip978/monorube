package collector

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Stats represent all collected statics at one point in time
type Stats struct {
	Zones Thermal
	Power
}

// Power related statistics
type Power struct {
	Volt float32
}

// Thermal collect basic thermal_zone_ info
type Thermal map[string]float32

// Config stores all information related to measurements statistics
type Config struct {
	ThermalPath string
	CpuInfoPath string
}

// PiInfo stores raspberry pi specific information: Serial, Model, Revision
type PiInfo struct {
	Hardware string
	Revision string
	Serial   string
	Model    string
}

// Collector gathers PiBox collector
type Collector struct {
	config   Config
	heatMap  map[string]string
	PiInfo   PiInfo
	Hostname string
}

// New return new collector
func New(config Config) *Collector {
	hostname, err := os.Hostname()
	if err != nil {
		log.Errorf("can't get hostname: %v", err)
	}

	return &Collector{
		config:   config,
		PiInfo:   piInfo(config.CpuInfoPath),
		Hostname: hostname,
	}
}

func piInfo(path string) PiInfo {
	var pi PiInfo
	cpuInfo, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("can't open %s : %v", path, err)
		return pi
	}

	var parse = func(keyword, line string) (string, error) {
		if strings.HasPrefix(line, keyword) {
			chunks := strings.Split(line, ": ")
			if len(chunks) == 2 {
				return chunks[1], nil
			}
		}

		return "", fmt.Errorf("can't find keyword")
	}

	scanner := bufio.NewScanner(bytes.NewReader(cpuInfo))
	for scanner.Scan() {
		line := scanner.Text()
		hardware, err := parse("Hardware", line)
		if err == nil {
			pi.Hardware = hardware
		}
		revision, err := parse("Revision", line)
		if err == nil {
			pi.Revision = revision
		}
		serial, err := parse("Serial", line)
		if err == nil {
			pi.Serial = serial
		}
		model, err := parse("Model", line)
		if err == nil {
			pi.Model = model
		}
	}

	return pi
}

func (c *Collector) temp() (Thermal, error) {
	thermal := make(map[string]float32)
	if c.heatMap == nil {
		files, err := ioutil.ReadDir(c.config.ThermalPath)
		if err != nil {
			return thermal, err
		}

		heapMap := make(map[string]string)
		for _, file := range files {
			if strings.HasPrefix(file.Name(), "thermal_zone") {
				typePath := filepath.Join(c.config.ThermalPath, file.Name(), "type")
				nf, err := ioutil.ReadFile(typePath)
				if err != nil {
					log.Errorf("can't open thermal type %s : %v", typePath, err)
					continue
				}
				heapMap[strings.TrimSpace(string(nf))] = filepath.Join(c.config.ThermalPath, file.Name(), "temp")
			}
		}

		if len(heapMap) == 0 {
			return thermal, fmt.Errorf("can't find any thermal_zone files")
		}

		c.heatMap = heapMap
	}

	for thrType, tempPath := range c.heatMap {
		temp, err := ioutil.ReadFile(tempPath)
		if err != nil {
			log.Errorf("can't open thermal temp file %s : %v", tempPath, err)
		}

		t, err := strconv.Atoi(strings.TrimSpace(string(temp)))
		if err != nil {
			log.Errorf("can't parse integer form thermal temp file %s : %v", tempPath, err)
		}
		thermal[thrType] = float32(t) / 1000
	}

	return thermal, nil
}

func (c *Collector) runVCG(arg string) ([]byte, error) {
	cmd := exec.Command("/usr/bin/vcgencmd", arg)
	return cmd.Output()
}

func (c *Collector) power() (Power, error) {
	var power Power
	b, err := c.runVCG("measure_volts")
	if err != nil {
		return power, err
	}

	v, err := parseVoltOutput(b)
	if err != nil {
		log.Errorf("can't parse volt ouput: %v", err)
		return power, err
	}
	power.Volt = v

	return power, nil
}

func parseVoltOutput(output []byte) (float32, error) {
	str := strings.TrimSpace(string(output))
	if !strings.Contains(str, "volt=") {
		return 0.0, fmt.Errorf("can't parse volt from vcgencmd")
	}

	chunks := strings.Split(str, "=")
	if len(chunks) != 2 {
		return 0.0, fmt.Errorf("can't parse volt from vcgencmd")
	}

	fp := chunks[1][:len(chunks[1])-1]
	v, err := strconv.ParseFloat(fp, 32)
	if err != nil {
		return 0.0, err
	}

	return float32(v), nil
}

// Measure return current collected system statistics
func (c *Collector) Measure() Stats {
	var stats Stats
	zones, err := c.temp()
	if err == nil {
		stats.Zones = zones
	} else {
		log.Debugf("can't get temp: %v", err)
	}

	power, err := c.power()
	if err == nil {
		stats.Power = power
	} else {
		log.Debugf("can't get power: %v", err)
	}

	return stats
}

// String return pretty printed string
func (s Stats) String() string {
	var zones []string
	for k := range s.Zones {
		zones = append(zones, k)
	}
	sort.Strings(zones)

	var buf bytes.Buffer
	needNewLine := false

	size := len(s.Zones)
	for i, zone := range zones {
		if i == size-1 {
			buf.WriteString(fmt.Sprintf("%16s  %2.2f \u2103", zone, s.Zones[zone]))
			needNewLine = true
			break
		}
		buf.WriteString(fmt.Sprintf("%16s  %2.2f \u2103\n", zone, s.Zones[zone]))
	}

	if s.Volt != 0.0 {
		if needNewLine {
			buf.WriteString("\n")
		}
		buf.WriteString(fmt.Sprintf("%16s  %2.2f V", "Volt", s.Volt))
	}

	return buf.String()
}
