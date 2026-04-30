package handler

import (
	"net/http"

	"github.com/Wizardsmile1412/hospital-middleware/internal/middleware"
	"github.com/Wizardsmile1412/hospital-middleware/internal/model"
	"github.com/Wizardsmile1412/hospital-middleware/internal/service"
	"github.com/Wizardsmile1412/hospital-middleware/internal/token"
	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	svc service.PatientService
}

func NewPatientHandler(svc service.PatientService) *PatientHandler {
	return &PatientHandler{svc: svc}
}

func (h *PatientHandler) SearchPatients(c *gin.Context) {
	claims := c.MustGet(middleware.ClaimsKey).(*token.Claims)

	params := model.SearchPatientParams{
		HospitalID:  claims.HospitalID,
		NationalID:  c.Query("national_id"),
		PassportID:  c.Query("passport_id"),
		FirstName:   c.Query("first_name"),
		MiddleName:  c.Query("middle_name"),
		LastName:    c.Query("last_name"),
		DateOfBirth: c.Query("date_of_birth"),
		PhoneNumber: c.Query("phone_number"),
		Email:       c.Query("email"),
	}

	result, err := h.svc.SearchPatients(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, result)
}
