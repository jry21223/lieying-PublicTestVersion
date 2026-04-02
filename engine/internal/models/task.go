package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusPaused    TaskStatus = "paused"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
)

type TaskType string

const (
	TaskTypeRecon      TaskType = "recon"
	TaskTypeScan       TaskType = "scan"
	TaskTypeFuzz       TaskType = "fuzz"
	TaskTypeExploit    TaskType = "exploit"
	TaskTypeFullChain  TaskType = "full_chain"
)

type Task struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Type        TaskType       `gorm:"size:50;not null" json:"type"`
	Status      TaskStatus     `gorm:"size:50;not null;default:pending" json:"status"`
	TargetID    uuid.UUID      `gorm:"type:uuid;index" json:"target_id"`
	Config      JSONB          `gorm:"type:jsonb" json:"config"`
	Progress    int            `gorm:"default:0" json:"progress"`
	ResultCount int            `gorm:"default:0" json:"result_count"`
	StartedAt   *time.Time     `json:"started_at"`
	CompletedAt *time.Time     `json:"completed_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Target   Target       `gorm:"foreignKey:TargetID" json:"target,omitempty"`
	Results  []ScanResult `gorm:"foreignKey:TaskID" json:"results,omitempty"`
}

func (t *Task) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type JSONB map[string]interface{}
