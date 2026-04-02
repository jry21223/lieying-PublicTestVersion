package scheduler

import (
	"context"
	"sync"
	"time"
)

type TaskType string

const (
	TaskTypeHTTP      TaskType = "http"
	TaskTypePortScan  TaskType = "portscan"
	TaskTypeNuclei    TaskType = "nuclei"
	TaskTypeDirScan   TaskType = "dirscan"
	TaskTypeSubdomain TaskType = "subdomain"
)

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusPaused    TaskStatus = "paused"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
	TaskStatusCancelled TaskStatus = "cancelled"
)

type Task struct {
	ID          string
	Type        TaskType
	Target      string
	Options     map[string]interface{}
	Status      TaskStatus
	Progress    float64
	Result      interface{}
	Error       error
	Callback    func(result interface{})
	CreatedAt   time.Time
	UpdatedAt   time.Time
	StartedAt   *time.Time
	FinishedAt  *time.Time
}

type Scheduler struct {
	maxConcurrent int
	semaphore    chan struct{}
	tasks         chan *Task
	wg            sync.WaitGroup
	ctx           context.Context
	cancel        context.CancelFunc
	mu            sync.RWMutex
	runningTasks  map[string]*Task
	pausedTasks   map[string]*Task
}

func NewScheduler(maxConcurrent int) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		maxConcurrent: maxConcurrent,
		semaphore:    make(chan struct{}, maxConcurrent),
		tasks:        make(chan *Task, 1000),
		ctx:          ctx,
		cancel:       cancel,
		runningTasks: make(map[string]*Task),
		pausedTasks:  make(map[string]*Task),
	}
}

func (s *Scheduler) Start() {
	for i := 0; i < s.maxConcurrent; i++ {
		s.wg.Add(1)
		go s.worker()
	}

	go s.reaper()
}

func (s *Scheduler) worker() {
	defer s.wg.Done()

	for {
		select {
		case task := <-s.tasks:
			s.executeTask(task)
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *Scheduler) executeTask(task *Task) {
	s.semaphore <- struct{}{}
	defer func() { <-s.semaphore }()

	s.mu.Lock()
	task.Status = TaskStatusRunning
	now := time.Now()
	task.StartedAt = &now
	task.UpdatedAt = now
	s.runningTasks[task.ID] = task
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.runningTasks, task.ID)
		if task.FinishedAt == nil {
			finished := time.Now()
			task.FinishedAt = &finished
		}
		task.UpdatedAt = time.Now()
		s.mu.Unlock()
	}()

	result := s.runTask(task)

	s.mu.Lock()
	task.Result = result
	if task.Error != nil {
		task.Status = TaskStatusFailed
	} else {
		task.Status = TaskStatusCompleted
	}
	s.mu.Unlock()

	if task.Callback != nil {
		task.Callback(result)
	}
}

func (s *Scheduler) runTask(task *Task) interface{} {
	switch task.Type {
	case TaskTypeHTTP:
		return s.runHTTPTask(task)
	case TaskTypePortScan:
		return s.runPortScanTask(task)
	case TaskTypeNuclei:
		return s.runNucleiTask(task)
	case TaskTypeDirScan:
		return s.runDirScanTask(task)
	case TaskTypeSubdomain:
		return s.runSubdomainTask(task)
	default:
		task.Error = ErrUnknownTaskType
		return nil
	}
}

func (s *Scheduler) runHTTPTask(task *Task) interface{} {
	return map[string]string{"message": "HTTP task executed", "target": task.Target}
}

func (s *Scheduler) runPortScanTask(task *Task) interface{} {
	return map[string]string{"message": "Port scan task executed", "target": task.Target}
}

func (s *Scheduler) runNucleiTask(task *Task) interface{} {
	return map[string]string{"message": "Nuclei task executed", "target": task.Target}
}

func (s *Scheduler) runDirScanTask(task *Task) interface{} {
	return map[string]string{"message": "Directory scan task executed", "target": task.Target}
}

func (s *Scheduler) runSubdomainTask(task *Task) interface{} {
	return map[string]string{"message": "Subdomain task executed", "target": task.Target}
}

func (s *Scheduler) reaper() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.cleanup()
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *Scheduler) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, task := range s.runningTasks {
		if task.Status == TaskStatusCancelled {
			delete(s.runningTasks, id)
		}
	}
}

func (s *Scheduler) AddTask(task *Task) error {
	select {
	case s.tasks <- task:
		return nil
	default:
		return ErrTaskQueueFull
	}
}

func (s *Scheduler) AddTaskWithContext(ctx context.Context, task *Task) error {
	select {
	case s.tasks <- task:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return ErrTaskQueueFull
	}
}

func (s *Scheduler) CancelTask(taskID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task, exists := s.runningTasks[taskID]; exists {
		task.Status = TaskStatusCancelled
		return true
	}

	if task, exists := s.pausedTasks[taskID]; exists {
		task.Status = TaskStatusCancelled
		delete(s.pausedTasks, taskID)
		return true
	}

	return false
}

func (s *Scheduler) PauseTask(taskID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task, exists := s.runningTasks[taskID]; exists {
		task.Status = TaskStatusPaused
		s.pausedTasks[taskID] = task
		delete(s.runningTasks, taskID)
		return true
	}

	return false
}

func (s *Scheduler) ResumeTask(taskID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task, exists := s.pausedTasks[taskID]; exists {
		task.Status = TaskStatusPending
		delete(s.pausedTasks, taskID)
		s.runningTasks[taskID] = task
		s.tasks <- task
		return true
	}

	return false
}

func (s *Scheduler) GetTaskStatus(taskID string) *TaskStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if task, exists := s.runningTasks[taskID]; exists {
		status := task.Status
		return &status
	}

	if task, exists := s.pausedTasks[taskID]; exists {
		status := task.Status
		return &status
	}

	return nil
}

func (s *Scheduler) GetRunningTasks() []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*Task, 0, len(s.runningTasks))
	for _, task := range s.runningTasks {
		tasks = append(tasks, task)
	}

	return tasks
}

func (s *Scheduler) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"max_concurrent": s.maxConcurrent,
		"running_tasks":  len(s.runningTasks),
		"paused_tasks":   len(s.pausedTasks),
		"queue_length":  len(s.tasks),
	}
}

func (s *Scheduler) Stop() {
	s.cancel()
	s.wg.Wait()
}

func (s *Scheduler) SetMaxConcurrent(max int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.maxConcurrent = max
}

var (
	ErrTaskQueueFull     = &SchedulerError{Message: "task queue is full"}
	ErrUnknownTaskType   = &SchedulerError{Message: "unknown task type"}
	ErrTaskNotFound      = &SchedulerError{Message: "task not found"}
)

type SchedulerError struct {
	Message string
}

func (e *SchedulerError) Error() string {
	return e.Message
}
