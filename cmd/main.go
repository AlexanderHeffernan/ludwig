package main

import (
	"ludwig/internal/orchestrator"
	"time"
)

func main() {
	orchestrator.Start()
	time.Sleep(10000 * time.Second)
	orchestrator.Stop()
}
