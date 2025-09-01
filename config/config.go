package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// Holds the configuration for InfluxDB.
type InfluxConfig struct {
	URL    string
	Token  string
	Org    string
	Bucket string
}

// Holds the data collection intervals.
type IntervalConfig struct {
	CPU time.Duration
	GPU time.Duration
}

// Reads and returns InfluxDB configuration from environment variables.
func LoadInfluxConfig() InfluxConfig {
	return InfluxConfig{
		URL:    os.Getenv("INFLUXDB_URL"),
		Token:  os.Getenv("INFLUXDB_TOKEN"),
		Org:    os.Getenv("INFLUXDB_ORG"),
		Bucket: os.Getenv("INFLUXDB_BUCKET"),
	}
}

// Reads and returns the collection intervals from environment variables.
func LoadIntervalConfig() IntervalConfig {
	cpuIntvStr := os.Getenv("CPU_INTERVAL_SECONDS")
	if cpuIntvStr == "" {
		cpuIntvStr = "10"
	}
	cpuIntv, err := strconv.Atoi(cpuIntvStr)
	if err != nil {
		log.Fatalf("CPU data collection interval is invalid: %v", err)
	}

	gpuIntvStr := os.Getenv("GPU_INTERVAL_SECONDS")
	if gpuIntvStr == "" {
		gpuIntvStr = "10"
	}
	gpuIntv, err := strconv.Atoi(gpuIntvStr)
	if err != nil {
		log.Fatalf("GPU data collection interval is invalid: %v", err)
	}

	return IntervalConfig{
		CPU: time.Duration(cpuIntv) * time.Second,
		GPU: time.Duration(gpuIntv) * time.Second,
	}
}
