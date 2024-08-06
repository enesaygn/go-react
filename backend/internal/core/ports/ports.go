package ports

import (
	"sasa-elterminali-service/internal/core/models"
)

// ArizaRepository is an interface that defines the methods that the ArizaRepository should implement
type ArizaRepository interface {
	GetAriza(ariza *models.GetArizaRequestBody) (*models.ArizaDetay, error)
}
