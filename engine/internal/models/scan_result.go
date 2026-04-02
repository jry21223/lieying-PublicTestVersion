package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScanResult struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	TaskID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"task_id"`
	TargetID       uuid.UUID      `gorm:"type:uuid;index" json:"target_id"`
	VulnID         *uuid.UUID     `gorm:"type:uuid;index" json:"vuln_id"`
	POCID          *uuid.UUID     `gorm:"type:uuid;index" json:"poc_id"`
	Title          string         `gorm:"size:500;not null" json:"title"`
	Description    string         `gorm:"type:text" json:"description"`
	Severity       Severity       `gorm:"size:20;not null" json:"severity"`
	URL            string         `gorm:"size:1000" json:"url"`
	Payload        string         `gorm:"type:text" json:"payload"`
	Evidence       string         `gorm:"type:text" json:"evidence"`
	Request        string         `gorm:"type:text" json:"request"`
	Response       string         `gorm:"type:text" json:"response"`
	Confirmed      bool           `gorm:"default:false" json:"confirmed"`
	FalsePositive  bool           `gorm:"default:false" json:"false_positive"`
	Metadata       JSONB          `gorm:"type:jsonb" json:"metadata"`
	DiscoveredAt   time.Time      `json:"discovered_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	Task          *Task          `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	Target        *Target        `gorm:"foreignKey:TargetID" json:"target,omitempty"`
	Vulnerability *Vulnerability `gorm:"foreignKey:VulnID" json:"vulnerability,omitempty"`
	POC           *POC           `gorm:"foreignKey:POCID" json:"poc,omitempty"`
}

func (s *ScanResult) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	if s.DiscoveredAt.IsZero() {
		s.DiscoveredAt = time.Now()
	}
	return nil
}
