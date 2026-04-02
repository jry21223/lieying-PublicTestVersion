package scheduler

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/lieying/engine/internal/config"
	"github.com/lieying/engine/internal/models"
	"github.com/lieying/engine/pkg/mq"
	"github.com/lieying/engine/pkg/redis"
	"gorm.io/gorm"
)

type TaskScheduler struct {
	db       *gorm.DB
	redis    *redis.Client
	mq       *mq.Client
	cfg      config.EngineConfig
	taskChan chan models.Task
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

func New(db *gorm.DB, redis *redis.Client, mq *mq.Client, cfg config.EngineConfig) *TaskScheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &TaskScheduler{
		db:       db,
		redis:    redis,
		mq:       mq,
		cfg:      cfg,
		taskChan: make(chan models.Task, cfg.TaskQueueSize),
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (s *TaskScheduler) Start() {
	log.Println("Starting task scheduler...")

	if _, err := s.mq.DeclareQueue("tasks"); err != nil {
		log.Printf("Failed to declare queue: %v", err)
	}

	for i := 0; i < s.cfg.Workers; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}

	go s.pollPendingTasks()
}

func (s *TaskScheduler) Stop() {
	log.Println("Stopping task scheduler...")
	s.cancel()
	s.wg.Wait()
	close(s.taskChan)
}

func (s *TaskScheduler) Submit(task *models.Task) error {
	task.Status = models.TaskStatusPending
	if err := s.db.Save(task).Error; err != nil {
		return err
	}

	body, _ := json.Marshal(task)
	return s.mq.Publish("tasks", body)
}

func (s *TaskScheduler) worker(id int) {
	defer s.wg.Done()
	log.Printf("Worker %d started", id)

	for {
		select {
		case <-s.ctx.Done():
			log.Printf("Worker %d stopped", id)
			return
		case task := <-s.taskChan:
			s.executeTask(&task)
		}
	}
}

func (s *TaskScheduler) executeTask(task *models.Task) {
	log.Printf("Executing task: %s (%s)", task.ID, task.Type)

	now := time.Now()
	task.Status = models.TaskStatusRunning
	task.StartedAt = &now
	s.db.Save(task)

	switch task.Type {
	case models.TaskTypeRecon:
		s.runReconTask(task)
	case models.TaskTypeScan:
		s.runScanTask(task)
	case models.TaskTypeFuzz:
		s.runFuzzTask(task)
	case models.TaskTypeExploit:
		s.runExploitTask(task)
	case models.TaskTypeFullChain:
		s.runFullChainTask(task)
	}

	completed := time.Now()
	task.Status = models.TaskStatusCompleted
	task.CompletedAt = &completed
	task.Progress = 100
	s.db.Save(task)

	log.Printf("Task completed: %s", task.ID)
}

func (s *TaskScheduler) pollPendingTasks() {
	msgs, err := s.mq.Consume("tasks")
	if err != nil {
		log.Printf("Failed to consume queue: %v", err)
		return
	}

	for {
		select {
		case <-s.ctx.Done():
			return
		case msg := <-msgs:
			var task models.Task
			if err := json.Unmarshal(msg.Body, &task); err == nil {
				s.taskChan <- task
			}
		}
	}
}

func (s *TaskScheduler) runReconTask(task *models.Task) {
	log.Println("Running recon task...")
	task.Progress = 50
	s.db.Save(task)
	time.Sleep(2 * time.Second)
	task.Progress = 100
}

func (s *TaskScheduler) runScanTask(task *models.Task) {
	log.Println("Running scan task...")
	task.Progress = 30
	s.db.Save(task)
	time.Sleep(3 * time.Second)
	task.Progress = 100
}

func (s *TaskScheduler) runFuzzTask(task *models.Task) {
	log.Println("Running fuzz task...")
	task.Progress = 20
	s.db.Save(task)
	time.Sleep(5 * time.Second)
	task.Progress = 100
}

func (s *TaskScheduler) runExploitTask(task *models.Task) {
	log.Println("Running exploit task...")
	task.Progress = 40
	s.db.Save(task)
	time.Sleep(2 * time.Second)
	task.Progress = 100
}

func (s *TaskScheduler) runFullChainTask(task *models.Task) {
	log.Println("Running full chain task...")
	task.Progress = 10
	s.db.Save(task)
	time.Sleep(8 * time.Second)
	task.Progress = 100
}
