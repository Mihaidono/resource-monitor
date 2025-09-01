package main

import (
	"log"

	"github.com/Mihaidono/resource-monitor/config"
	"github.com/Mihaidono/resource-monitor/workers"
	"github.com/NVIDIA/go-nvml/pkg/nvml"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Using system environment variables!")
	} else {
		log.Println("Loaded custom environment variables!")
	}

	influxCfg := config.LoadInfluxConfig()
	intervalCfg := config.LoadIntervalConfig()

	client := influxdb2.NewClient(influxCfg.URL, influxCfg.Token)
	defer client.Close()
	writeAPI := client.WriteAPI(influxCfg.Org, influxCfg.Bucket)

	if result := nvml.Init(); result != nvml.SUCCESS {
		log.Printf("Failed to initialize NVML: %v\nGPU monitoring will be disabled...", result)
	} else {
		defer nvml.Shutdown()
		log.Println("NVML initialized.")
		go workers.GpuWorker(writeAPI, intervalCfg.GPU)
	}

	go workers.CpuWorker(writeAPI, intervalCfg.CPU)

	select {}
}
