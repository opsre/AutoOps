// Agent调度引擎 - 工作池
package service

import (
	"context"
	"sync"
	"time"

	"dodevops-api/api/agent/model"
)

// WorkerPool 工作池
type WorkerPool struct {
	workers  int                          // 工作协程数
	taskChan chan *model.AgentTask        // 任务通道
	handler  func(*model.AgentTask)       // 任务处理函数
	wg       sync.WaitGroup               // 等待组
	ctx      context.Context              // 上下文
	cancel   context.CancelFunc           // 取消函数
	mu       sync.RWMutex                 // 读写锁
	running  bool                         // 运行状态
}

// NewWorkerPool 创建工作池
func NewWorkerPool(workers int) *WorkerPool {
	if workers <= 0 {
		workers = 10 // 默认10个工作协程
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		workers:  workers,
		taskChan: make(chan *model.AgentTask, 1000), // 缓冲队列
		ctx:      ctx,
		cancel:   cancel,
		running:  false,
	}
}

// Start 启动工作池
func (wp *WorkerPool) Start(handler func(*model.AgentTask)) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.running {
		return
	}

	wp.handler = handler
	wp.running = true

	// 启动工作协程
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// Stop 停止工作池
func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if !wp.running {
		return
	}

	wp.running = false
	wp.cancel()
	close(wp.taskChan)
	wp.wg.Wait()
}

// Submit 提交任务到工作池
func (wp *WorkerPool) Submit(task *model.AgentTask) error {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	if !wp.running {
		return nil
	}

	select {
	case wp.taskChan <- task:
		return nil
	case <-time.After(5 * time.Second):
		return nil
	case <-wp.ctx.Done():
		return wp.ctx.Err()
	}
}

// worker 工作协程
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for {
		select {
		case task, ok := <-wp.taskChan:
			if !ok {
				return
			}
			if wp.handler != nil {
				wp.handler(task)
			}
		case <-wp.ctx.Done():
			return
		}
	}
}

// GetStats 获取工作池统计信息
func (wp *WorkerPool) GetStats() map[string]interface{} {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	return map[string]interface{}{
		"workers":     wp.workers,
		"running":     wp.running,
		"queue_size":  len(wp.taskChan),
		"queue_cap":   cap(wp.taskChan),
	}
}
