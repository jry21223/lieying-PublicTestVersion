package service

import (
	"github.com/kunlun-sec/lunying/internal/repository"
)

type TargetService struct {
	targetRepo *repository.TargetRepository
}

func NewTargetService(targetRepo *repository.TargetRepository) *TargetService {
	return &TargetService{targetRepo: targetRepo}
}

type VulnerabilityService struct {
	vulnRepo *repository.VulnerabilityRepository
}

func NewVulnerabilityService(vulnRepo *repository.VulnerabilityRepository) *VulnerabilityService {
	return &VulnerabilityService{vulnRepo: vulnRepo}
}

type AssetService struct {
	assetRepo *repository.AssetRepository
}

func NewAssetService(assetRepo *repository.AssetRepository) *AssetService {
	return &AssetService{assetRepo: assetRepo}
}
