package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/tirthankarkundu17/pdf-api/internal/models"
	"github.com/tirthankarkundu17/pdf-api/pdfservice"
	"go.uber.org/zap"
)

var (
	// ErrBadRequest error.
	ErrBadRequest = errors.New("bad request")
)

// Score for score endpoints.
type PDFGenerator struct {
	logger     zap.Logger
	pdfService pdfservice.Service
}

// handle POST /generate-from-image
func (h PDFGenerator) GeneratePDFFromImage(c *gin.Context) {
	var pdfData models.PDFData
	if err := c.BindJSON(&pdfData); err != nil {
		h.logger.Error("Error", zap.Error(err))
		c.JSON(400, models.Error{
			Error: err.Error(),
		})
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(pdfData.Data)
	if err != nil {
		panic(err)
	}

	fPath := fmt.Sprintf("./uploads/%s.png", uuid.New().String())

	os.WriteFile(fPath, decoded, fs.FileMode(os.O_CREATE))

	data, err := h.pdfService.CreateFromImage(c, fPath)
	if err != nil {
		h.logger.Error("Error", zap.Error(err))
		c.JSON(400, models.Error{
			Error: err.Error(),
		})
		return
	}

	os.Remove(fPath)
	c.Header("Content-Disposition", "attachment; filename=file.pdf")
	c.Data(http.StatusOK, "application/octet-stream", data)
}

// handle POST /generate-from-text
func (h PDFGenerator) GeneratePDFFromText(c *gin.Context) {

}

// Mount handlers to router group.
func (h PDFGenerator) Mount(router *gin.RouterGroup) {
	router.POST("/generate-from-image", h.GeneratePDFFromImage)
}

// PDF handler.
func NewPDFGenerator(logger zap.Logger) PDFGenerator {
	s := pdfservice.New()
	return PDFGenerator{
		logger:     logger,
		pdfService: s,
	}
}
