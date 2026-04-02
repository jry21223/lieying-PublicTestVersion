package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TargetType string

const (
	TargetTypeDomain   TargetType = "domain"
	TargetTypeIP       TargetType = "ip"
	TargetTypeCIDR     TargetType = "cidr"
	TargetTypeURL      TargetType = "url"
	TargetTypePlatform TargetType = "platform"
)

type Target struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Type        TargetType     `gorm:"size:50;not null" json:"type"`
	Value       string         `gorm:"size:500;not null;index" json:"value"`
	Description string         `gorm:"type:text" json:"description"`
	Tags        JSONB          `gorm:"type:jsonb" json:"tags"`
	Metadata    JSONB          `gorm:"type:jsonb" json:"metadata"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Tasks []Task `gorm:"foreignKey:TargetID" json:"tasks,omitempty"`
}

func (t *Target) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
