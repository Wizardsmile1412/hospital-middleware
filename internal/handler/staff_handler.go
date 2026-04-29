package handler

import (
	"errors"
	"net/http"

	"github.com/Wizardsmile1412/hospital-middleware/internal/config"
	"github.com/Wizardsmile1412/hospital-middleware/internal/model"
	"github.com/Wizardsmile1412/hospital-middleware/internal/service"
	"github.com/gin-gonic/gin"
)

type StaffHandler struct {
	svc service.StaffService
	cfg *config.Config
}

func NewStaffHandler(svc service.StaffService, cfg *config.Config) *StaffHandler {
	return &StaffHandler{svc: svc, cfg: cfg}
}

func (h *StaffHandler) CreateStaff(c *gin.Context) {
	var req model.CreateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	staffID, err := h.svc.CreateStaff(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrHospitalNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "hospital not found"})
		case errors.Is(err, service.ErrUsernameTaken):
			c.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Staff created successfully", "staff_id": staffID})
}

func (h *StaffHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenResp, refreshToken, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrHospitalNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "hospital not found"})
		case errors.Is(err, service.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	secure := h.cfg.AppEnv == "production"
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("refresh_token", refreshToken, 7*24*3600, "/", "", secure, true)
	c.JSON(http.StatusOK, tokenResp)
}

func (h *StaffHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found"})
		return
	}

	tokenResp, err := h.svc.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	c.JSON(http.StatusOK, tokenResp)
}
