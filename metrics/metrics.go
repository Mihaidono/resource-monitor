package metrics

import (
	"log"
	"strconv"
	"time"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

const MEGABYTE = 1024 * 1024

// Collects CPU and memory metrics and returns an InfluxDB point.
func CollectCPUMetricsPoint() *write.Point {
	metrics := make(map[string]any)

	percentages, err := cpu.Percent(0, false)
	if err != nil {
		log.Printf("Error getting CPU usage: %v", err)
	} else {
		metrics["cpu_usage_percentage"] = percentages[0]
	}

	v, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Error getting memory usage: %v", err)
	} else {
		metrics["memory_total_megabytes"] = v.Total / MEGABYTE
		metrics["memory_used_megabytes"] = v.Used / MEGABYTE
		metrics["memory_usage_percent"] = v.UsedPercent
	}

	point := influxdb2.NewPoint("cpu", nil, metrics, time.Now())
	return point
}

// Collects metrics for all detected GPUs and returns a slice of InfluxDB points.
func CollectGPUMetricsPoints() []*write.Point {
	var points []*write.Point
	count, result := nvml.DeviceGetCount()
	if result != nvml.SUCCESS {
		log.Printf("Failed to get GPU count from NVML: %v", result)
		return nil
	}

	if count == 0 {
		log.Println("No GPUs detected on this system")
		return nil
	}

	for i := 0; i < int(count); i++ {
		device, result := nvml.DeviceGetHandleByIndex(i)
		if result != nvml.SUCCESS {
			log.Printf("Failed to get GPU %d handle: %v", i, result)
			continue
		}

		util, _ := device.GetUtilizationRates()
		memInfo, _ := device.GetMemoryInfo()
		temp, _ := device.GetTemperature(nvml.TEMPERATURE_GPU)
		name, _ := device.GetName()

		metrics := map[string]any{
			"utilization_percent":    util.Gpu,
			"memory_used_megabytes":  memInfo.Used / MEGABYTE,
			"memory_total_megabytes": memInfo.Total / MEGABYTE,
			"temperature_celsius":    temp,
		}

		tags := map[string]string{
			"gpu_name":  name,
			"gpu_index": strconv.Itoa(i),
		}

		point := influxdb2.NewPoint("gpu", tags, metrics, time.Now())
		points = append(points, point)
	}

	return points
}
