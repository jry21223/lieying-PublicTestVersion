package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

type TargetRepository struct {
	*Repository
}

type VulnerabilityRepository struct {
	*Repository
}

type AssetRepository struct {
	*Repository
}

func NewTargetRepository(db *gorm.DB) *TargetRepository {
	return &TargetRepository{Repository: NewRepository(db)}
}

func NewVulnerabilityRepository(db *gorm.DB) *VulnerabilityRepository {
	return &VulnerabilityRepository{Repository: NewRepository(db)}
}

func NewAssetRepository(db *gorm.DB) *AssetRepository {
	return &AssetRepository{Repository: NewRepository(db)}
}
