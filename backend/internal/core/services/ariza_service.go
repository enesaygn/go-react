package services

import (
	"sasa-elterminali-service/internal/core/models"
	"sasa-elterminali-service/internal/core/ports"
)

type ArizaService struct {
	repo ports.ArizaRepository
}

func NewArizaService(repo ports.ArizaRepository) *ArizaService {
	return &ArizaService{repo: repo}
}
func (s *ArizaService) GetAriza(ariza *models.GetArizaRequestBody) (*models.ArizaDetay, error) {
	return s.repo.GetAriza(ariza)
}
