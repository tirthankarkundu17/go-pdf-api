package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/tirthankarkundu17/pdf-api/api/handler"
	"go.uber.org/zap"
)

// New api.
func New() *gin.Engine {
	var (
		logger, _           = zap.NewProduction()
		router              = gin.New()
		pdfGeneratorHandler = handler.NewPDFGenerator(*logger)
	)

	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))
	router.Use(requestid.New())
	router.Use(cors.Default())

	pdfGeneratorHandler.Mount(router.Group("/pdfs"))

	return router
}
