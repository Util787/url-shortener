package rest

import (
	"log/slog"

	"github.com/Util787/url-shortener/internal/common"
	"github.com/gin-gonic/gin"
)

type shortenerUsecase interface {
	SaveURL(longURL string) error
	GetLongURL(shortURL string) (string, error)
	GetRandomURL() (string, error)
	DeleteURL(id *string, longURL *string, shortURL *string) error
}

type Handler struct {
	log              *slog.Logger
	ShortenerUsecase shortenerUsecase
}

type SaveURLRequest struct {
	LongURL string `json:"long_url" binding:"required,url"`
}

func (h *Handler) SaveURL(c *gin.Context) {
	log := common.LogOpAndId(c.Request.Context(), common.GetOperationName(), h.log)

	var req SaveURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("failed to bind JSON", "error", err)
		newErrorResponse(c, log, 400, "failed to bind JSON", err)
		return
	}

	if err := h.ShortenerUsecase.SaveURL(req.LongURL); err != nil {
		log.Error("failed to save URL", "error", err)
		newErrorResponse(c, log, 500, "failed to save URL", err)
		return
	}

	c.JSON(200, gin.H{"message": "URL saved successfully"})
}

func (h *Handler) GetRandomURL(c *gin.Context) {
	log := common.LogOpAndId(c.Request.Context(), common.GetOperationName(), h.log)

	shortURL, err := h.ShortenerUsecase.GetRandomURL()
	if err != nil {
		log.Error("failed to get random URL", "error", err)
		newErrorResponse(c, log, 500, "failed to get random URL", err)
		return
	}

	c.JSON(200, gin.H{"short_url": shortURL})
}

func (h *Handler) RedirectURL(c *gin.Context) {
	log := common.LogOpAndId(c.Request.Context(), common.GetOperationName(), h.log)

	shortURLId := c.Param("short_url_id")
	longURL, err := h.ShortenerUsecase.GetLongURL(common.RedirectBaseURL + shortURLId)
	if err != nil {
		log.Error("failed to get long URL", "error", err)
		newErrorResponse(c, log, 500, "failed to get long URL", err)
		return
	}

	c.Redirect(302, longURL)
}

func (h *Handler) DeleteURL(c *gin.Context) {
	log := common.LogOpAndId(c.Request.Context(), common.GetOperationName(), h.log)

	shortURL := c.Param("short_url")
	if err := h.ShortenerUsecase.DeleteURL(nil, nil, &shortURL); err != nil {
		log.Error("failed to delete URL", "error", err)
		newErrorResponse(c, log, 500, "failed to delete URL", err)
		return
	}

	c.JSON(200, gin.H{"message": "URL deleted successfully"})
}
