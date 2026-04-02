package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type POCType string

const (
	POCTypeNuclei  POCType = "nuclei"
	POCTypeXray    POCType = "xray"
	POCTypeCustom  POCType = "custom"
	POCTypePython  POCType = "python"
	POCTypeGo      POCType = "go"
)

type POC struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Type        POCType        `gorm:"size:50;not null" json:"type"`
	Content     string         `gorm:"type:text;not null" json:"content"`
	VulnID      uuid.UUID      `gorm:"type:uuid;index" json:"vuln_id"`
	Author      string         `gorm:"size:255" json:"author"`
	Version     string         `gorm:"size:50" json:"version"`
	Tags        JSONB          `gorm:"type:jsonb" json:"tags"`
	Verified    bool           `gorm:"default:false" json:"verified"`
	Downloads   int            `gorm:"default:0" json:"downloads"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Vulnerability *Vulnerability `gorm:"foreignKey:VulnID" json:"vulnerability,omitempty"`
}

func (p *POC) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
