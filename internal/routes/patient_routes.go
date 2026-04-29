package routes

import (
	"github.com/Wizardsmile1412/hospital-middleware/internal/config"
	"github.com/Wizardsmile1412/hospital-middleware/internal/handler"
	"github.com/Wizardsmile1412/hospital-middleware/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterPatientRoutes(api *gin.RouterGroup, patientHandler *handler.PatientHandler, cfg *config.Config) {
	patient := api.Group("/patient")
	patient.Use(middleware.AuthRequired(cfg.JWTSecret))
	{
		patient.GET("/search", patientHandler.SearchPatients)
	}
}
