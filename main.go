package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type CallbackData struct {
	TaskID        string  `json:"taskId"`
	DeviceID      int     `json:"deviceId"`
	WorkloadSize  int     `json:"workloadSize"`
	ExecutionSite string  `json:"executionSite"`
	CreatedAt     int64   `json:"createdAt"`
	Duration      float64 `json:"duration"`
}

func main() {
	log.Println("Running task...")

	callbackData := new(CallbackData)

	workloadSize, err := strconv.Atoi(os.Getenv("WORKLOAD_SIZE"))
	callbackData.WorkloadSize = workloadSize

	if err != nil {
		log.Fatal("Failed to parse workload size: ", err)
	}

	callbackData.CreatedAt = time.Now().UnixNano()
	callbackData.TaskID = os.Getenv("TASK_ID")
	callbackData.ExecutionSite = os.Getenv("EXECUTION_SITE")
	callbackData.DeviceID, err = strconv.Atoi(os.Getenv("DEVICE_ID"))

	if err != nil {
		log.Println("Failed to parse device id: ", err, " Using default value: -1")
		callbackData.DeviceID = -1
	}

	callbackData.Duration = cpuBoundWork(workloadSize).Seconds()

	callbackAddr := os.Getenv("CALLBACK_ADDR")
	sendCallback(callbackData, callbackAddr)

	log.Println("Task finished.")

	os.Exit(0)
}

func cpuBoundWork(n int) time.Duration {
	start := time.Now()

	x := uint64(0xDEADBEEF)
	for i := 0; i < n; i++ {
		x ^= uint64(i) * 0x9e3779b97f4a7c15
		x ^= x >> 33
		x *= 0xff51afd7ed558ccd
		x ^= x >> 33
		x *= 0xc4ceb9fe1a85ec53
		x ^= x >> 33
	}

	_ = x

	return time.Since(start)
}

func sendCallback(data *CallbackData, addr string) {
	log.Printf("Sending callback: %+v\n", *data)

	body, err := json.Marshal(*data)

	if err != nil {
		log.Fatal("Failed to parse body: ", err)
	}

	resp, err := http.Post(addr, "application/json", bytes.NewBuffer(body))

	if err != nil {
		log.Fatal("Failed to send callback: ", err)
	}

	log.Println("Response status:", resp.StatusCode)
	log.Println("Callback sent successfully.")
}
