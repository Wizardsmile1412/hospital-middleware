package routes

import (
	"github.com/Wizardsmile1412/hospital-middleware/internal/config"
	"github.com/Wizardsmile1412/hospital-middleware/internal/handler"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, staffHandler *handler.StaffHandler, patientHandler *handler.PatientHandler, cfg *config.Config) {
	api := r.Group("/api/v1")
	RegisterStaffRoutes(api, staffHandler)
	RegisterPatientRoutes(api, patientHandler, cfg)
}
