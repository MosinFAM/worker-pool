package pool

import (
	"context"
	"sync"

	"github.com/MosinFAM/worker-pool/logger"
	"github.com/sirupsen/logrus"
)

type Worker struct {
	id     int
	ctx    context.Context
	cancel context.CancelFunc
}

type WorkerPool struct {
	mu      sync.Mutex
	workers map[int]*Worker
	input   chan string
	nextID  int
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewWorkerPool(bufferSize int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		workers: make(map[int]*Worker),
		input:   make(chan string, bufferSize),
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (wp *WorkerPool) AddWorker() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	id := wp.nextID
	wp.nextID++

	ctx, cancel := context.WithCancel(wp.ctx)
	worker := &Worker{id: id, ctx: ctx, cancel: cancel}
	wp.workers[id] = worker

	wp.wg.Add(1)
	go func(w *Worker) {
		defer wp.wg.Done()
		logger.LogInfo("Worker started", logrus.Fields{"worker_id": w.id})
		for {
			select {
			case <-w.ctx.Done():
				logger.LogInfo("Worker stopped", logrus.Fields{"worker_id": w.id})
				return
			case data, ok := <-wp.input:
				if !ok {
					return
				}
				logger.LogInfo("Processing data", logrus.Fields{
					"worker_id": w.id,
					"data":      data,
				})
			}
		}
	}(worker)
}

func (wp *WorkerPool) RemoveWorker() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	for id, worker := range wp.workers {
		worker.cancel()
		delete(wp.workers, id)
		logger.LogInfo("Worker removed", logrus.Fields{"worker_id": id})
		break
	}
}

func (wp *WorkerPool) Submit(data string) {
	wp.input <- data
}

func (wp *WorkerPool) StopAll() {
	logger.LogInfo("Stopping all workers", logrus.Fields{
		"total_workers": len(wp.workers),
	})
	wp.cancel()

	wp.mu.Lock()
	for id, worker := range wp.workers {
		worker.cancel()
		delete(wp.workers, id)
	}
	wp.mu.Unlock()

	close(wp.input)
	wp.wg.Wait()
	logger.LogInfo("All workers stopped.", logrus.Fields{"status": "OK"})
}
