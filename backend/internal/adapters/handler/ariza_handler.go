package handler

import (
	"encoding/json"
	"sasa-elterminali-service/internal/core/models"
	"sasa-elterminali-service/internal/core/services"
	"sasa-elterminali-service/internal/messaging"

	"github.com/gofiber/fiber/v2"
)

type ArizaHandler struct {
	service   services.ArizaService
	messaging messaging.Messaging
}

func NewArizaHandler(arizaService services.ArizaService, messaging messaging.Messaging) *ArizaHandler {
	return &ArizaHandler{
		service:   arizaService,
		messaging: messaging,
	}
}

// GetAriza godoc
// @Summary Get ariza by ID
// @Description Ariza getir
// @Tags ariza
// @Accept  json
// @Produce  json
// @Param body body models.GetArizaRequestBody true "ArizaID"
// @Success 200 {object} models.ArizaDetay
// @Router /api/v1/ariza/get [post]
// @Security Bearer
func (h *ArizaHandler) GetAriza(c *fiber.Ctx) error {
	var req models.GetArizaRequestBody

	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}
	ariza, err := h.service.GetAriza(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(ariza)
}
