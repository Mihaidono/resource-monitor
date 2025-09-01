package workers

import (
	"time"

	"github.com/Mihaidono/resource-monitor/metrics"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

// Background routine that periodically collects and writes CPU metrics.
func CpuWorker(writeAPI api.WriteAPI, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		point := metrics.CollectCPUMetricsPoint()
		writeAPI.WritePoint(point)
	}
}

// Background routine that periodically collects and writes GPU metrics.
func GpuWorker(writeAPI api.WriteAPI, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		points := metrics.CollectGPUMetricsPoints()
		for _, point := range points {
			writeAPI.WritePoint(point)
		}
	}
}
