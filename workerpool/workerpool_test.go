package pool_test

import (
	"testing"
	"time"

	"github.com/MosinFAM/worker-pool/logger"
	pool "github.com/MosinFAM/worker-pool/workerpool"
)

func TestWorkerPool(t *testing.T) {
	logger.Init()
	wp := pool.NewWorkerPool(5)

	wp.AddWorker()
	wp.AddWorker()
	wp.Submit("task-1")
	wp.Submit("task-2")

	time.Sleep(1 * time.Second)

	wp.AddWorker()
	wp.Submit("task-3")
	wp.Submit("task-4")

	time.Sleep(1 * time.Second)
	wp.RemoveWorker()
	wp.StopAll()
}
