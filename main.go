package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/MosinFAM/worker-pool/logger"
	pool "github.com/MosinFAM/worker-pool/workerpool"
)

func main() {
	logger.Init()
	wp := pool.NewWorkerPool(10)

	wp.AddWorker()
	wp.AddWorker()
	wp.AddWorker()

	go func() {
		for i := 0; i < 10; i++ {
			wp.Submit("task-" + time.Now().Format("15:04:05"))
			time.Sleep(500 * time.Millisecond)
		}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		wp.AddWorker()
		time.Sleep(2 * time.Second)
		wp.RemoveWorker()
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	println("\nReceived interrupt signal, shutting down...")
	wp.StopAll()
	println("Clean shutdown completed")
}
