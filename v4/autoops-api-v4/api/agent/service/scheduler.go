// Agent调度引擎 - 核心调度器
package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"dodevops-api/api/agent/dao"
	"dodevops-api/api/agent/model"
	"dodevops-api/common/util"
)

// AgentHandler Agent处理器接口
type AgentHandler interface {
	Name() string                                    // 处理器名称
	Execute(ctx context.Context, task *model.AgentTask) error // 执行任务
	Validate(task *model.AgentTask) error            // 验证任务
	CanRetry(task *model.AgentTask, err error) bool  // 是否可以重试
}

// AgentScheduler Agent调度器
type AgentScheduler struct {
	taskQueue    chan *model.AgentTask            // 任务队列
	workerPool   *WorkerPool                      // 工作池
	taskRegistry map[string]AgentHandler          // 处理器注册表
	taskDao      dao.AgentTaskDao                 // 任务DAO
	mu           sync.RWMutex                     // 读写锁
	running      bool                             // 运行状态
	ctx          context.Context                  // 上下文
	cancel       context.CancelFunc               // 取消函数
}

var (
	schedulerInstance *AgentScheduler
	schedulerOnce     sync.Once
)

// GetAgentScheduler 获取调度器单例
func GetAgentScheduler() *AgentScheduler {
	schedulerOnce.Do(func() {
		schedulerInstance = NewAgentScheduler(20) // 默认20个工作协程
	})
	return schedulerInstance
}

// NewAgentScheduler 创建调度器
func NewAgentScheduler(workers int) *AgentScheduler {
	ctx, cancel := context.WithCancel(context.Background())

	scheduler := &AgentScheduler{
		taskQueue:    make(chan *model.AgentTask, 1000),
		workerPool:   NewWorkerPool(workers),
		taskRegistry: make(map[string]AgentHandler),
		taskDao:      dao.NewAgentTaskDao(),
		running:      false,
		ctx:          ctx,
		cancel:       cancel,
	}

	return scheduler
}

// RegisterHandler 注册处理器
func (s *AgentScheduler) RegisterHandler(handler AgentHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.taskRegistry[handler.Name()] = handler
}

// Start 启动调度器
func (s *AgentScheduler) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	// 启动工作池
	s.workerPool.Start(s.processTask)

	// 启动任务分发
	go s.dispatchTasks()

	// 启动定时任务扫描
	go s.scanPendingTasks()

	// 启动定时巡检扫描
	go s.scanScheduledInspections()
}

// Stop 停止调度器
func (s *AgentScheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	s.running = false
	s.cancel()
	close(s.taskQueue)
	s.workerPool.Stop()
}

// SubmitTask 提交任务
func (s *AgentScheduler) SubmitTask(task *model.AgentTask) error {
	// 1. 验证任务
	handler := s.getHandler(task.Type)
	if handler == nil {
		return fmt.Errorf("未找到处理器: %s", task.Type)
	}

	if err := handler.Validate(task); err != nil {
		return fmt.Errorf("任务验证失败: %v", err)
	}

	// 2. 保存任务
	task.Status = 1 // 待执行
	task.CreateTime = util.HTime{Time: time.Now()}
	task.UpdateTime = util.HTime{Time: time.Now()}

	if err := s.taskDao.Create(task); err != nil {
		return fmt.Errorf("保存任务失败: %v", err)
	}

	// 3. 放入队列
	select {
	case s.taskQueue <- task:
		return nil
	default:
		return fmt.Errorf("任务队列已满")
	}
}

// dispatchTasks 分发任务
func (s *AgentScheduler) dispatchTasks() {
	for {
		select {
		case task, ok := <-s.taskQueue:
			if !ok {
				return
			}
			// 提交到工作池
			s.workerPool.Submit(task)
		case <-s.ctx.Done():
			return
		}
	}
}

// processTask 处理任务
func (s *AgentScheduler) processTask(task *model.AgentTask) {
	// 1. 更新状态为执行中
	task.Status = 2 // 执行中
	task.StartTime = util.HTime{Time: time.Now()}
	s.taskDao.Update(task)

	// 2. 获取处理器
	handler := s.getHandler(task.Type)
	if handler == nil {
		s.failTask(task, fmt.Errorf("未找到处理器: %s", task.Type))
		return
	}

	// 3. 执行任务
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	err := handler.Execute(ctx, task)

	// 4. 处理结果
	if err != nil {
		// 判断是否可以重试
		if task.Retry < task.MaxRetry && handler.CanRetry(task, err) {
			task.Retry++
			task.Status = 1 // 重新排队
			s.taskDao.Update(task)

			// 指数退避
			backoff := time.Duration(task.Retry) * 10 * time.Second
			time.Sleep(backoff)

			// 重新提交
			select {
			case s.taskQueue <- task:
			default:
			}
		} else {
			s.failTask(task, err)
		}
	} else {
		s.completeTask(task)
	}
}

// scanPendingTasks 扫描待执行任务
func (s *AgentScheduler) scanPendingTasks() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 查询待执行的任务
			tasks, err := s.taskDao.GetPendingTasks(100)
			if err != nil {
				continue
			}

			for _, task := range tasks {
				select {
				case s.taskQueue <- &task:
				default:
					// 队列满了，下次再试
					break
				}
			}
		case <-s.ctx.Done():
			return
		}
	}
}

// scanScheduledInspections 扫描定时巡检任务
func (s *AgentScheduler) scanScheduledInspections() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// TODO: 实现定时巡检扫描逻辑
	for {
		select {
		case <-ticker.C:
			// 这里将在实现巡检Service时添加逻辑
		case <-s.ctx.Done():
			return
		}
	}
}

// completeTask 任务完成
func (s *AgentScheduler) completeTask(task *model.AgentTask) {
	task.Status = 3 // 成功
	task.EndTime = util.HTime{Time: time.Now()}
	task.UpdateTime = util.HTime{Time: time.Now()}
	s.taskDao.Update(task)
}

// failTask 任务失败
func (s *AgentScheduler) failTask(task *model.AgentTask, err error) {
	task.Status = 4 // 失败
	task.ErrorMsg = err.Error()
	task.EndTime = util.HTime{Time: time.Now()}
	task.UpdateTime = util.HTime{Time: time.Now()}
	s.taskDao.Update(task)
}

// getHandler 获取处理器
func (s *AgentScheduler) getHandler(taskType string) AgentHandler {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.taskRegistry[taskType]
}

// GetTaskStatus 获取任务状态
func (s *AgentScheduler) GetTaskStatus(taskID uint) (*model.AgentTask, error) {
	return s.taskDao.GetByID(taskID)
}

// CancelTask 取消任务
func (s *AgentScheduler) CancelTask(taskID uint) error {
	task, err := s.taskDao.GetByID(taskID)
	if err != nil {
		return err
	}

	if task.Status == 2 {
		return fmt.Errorf("任务正在执行中，无法取消")
	}

	task.Status = 5 // 取消
	task.UpdateTime = util.HTime{Time: time.Now()}
	return s.taskDao.Update(task)
}

// RetryTask 重试任务
func (s *AgentScheduler) RetryTask(taskID uint) error {
	task, err := s.taskDao.GetByID(taskID)
	if err != nil {
		return err
	}

	if task.Status != 4 {
		return fmt.Errorf("只能重试失败的任务")
	}

	task.Status = 1 // 待执行
	task.Retry = 0
	task.ErrorMsg = ""
	task.UpdateTime = util.HTime{Time: time.Now()}

	if err := s.taskDao.Update(task); err != nil {
		return err
	}

	// 重新提交
	select {
	case s.taskQueue <- task:
		return nil
	default:
		return fmt.Errorf("任务队列已满")
	}
}

// GetStats 获取调度器统计信息
func (s *AgentScheduler) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := map[string]interface{}{
		"running":       s.running,
		"queue_size":    len(s.taskQueue),
		"queue_cap":     cap(s.taskQueue),
		"handlers":      len(s.taskRegistry),
		"worker_pool":   s.workerPool.GetStats(),
	}

	// 获取任务统计
	taskStats, _ := s.taskDao.GetStatistics()
	stats["task_stats"] = taskStats

	return stats
}
